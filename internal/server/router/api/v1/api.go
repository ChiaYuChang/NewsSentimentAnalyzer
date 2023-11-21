package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	cm "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/cookieMaker"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/cache"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/go-playground/mold/v4"
	val "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type APIRepo struct {
	Version      string
	Service      service.Service
	View         view.View
	Cache        *cache.RedsiStore
	TokenMaker   tokenmaker.TokenMaker
	CookieMaker  *cm.CookieMaker
	Validator    *val.Validate
	FormDecoder  *form.Decoder
	FormModifier *mold.Transformer
}

func NewAPIRepo(
	ver string, srvc service.Service, view view.View, cache *cache.RedsiStore,
	tokenmaker tokenmaker.TokenMaker, cookiemaker *cm.CookieMaker,
	validator *val.Validate, decoder *form.Decoder, modifier *mold.Transformer) APIRepo {
	return APIRepo{
		Version:      ver,
		Service:      srvc,
		View:         view,
		Cache:        cache,
		TokenMaker:   tokenmaker,
		CookieMaker:  cookiemaker,
		Validator:    validator,
		FormDecoder:  decoder,
		FormModifier: modifier,
	}
}

func (repo APIRepo) HealthCheck(w http.ResponseWriter, req *http.Request) {
	j, _ := json.Marshal(map[string]any{
		"status code": http.StatusOK,
		"status":      "OK",
		"message":     "News Sentiment Analyzer (nsa)",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(j))
}

func (repo APIRepo) GetWelcome(w http.ResponseWriter, req *http.Request) {
	payload, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	pageData := object.WelcomePage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "Welcome",
		},
		Name:             payload.GetUsername(),
		Role:             payload.GetRole().String(),
		PageEndpoint:     strings.TrimLeft(global.AppVar.App.RoutePattern.Page["endpoints"], "/"),
		PageChangePWD:    strings.TrimLeft(global.AppVar.App.RoutePattern.Page["change-password"], "/"),
		PageManageAPIKey: strings.TrimLeft(global.AppVar.App.RoutePattern.Page["apikey"], "/"),
		PageSeeResult:    strings.TrimLeft(global.AppVar.App.RoutePattern.Page["job"], "/"),
		PageAdmin:        strings.TrimLeft(global.AppVar.App.RoutePattern.Page["admin"], "/"),
		PageSignOut:      global.AppVar.App.RoutePattern.Page["sign-out"],
	}

	if err := repo.View.ExecuteTemplate(w, "welcome.gotmpl", pageData); err != nil {
		global.Logger.
			Err(err).
			Msg("error while ExecuteTemplate welcome.gotmpl")
	}
	return
}

func (repo APIRepo) GetAPIKey(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	apis, err := repo.Service.API().List(req.Context(), 100)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}
	apiOpts := object.APIOption{
		Source:   map[int16]string{},
		Analyzer: map[int16]string{},
	}
	for _, api := range apis {
		if api.Type == model.ApiTypeSource {
			apiOpts.Source[api.ID] = api.Name
		}
		if api.Type == model.ApiTypeLanguageModel {
			apiOpts.Analyzer[api.ID] = api.Name
		}
	}

	apikey, err := repo.Service.APIKey().List(req.Context(), userInfo.GetUserID())
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	pageData := object.APIKeyPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "API Key",
		},
		APIVersion:   repo.Version,
		APIOption:    apiOpts,
		NewsAPIs:     []*object.APIKey{},
		AnalyzerAPIs: []*object.APIKey{},
	}

	for _, a := range apikey {
		obj := &object.APIKey{
			ID:   a.ApiID,
			Name: a.Name,
			Key:  a.Key,
			Icon: a.Icon,
		}

		switch a.Type {
		case model.ApiTypeSource:
			pageData.NewsAPIs = append(pageData.NewsAPIs, obj)
		case model.ApiTypeLanguageModel:
			pageData.AnalyzerAPIs = append(pageData.AnalyzerAPIs, obj)
		}
	}

	w.WriteHeader(http.StatusOK)
	err = repo.View.ExecuteTemplate(w, "apikey.gotmpl", pageData)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (repo APIRepo) DeleteAPIKey(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	apiIDRaw := chi.URLParam(req, "id")
	apiID, err := convert.StrTo(apiIDRaw).Int()
	if err != nil {
		ecErr := *ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	global.Logger.Info().
		Str("name", userInfo.GetUsername()).
		Int("api_id", apiID).
		Str("user_role", userInfo.GetRole().String()).
		Str("path", req.URL.Path)

	_, err = repo.Service.APIKey().Delete(req.Context(), &service.APIKeyDeleteRequest{
		Owner: userInfo.GetUserID(), ApiID: int16(apiID),
	})
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return
}

func (repo APIRepo) PostAPIKey(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if err := req.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}

	var apikey pageform.APIKeyPost
	if err := repo.FormDecoder.Decode(&apikey, req.PostForm); err != nil {
		return
	}

	if err := repo.FormModifier.Struct(req.Context(), &apikey); err != nil {
		return
	}

	if err := repo.Validator.StructCtx(req.Context(), &apikey); err != nil {
		return
	}

	apikey.Key = strings.TrimSpace(apikey.Key)
	if err := validator.Validate.Struct(apikey); err != nil {
		ecErr := *ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if results, err := repo.Service.APIKey().CreateOrUpdate(req.Context(),
		&service.APIKeyCreateOrUpdateRequest{
			Owner: userInfo.GetUserID(),
			ApiID: apikey.ApiID,
			Key:   apikey.Key,
		}); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
	} else {
		global.Logger.Info().
			Int16("api_id", apikey.ApiID).
			Int64("N", results.N).
			Msg("API key successfully updated/created")
	}

	http.Redirect(w, req, req.URL.Path, http.StatusSeeOther)
	return
}

