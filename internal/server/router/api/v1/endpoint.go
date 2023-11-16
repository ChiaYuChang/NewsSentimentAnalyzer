package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pf "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	val "github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var ErrEndpointNotFount = errors.New("unregistered endpoint")

type EndpointRepoKey [2]string

func NewRepoMapKey(apiName string, endpointName string) EndpointRepoKey {
	return EndpointRepoKey{apiName, endpointName}
}

func (k EndpointRepoKey) APIName() string {
	return k[0]
}

func (k EndpointRepoKey) EndpointName() string {
	return k[1]
}

func (k EndpointRepoKey) String() string {
	return fmt.Sprintf("%s-%s", k[0], k[1])
}

func (k EndpointRepoKey) CacheKey(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)

	return fmt.Sprintf(
		"%s-%s-%s",
		k[0], k[1],
		base64.StdEncoding.EncodeToString(hasher.Sum(nil)))
}

type EndpointRepo struct {
	apiRepo          APIRepo
	val              *val.Validate
	PageView         map[EndpointRepoKey]string
	PageSelectOption map[EndpointRepoKey]string
	PageTemplateName map[EndpointRepoKey]string
	ApiID            map[EndpointRepoKey]int16
}

func NewEndpointRepo(apiRepo APIRepo, v *val.Validate) EndpointRepo {
	return EndpointRepo{
		apiRepo:          apiRepo,
		val:              v,
		PageView:         make(map[EndpointRepoKey]string),
		PageSelectOption: make(map[EndpointRepoKey]string),
		PageTemplateName: make(map[EndpointRepoKey]string),
		ApiID:            make(map[EndpointRepoKey]int16),
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

	f, err := pf.PageFormRepo.Get(apiName, endpointName)
	if err != nil {
		return err
	}
	pageData.SelectOpts = f.SelectionOpts()

	key := NewRepoMapKey(apiName, endpointName)

	repo.PageTemplateName[key] = templateName

	buf := bytes.NewBuffer(nil)
	gz, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
	err = repo.apiRepo.View.ExecuteTemplate(gz, templateName, pageData)
	if err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return fmt.Errorf("error while compressing page: %w", err)
	}
	ctext := buf.Bytes()
	repo.PageView[key] = key.CacheKey(ctext)
	repo.apiRepo.Cache.Set(
		context.Background(), repo.PageView[key],
		ctext, global.CacheExpireDefault)

	buf = bytes.NewBuffer(nil)
	err = repo.apiRepo.View.ExecuteTemplate(buf, "options.gotmpl", pageData.SelectOpts)
	if err != nil {
		return err
	}
	repo.PageSelectOption[key] = key.CacheKey(buf.Bytes())
	repo.apiRepo.Cache.Set(
		context.Background(), repo.PageSelectOption[key],
		buf.Bytes(), global.CacheExpireDefault)

	repo.ApiID[key] = apiId
	return nil
}

func (repo EndpointRepo) getCache(ctx context.Context, key EndpointRepoKey, ckey string) ([]byte, error) {
	page, err := repo.apiRepo.Cache.Get(context.Background(), ckey).Bytes()
	switch err {
	default:
		return nil, fmt.Errorf("error while read cache: %w", err)
	case redis.Nil:
		if err := repo.RegisterEndpointsPageView(
			key.APIName(),
			repo.ApiID[key],
			key.EndpointName(),
			repo.PageTemplateName[key]); err != nil {
			return nil, fmt.Errorf("error excute template: %w", err)
		} else {
			page, _ = repo.apiRepo.Cache.Get(context.Background(), ckey).Bytes()
		}
	case nil:
		repo.apiRepo.Cache.ExpireGT(context.Background(), ckey, global.CacheExpireLong)
	}
	return page, nil
}

func (repo EndpointRepo) GetAPIEndpoints(key EndpointRepoKey) (http.HandlerFunc, error) {
	ckey, ok := repo.PageView[key]
	if !ok {
		return nil, ErrEndpointNotFount
	}

	page, err := repo.getCache(context.Background(), key, ckey)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(page)
	}, nil
}

func (repo EndpointRepo) GetAPISelectOptions(key EndpointRepoKey) (http.HandlerFunc, error) {
	ckey, ok := repo.PageSelectOption[key]
	if !ok {
		return nil, ErrEndpointNotFount
	}

	page, err := repo.getCache(context.Background(), key, ckey)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(page)
	}, nil
}

func (repo EndpointRepo) PostAPIEndpoints(key EndpointRepoKey) (http.HandlerFunc, error) {
	pf, err := pf.Get(key.APIName(), key.EndpointName())
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
		userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
		if !ok {
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			ecErr.WithDetails("user information not found")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}
		postEndpoints(repo, pf, userInfo, w, req)
	}, nil
}

