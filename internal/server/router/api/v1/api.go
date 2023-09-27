package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	hc := view.NewHeadContent()
	hc.Script.NewHTMLElement().AddPair("src", "/static/js/job_funcs.js")
	hc.Script.NewHTMLElement().AddPair("src", "//cdnjs.cloudflare.com/ajax/libs/list.js/2.3.1/list.min.js")

	pageData := object.APIResultPage{
		Page: object.Page{
			HeadConent: hc,
			Title:      "job",
		},
		PageSize: 10,
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

	pageData.NCreated = jobSummary.JobGroup[model.JobStatusCreated].NJob
	pageData.NRunning = jobSummary.JobGroup[model.JobStatusRunning].NJob
	pageData.NDone = jobSummary.JobGroup[model.JobStatusDone].NJob
	pageData.NFailed = jobSummary.JobGroup[model.JobStatusFailure].NJob
	pageData.NCanceled = jobSummary.JobGroup[model.JobStatusCanceled].NJob
	pageData.TotalJobs = jobSummary.TotalJob

	w.WriteHeader(http.StatusOK)
	_ = repo.View.ExecuteTemplate(w, "result.gotmpl", pageData)
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

	time.Sleep(2 * time.Second)
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

	if pager.FromJId < 0 || pager.ToJId < 0 || pager.JStatusStr == "" || pager.ParseJIds() != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
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
