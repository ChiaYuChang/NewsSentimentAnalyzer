package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type APIRepo struct {
	Version     string
	Service     service.Service
	View        view.View
	TokenMaker  tokenmaker.TokenMaker
	Validate    *val.Validate
	FormDecoder *form.Decoder
}

func NewAPIRepo(ver string, srvc service.Service, view view.View,
	tokenmaker tokenmaker.TokenMaker, decoder *form.Decoder) APIRepo {
	return APIRepo{
		Version:     ver,
		Service:     srvc,
		View:        view,
		TokenMaker:  tokenmaker,
		FormDecoder: decoder,
	}
}

func (repo APIRepo) HealthCheck(w http.ResponseWriter, req *http.Request) {
	m := make(map[string]string)
	m["status code"] = strconv.Itoa(http.StatusOK)
	m["status"] = "OK"
	m["message"] = "News Sentiment Analyzer (nsa)"

	j, _ := json.Marshal(m)
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
			HeadConent: view.SharedHeadContent,
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

	err := repo.View.ExecuteTemplate(w, "welcome.gotmpl", pageData)
	if err != nil {
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
			HeadConent: view.SharedHeadContent,
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

	fmt.Printf("name: %s(%d), api: %d path: %s\n",
		userInfo.GetUsername(),
		userInfo.GetUserID(),
		apiID, req.URL.Path)

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
		fmt.Printf("API ID: %d, N: %d\n", results.ApiKeyId, results.N)
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
			HeadConent: view.SharedHeadContent,
			Title:      "Endpoints",
		},
		repo.Version,
		eps,
	)

	w.WriteHeader(http.StatusOK)
	_ = repo.View.ExecuteTemplate(w, "endpoint.gotmpl", pageData)
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
			HeadConent: view.SharedHeadContent,
			Title:      "admin",
		},
	}
	w.WriteHeader(http.StatusOK)
	_ = repo.View.ExecuteTemplate(w, "admin.gotmpl", pageData)
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

	if err := req.ParseForm(); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECInvalidParams)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if jobs, err := repo.Service.Job().GetByOwner(req.Context(), &service.JobGetByOwnerRequest{
		Owner: userInfo.GetUserID(),
		Next:  0,
		N:     10,
	}); err != nil && err != pgx.ErrNoRows {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	} else {
		hc := view.NewHeadContent()
		hc.Script.NewHTMLElement().
			AddPair("src", "/static/js/wasm_exec.js")
		hc.Script.NewHTMLElement().
			AddPair("src", "/static/js/wasm_go.js")

		pageData := object.APIResultPage{
			Page: object.Page{
				HeadConent: hc,
				Title:      "job",
			},
		}
		pageData.SetJobs(jobs)

		w.WriteHeader(http.StatusOK)
		_ = repo.View.ExecuteTemplate(w, "result.gotmpl", pageData)

	}
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

	req.ParseForm()
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

	job, err := repo.Service.Job().GetByJobId(req.Context(), &service.JobGetByJobIdRequest{
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

	analyzerQuery := map[string]string{}
	json.Unmarshal(job.LlmQuery, &analyzerQuery)
	jsn, _ := json.Marshal(object.NewJobDetails(userInfo.GetUsername(), job))
	w.WriteHeader(http.StatusOK)
	w.Write(jsn)
	return
}

func (repo APIRepo) EndpointRepo() EndpointRepo {
	return NewEndpointRepo(repo, validator.Validate)
}
