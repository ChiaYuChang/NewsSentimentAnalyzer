package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/base64"
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
	PageTemplateName map[EndpointRepoKey]string
	ApiID            map[EndpointRepoKey]int16
	EndpointID       map[EndpointRepoKey]int32
	PageSelectOption map[string]string
}

func NewEndpointRepo(apiRepo APIRepo, v *val.Validate) EndpointRepo {
	return EndpointRepo{
		apiRepo:          apiRepo,
		val:              v,
		PageView:         make(map[EndpointRepoKey]string),
		PageTemplateName: make(map[EndpointRepoKey]string),
		ApiID:            make(map[EndpointRepoKey]int16),
		EndpointID:       make(map[EndpointRepoKey]int32),
		PageSelectOption: make(map[string]string),
	}
}

func (repo *EndpointRepo) RegisterEndpointsPageView(
	apiName string, apiId int16,
	endpointName string, endpointId int32,
	templateName string) error {

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
	if err = repo.genPageView(key, templateName, pageData); err != nil {
		return fmt.Errorf("error while generating page view: %w", err)
	}

	if err = repo.genSelectOptionsScript(key, pageData, true); err != nil {
		return fmt.Errorf("error while generating select options: %w", err)
	}

	repo.ApiID[key] = apiId
	repo.EndpointID[key] = endpointId
	return nil
}

func (repo EndpointRepo) genPageView(key EndpointRepoKey, templateName string, pageData object.APIEndpointPage) error {
	buf := bytes.NewBuffer(nil)
	gz, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
	err := repo.apiRepo.View.ExecuteTemplate(gz, templateName, pageData)
	if err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return fmt.Errorf("error while compressing page: %w", err)
	}

	ctext := buf.Bytes()
	repo.PageView[key] = key.CacheKey(ctext)
	cmd := repo.apiRepo.Cache.Set(
		context.Background(), repo.PageView[key],
		ctext, global.CacheExpireDefault)
	return cmd.Err()
}

func (repo EndpointRepo) genSelectOptionsScript(key EndpointRepoKey, pageData object.APIEndpointPage, minify bool) error {
	if _, ok := repo.PageSelectOption[key.APIName()]; ok {
		return nil
	}

	originalScriptBuf := bytes.NewBuffer(nil)
	err := repo.apiRepo.View.ExecuteTemplate(originalScriptBuf, "options.gotmpl", pageData.SelectOpts)
	if err != nil {
		return err
	}

	if minify {
		minifiedScrptbuf := bytes.NewBuffer(nil)
		if err := global.Minifier().Minify("application/javascript", minifiedScrptbuf, originalScriptBuf); err != nil {
			return fmt.Errorf("error while minifying options: %w", err)
		}
		repo.PageSelectOption[key.APIName()] = key.CacheKey(minifiedScrptbuf.Bytes())
		cmd := repo.apiRepo.Cache.Set(
			context.Background(), repo.PageSelectOption[key.APIName()],
			minifiedScrptbuf.Bytes(), global.CacheExpireDefault)
		return cmd.Err()
	}

	repo.PageSelectOption[key.APIName()] = key.CacheKey(originalScriptBuf.Bytes())
	cmd := repo.apiRepo.Cache.Set(
		context.Background(), repo.PageSelectOption[key.APIName()],
		originalScriptBuf.Bytes(), global.CacheExpireDefault)
	return cmd.Err()
}

func (repo EndpointRepo) getCache(ctx context.Context, key EndpointRepoKey, ckey string) ([]byte, error) {
	page, err := repo.apiRepo.Cache.Get(context.Background(), ckey).Bytes()
	switch err {
	default:
		return nil, fmt.Errorf("error while read cache: %w", err)
	case redis.Nil:
		if err := repo.RegisterEndpointsPageView(
			key.APIName(), repo.ApiID[key],
			key.EndpointName(), repo.EndpointID[key],
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
	ckey, ok := repo.PageSelectOption[key.APIName()]
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

func (repo EndpointRepo) Do(req api.Request, handler client.Handler) (api.Response, error) {
	httpReq, err := req.ToHttpRequest()
	if err != nil {
		return nil, err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return handler.Parse(httpResp)
}

func postEndpoints(repo EndpointRepo, pageform pf.PageForm, userInfo tokenmaker.Payload,
	w http.ResponseWriter, httpReq *http.Request) {

	userInfo, ok := httpReq.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError).
			WithDetails("user info missing")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if err := httpReq.ParseForm(); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while .ParseForm").
			WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	pageform, err := pageform.FormDecodeAndValidate(
		repo.apiRepo.FormDecoder, repo.val, httpReq.PostForm)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while .FormDecodeAndValidate").
			WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	apikey, _ := repo.apiRepo.Service.APIKey().Get(
		httpReq.Context(), &service.APIKeyGetRequest{
			Owner: userInfo.GetUserID(),
			ApiID: repo.ApiID[NewRepoMapKey(pageform.API(), pageform.Endpoint())],
		},
	)

	handler, err := client.HandlerRepo.Get(pageform.API(), pageform.Endpoint())
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while calling .PageFormHandlerRepo.Get method").
			WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	ckey, cache, err := handler.Handle(apikey.Key, userInfo.GetUserID(), pageform)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while calling .NewQueryFromPageFrom method").
			WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	_, _ = repo.apiRepo.Cache.JSONSet(ckey, ".", cache)
	_ = repo.apiRepo.Cache.Expire(httpReq.Context(), ckey, 10*time.Minute)

	aid := repo.ApiID[NewRepoMapKey(pageform.API(), pageform.Endpoint())]
	eid := repo.EndpointID[NewRepoMapKey(pageform.API(), pageform.Endpoint())]
	http.Redirect(w, httpReq,
		fmt.Sprintf("/v1/preview/%s?aid=%d&eid=%d", ckey, aid, eid),
		http.StatusSeeOther)
	return
}