func (repo APIRepo) GetEndpoints(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	eps, err := repo.Service.Endpoint().
		ListEndpointByOwner(req.Context(), userInfo.GetUserID())

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	pageData := object.APIEndpointFromDBModel(
		object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "Endpoints",
		},
		repo.Version,
		eps,
	)

	w.WriteHeader(http.StatusOK)
	if err = repo.View.ExecuteTemplate(w, "endpoint.gotmpl", pageData); err != nil {
		global.Logger.
			Error().
			Err(err).
			Msg("error while ExecuteTemplate endpoint.gotmpl")
	}
	return
}

func (repo APIRepo) GetAdmin(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if userInfo.GetRole() != tokenmaker.RAdmin {
		http.Redirect(w, req, "forbidden", http.StatusSeeOther)
	}

	pageData := object.APIAdminPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "admin",
		},
	}
	w.WriteHeader(http.StatusOK)
	if err := repo.View.ExecuteTemplate(w, "admin.gotmpl", pageData); err != nil {
		global.Logger.
			Error().
			Err(err).
			Msg("error executing template admin.gotmpl")
	}
	return
}

func (repo APIRepo) GetJob(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	pageData := object.APIResultPage{
		Page: object.Page{
			HeadConent: view.JobPageHeadContent(),
			Title:      "job",
		},
		NJobs:       map[string]int{},
		TotalJobKey: "all",
		PageSize:    15,
	}

	jobSummary, err := repo.Service.Job().Count(req.Context(), userInfo.GetUserID())
	if err != nil {
		if err == pgx.ErrNoRows {
			jobSummary = &model.CountUserJobTxResult{
				JobGroup:  make(map[model.JobStatus]model.JobGroup),
				TotalJob:  0,
				LastJobId: 0,
			}
		} else {
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			ecErr.WithDetails(err.Error())
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}
	}

	for k, v := range jobSummary.JobGroup {
		pageData.NJobs[string(k)] = v.NJob
	}
	pageData.NJobs[pageData.TotalJobKey] = jobSummary.TotalJob

	w.WriteHeader(http.StatusOK)
	if err = repo.View.ExecuteTemplate(w, "result.gotmpl", pageData); err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error executing template result.gotmpl")
	}
}

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func (repo APIRepo) PostJob(w http.ResponseWriter, req *http.Request) {
	global.Logger.Info().Msg("Call Post Job API")

	// time.Sleep(2 * time.Second)
	w.Header().Set("Content-Type", "application/json")
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if err := req.ParseForm(); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	var pager pageform.JobPager
	if err := repo.FormDecoder.Decode(&pager, req.Form); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if err := repo.FormModifier.Struct(req.Context(), &pager); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if err := repo.Validator.StructCtx(req.Context(), &pager); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	global.Logger.Debug().
		Str("owner", userInfo.GetUserID().String()).
		Str("jids", pager.JIdsStr).
		Int32("fjid", pager.FromJId).
		Int32("tjid", pager.ToJId).
		Str("status", pager.JStatusStr).
		Int("page", pager.Page).
		Msg("Get update query")

	rows, err := repo.Service.Job().Get(req.Context(), userInfo.GetUserID(),
		pager.JIds, pager.FromJId, pager.ToJId, pager.JStatusStr, 15, pager.Page)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		global.Logger.Debug().Msg(ecErr.Error())
		return
	}

	global.Logger.Debug().Int("N", len(rows)).Msg("get n rows")
	job := make([]object.Job, len(rows))
	for i, r := range rows {
		job[i] = object.Job{
			Id:        r.ID,
			Status:    object.StatusToClass(r.Status),
			NewsSrc:   r.NewsSrc,
			Analyzer:  r.Analyzer,
			CreatedAt: r.CreatedAt.Time.UTC().Format(time.DateTime),
			UpdatedAt: r.UpdatedAt.Time.UTC().Format(time.DateTime),
		}
	}

	jsn, err := json.Marshal(job)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsn)
	return
}

