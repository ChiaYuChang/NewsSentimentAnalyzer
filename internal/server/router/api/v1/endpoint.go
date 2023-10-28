package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/google/uuid"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	val "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var ErrEndpointNotFount = errors.New("unregistered endpoint")

type repoMapKey [2]string

func newRepoMapKey(apiName string, endpointName string) repoMapKey {
	return repoMapKey{apiName, endpointName}
}

func (k repoMapKey) APIName() string {
	return k[0]
}

func (k repoMapKey) EndpointName() string {
	return k[1]
}

func (k repoMapKey) String() string {
	return fmt.Sprintf("%s-%s", k[0], k[1])
}

type EndpointRepo struct {
	apiRepo  APIRepo
	val      *val.Validate
	PageView map[repoMapKey][]byte
	ApiID    map[repoMapKey]int16
}

func NewEndpointRepo(apiRepo APIRepo, v *val.Validate) EndpointRepo {
	return EndpointRepo{
		apiRepo:  apiRepo,
		val:      v,
		PageView: make(map[repoMapKey][]byte),
		ApiID:    make(map[repoMapKey]int16),
	}
}

func (repo *EndpointRepo) RegisterEndpointsPageView(apiName string, apiId int16, endpointName, templateName string) error {
	pageData := object.APIEndpointPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      endpointName,
		},
		API:      apiName,
		Version:  viper.GetString("APP_API_VERSION"),
		Endpoint: endpointName,
	}

	buffer := bytes.NewBufferString("")
	err := repo.apiRepo.View.ExecuteTemplate(buffer, templateName, pageData)
	if err != nil {
		return err
	}

	key := newRepoMapKey(apiName, endpointName)
	repo.PageView[key] = buffer.Bytes()
	repo.ApiID[key] = apiId
	return nil
}

func (repo *EndpointRepo) WritePageViewTo(apiName, endpointName string, w io.Writer) error {
	var view []byte
	var ok bool
	var err error
	if view, ok = repo.PageView[newRepoMapKey(apiName, endpointName)]; !ok {
		return ErrEndpointNotFount
	}
	_, err = w.Write(view)
	return err
}

func (repo EndpointRepo) GetAPIEndpoints(key repoMapKey) (http.HandlerFunc, error) {
	page, ok := repo.PageView[key]
	if !ok {
		return nil, ec.MustGetEcErr(ec.ECNotFound)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(page)
	}, nil
}

func (repo EndpointRepo) PostAPIEndpoints(key repoMapKey) (http.HandlerFunc, error) {
	pf, err := pageform.Get(key.APIName(), key.EndpointName())
	if err != nil || pf == nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		if err != nil {
			ecErr.WithDetails(err.Error())
			ecErr.WithDetails(key.String() + " not found")
		}
		if pf == nil {
			ecErr.WithDetails("pageform object is nil")
		}
		return nil, ecErr
	}

	return func(w http.ResponseWriter, req *http.Request) {
		postEndpoints(repo, pf, w, req)
	}, nil
}

func postEndpoints(repo EndpointRepo, obj pageform.PageForm, w http.ResponseWriter, req *http.Request) {
	type resp struct {
		StatusCode      int       `json:"status_code"`
		Message         string    `json:"message"`
		Details         []string  `json:"details"`
		JobId           int32     `json:"job_id"`
		UserName        string    `json:"user_name"`
		UserId          uuid.UUID `json:"user_id"`
		NewsSourceId    int16     `json:"news_source_id"`
		NewsSourceQuery string    `json:"news_source_query"`
		LLMId           int16     `json:"llm_id"`
		LLMQuery        string    `json:"llm_query"`
	}

	r := resp{Details: []string{}}
	w.Header().Add("Content-Type", "application/json")

	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		r.StatusCode = ecErr.HttpStatusCode
		r.Message = ecErr.Message
		r.Details = append(r.Details, "user information not found")
		b, _ := json.MarshalIndent(r, "", "    ")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	global.Logger.Info().
		Str("user_id", userInfo.GetUserID().String()).
		Str("user_role", userInfo.GetRole().String()).
		Msg("Get User Info OK")

	if err := req.ParseForm(); err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error while calling .ParseForm method")
		return
	}
	global.Logger.Info().
		Msg("Parse form OK")

	obj, err := obj.FormDecodeAndValidate(repo.apiRepo.FormDecoder, repo.val, req.PostForm)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		r.StatusCode = ecErr.HttpStatusCode
		r.Message = ecErr.Message
		r.Details = append(r.Details, err.Error())
		b, _ := json.MarshalIndent(r, "", "    ")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}
	global.Logger.Info().
		Msg("Decode and Validate OK")

	apikey, _ := repo.apiRepo.Service.APIKey().Get(
		req.Context(), &service.APIKeyGetRequest{
			Owner: userInfo.GetUserID(),
			ApiID: repo.ApiID[newRepoMapKey(obj.API(), obj.Endpoint())],
		},
	)

	q, err := client.PageFormHandlerRepo.NewQueryFromPageFrom(apikey.Key, obj)
	if err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error while calling .NewQueryFromPageFrom method")

		ecErr := ec.MustGetEcErr(ec.ECServerError)
		r.StatusCode = ecErr.HttpStatusCode
		r.Message = ecErr.Message
		r.Details = append(r.Details, err.Error())
		b, _ := json.MarshalIndent(r, "", "    ")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	// jsn, err := json.Marshal(q.Params())
	// if err != nil {
	// 	global.Logger.Error().
	// 		Err(err).
	// 		Msg("error while calling Marshaling params")
	// 	return
	// }
	llmId, llmQuery := int16(4), "{}"
	id, err := repo.apiRepo.Service.Job().Create(
		req.Context(), &service.JobCreateRequest{
			Owner:    userInfo.GetUserID(),
			Status:   string(model.JobStatusCreated),
			SrcApiID: apikey.ApiID,
			// SrcQuery: q.Params().ToQueryString(),
			// SrcQuery: string(jsn),
			SrcQuery: q.Params().ToQueryString(),
			LlmApiID: llmId,
			LlmQuery: llmQuery,
		},
	)

	if err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error while creating job")
		return
	}

	global.Logger.Info().
		Str("user", userInfo.GetUsername()).
		Str("role", userInfo.GetRole().String()).
		Str("api", obj.API()).
		Str("endpoint", obj.Endpoint()).
		Str("query", q.Params().ToQueryString()).
		Int32("job", id).
		Msg("Job Created OK")

	r.StatusCode = http.StatusOK
	r.Message = "OK"
	r.JobId = id
	r.LLMId = llmId
	r.LLMQuery = llmQuery
	r.NewsSourceId = apikey.ApiID
	r.NewsSourceQuery = q.Params().ToQueryString()
	r.UserId = userInfo.GetUserID()
	r.UserName = userInfo.GetUsername()

	b, _ := json.MarshalIndent(r, "", "    ")
	w.WriteHeader(r.StatusCode)
	w.Write(b)
	return
}