func (repo EndpointRepo) PatchAPIEndpoints(key EndpointRepoKey) (http.HandlerFunc, error) {
	pf, err := pf.Get(key.APIName(), key.EndpointName())
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
		// userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
		// if !ok {
		// 	ecErr := ec.MustGetEcErr(ec.ECServerError)
		// 	ecErr.WithDetails("user information not found")
		// 	w.WriteHeader(ecErr.HttpStatusCode)
		// 	w.Write(ecErr.MustToJson())
		return
	}, nil
}

func postEndpoints(repo EndpointRepo, pageform pf.PageForm, userInfo tokenmaker.Payload, w http.ResponseWriter, httpReq *http.Request) {
	type JSONResponse struct {
		CacheKey          string            `json:"cache_key"`
		StatusCode        int               `json:"status_code"`
		Message           string            `json:"message"`
		Details           []string          `json:"details"`
		NewsSourseRequest api.Request       `json:"news-src-request"`
		LLMRequest        map[string]string `json:"llm-request"`
	}

	var jsonResp JSONResponse
	w.Header().Add("Content-Type", "application/json")

	userInfo, ok := httpReq.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		jsonResp.StatusCode = ecErr.HttpStatusCode
		jsonResp.Message = ecErr.Message
		jsonResp.Details = append(jsonResp.Details, "user information not found")
		b, _ := json.MarshalIndent(jsonResp, "", "    ")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	global.Logger.Info().
		Str("user_id", userInfo.GetUserID().String()).
		Str("user_role", userInfo.GetRole().String()).
		Msg("Get User Info OK")

	if err := httpReq.ParseForm(); err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error while calling .ParseForm method")
		return
	}
	global.Logger.Info().
		Msg("Parse form OK")

	pageform, err := pageform.FormDecodeAndValidate(
		repo.apiRepo.FormDecoder, repo.val, httpReq.PostForm)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		jsonResp.StatusCode = ecErr.HttpStatusCode
		jsonResp.Message = ecErr.Message
		jsonResp.Details = append(jsonResp.Details, err.Error())
		b, _ := json.MarshalIndent(jsonResp, "", "    ")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}
	global.Logger.Info().
		Msg("Decode and Validate OK")

	apikey, _ := repo.apiRepo.Service.APIKey().Get(
		httpReq.Context(), &service.APIKeyGetRequest{
			Owner: userInfo.GetUserID(),
			ApiID: repo.ApiID[NewRepoMapKey(pageform.API(), pageform.Endpoint())],
		},
	)

	req, err := client.PageFormHandlerRepo.Handle(apikey.Key, pageform)
	if err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error while calling .NewQueryFromPageFrom method")

		ecErr := ec.MustGetEcErr(ec.ECServerError)
		jsonResp.StatusCode = ecErr.HttpStatusCode
		jsonResp.Message = ecErr.Message
		jsonResp.Details = append(jsonResp.Details, err.Error())
		b, _ := json.MarshalIndent(jsonResp, "", "    ")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	ckey, cache := req.ToPreviewCache(userInfo.GetUserID())
	if repo.apiRepo.Cache != nil {
		_, _ = repo.apiRepo.Cache.JSONSet(ckey, ".", cache)
		_ = repo.apiRepo.Cache.Expire(httpReq.Context(), ckey, 10*time.Minute)
		jsonResp.CacheKey = ckey
	} else {
		global.Logger.Warn().
			Str("uid", cache.Query.UserId.String()).
			Time("created_at", cache.CreatedAt).
			Msg("Cache is nil")
	}

	llmRequest := map[string]string{
		"id":    "1",
		"query": "{}",
	}

	jsonResp.StatusCode = http.StatusOK
	jsonResp.Message = "OK"

	jsonResp.LLMRequest = llmRequest
	jsonResp.NewsSourseRequest = req

	b, err := json.MarshalIndent(jsonResp, "", "    ")
	if err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error while calling .NewQueryFromPageFrom method")

		ecErr := ec.MustGetEcErr(ec.ECServerError)
		jsonResp.StatusCode = ecErr.HttpStatusCode
		jsonResp.Message = ecErr.Message
		jsonResp.Details = append(jsonResp.Details, err.Error())
		b, _ := json.MarshalIndent(jsonResp, "", "    ")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	cookie := repo.apiRepo.CookieMaker.NewCookie(COOKIE_PREVIEW_CID, ckey)
	http.SetCookie(w, cookie)

	w.WriteHeader(jsonResp.StatusCode)
	w.Write(b)
	return
}

// func patchEndpoints() {
// 	client.PageFormHandlerRepo.
// }