func (repo APIRepo) GetJobDetail(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	_ = req.ParseForm()
	jIdStr := chi.URLParam(req, "jId")
	jId, err := convert.StrTo(jIdStr).Int()
	w.Header().Set("Content-Type", "application/json")

	if jId <= 0 || err != nil {
		err := ec.MustGetEcErr(ec.ECBadRequest)
		err.WithDetails("jid not found")
		w.WriteHeader(err.HttpStatusCode)
		w.Write(err.MustToJson())
		return
	}

	job, err := repo.Service.Job().GetDetails(req.Context(), &service.JobGetByJobIdRequest{
		Owner: userInfo.GetUserID(),
		Id:    int32(jId),
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			ecErr := ec.MustGetEcErr(ec.ECForbidden)
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	jsn, _ := json.Marshal(object.NewJobDetails(userInfo.GetUsername(), job))
	w.WriteHeader(http.StatusOK)
	w.Write(jsn)
	return
}

func (repo APIRepo) EndpointRepo() EndpointRepo {
	return NewEndpointRepo(repo, validator.Validate)
}

func (repo APIRepo) GetResultSelector(w http.ResponseWriter, req *http.Request) {
	pageData := object.ResultSecectorPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "Result Selector",
		},
	}

	err := repo.View.ExecuteTemplate(w, "result_selector.gotmpl", pageData)
	if err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error executing template result_selector.gotmpl")
	}
}

func (repo APIRepo) GetPreview(w http.ResponseWriter, req *http.Request) {
	pageData := object.ResultSecectorPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "Result Selector",
		},
	}

	_ = req.ParseForm()
	err := repo.View.ExecuteTemplate(w, "preview.gotmpl", pageData)
	if err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error executing template preivew.gotmpl")
	}
}

type PreviewResponse struct {
	Error *PreviewError     `json:"error,omitempty"`
	Items []api.NewsPreview `json:"items"`
}

type PreviewError struct {
	Code        int      `json:"code,omitempty"`
	Message     string   `json:"message,omitempty"`
	Detail      []string `json:"details,omitempty"`
	RedirectURL string   `json:"url,omitempty"`
}

func (repo APIRepo) GetFetchNextPage(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseForm()
	pcid := chi.URLParam(req, "pcid")

	prev, ecErr := repo.getFetchNextPage(pcid)
	for i := range prev {
		// remove content before marshaling
		prev[i].Content = ""
	}

	respObj := PreviewResponse{Items: prev}
	if ecErr != nil {
		respObj.Error = &PreviewError{
			Code:    ecErr.HttpStatusCode,
			Message: ecErr.Message,
			Detail:  ecErr.Details,
		}
		switch ecErr.HttpStatusCode {
		case http.StatusBadRequest:
			respObj.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["bad-request"]
		case http.StatusGone:
			respObj.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["gone"]
		case http.StatusInternalServerError:
			respObj.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"]
		}
	}

	bprev, _ := json.Marshal(respObj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bprev)
}

func (repo APIRepo) getFetchNextPage(pcid string) ([]api.NewsPreview, *ec.Error) {
	b, err := repo.Cache.JSONGet(pcid, ".query")
	if err != nil {
		var ecErr *ec.Error
		if err == redis.Nil {
			ecErr = ec.MustGetEcErr(ec.ECGone).
				WithDetails("cache expired")
		} else {
			ecErr = ec.MustGetEcErr(ec.ECServerError).
				WithDetails("error getting query cache").
				WithDetails(err.Error())
		}
		return nil, ecErr
	}

	// read cache.query
	var cq api.CacheQuery
	err = json.Unmarshal(b.([]byte), &cq)
	if err != nil {
		global.Logger.Debug().Err(err).Msg("error unmarshal cq")
	}

	handler, err := client.HandlerRepo.GetByCacheQuery(cq)
	if err != nil {
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while get handler from client.HandlerRepo").
			WithDetails(err.Error())
	}

	// rebuild query from cache
	req, err := handler.RequestFromCacheQuery(cq)
	if err != nil {
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while .RequestFromCacheQuery").
			WithDetails(err.Error())
	}

	global.Logger.Info().
		Str("uid", cq.UserId.String()).
		Str("salt", cq.Salt).
		Str("rawQuery", cq.RawQuery).
		Str("nextPage", cq.NextPage.String()).
		Msg("rebuild cache query ok")

	// do request
	resp, err := client.HandlerRepo.Do(req, handler)
	if err != nil {
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while .Do").
			WithDetails(err.Error())
	}

	// append prev to cache
	next, prev := resp.ToNewsItemList()
	if _, err := repo.Cache.JSONArrAppend(pcid, ".news_item", prev); err != nil {
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while appending preview items to cache").
			WithDetails("error append prev to cache").
			WithDetails(err.Error())
	}

	// set next page token to cache
	if _, err := repo.Cache.JSONSet(pcid, ".query.next_page", next); err != nil {
		global.Logger.Debug().Err(err).Msg("error ")
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while set next page token to cache").
			WithDetails("error append prev to cache").
			WithDetails(err.Error())
	}
	return prev, nil
}
