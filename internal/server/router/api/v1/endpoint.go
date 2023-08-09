package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/GNews"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/NEWSDATA"
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
			HeadConent: view.SharedHeadContent,
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
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}
	global.Logger.Info().
		Str("user_id", userInfo.GetUserID().String()).
		Str("user_role", userInfo.GetRole().String()).
		Msg("Get User Info OK")

	if err := req.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}
	global.Logger.Info().
		Msg("Parse form OK")

	obj, err := obj.FormDecodeAndValidate(repo.apiRepo.FormDecoder, repo.val, req.PostForm)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}

	repo.apiRepo.Service.Job().Create(
		req.Context(), &service.JobCreateRequest{
			Owner:    userInfo.GetUserID(),
			Status:   string(model.JobStatusCreated),
			SrcApiID: apikey.ApiID,
			SrcQuery: q.Params().ToQueryString(),
		},
	)

	fmt.Fprintf(w,
		"User: %s(%s)\nAPI: %s\nEndpoint: %s\nQuery: %s\n%s\n",
		userInfo.GetUsername(),
		userInfo.GetRole(),
		obj.API(),
		obj.Endpoint(),
		q.Params().ToQueryString(),
		obj.String(),
	)
	return
}
