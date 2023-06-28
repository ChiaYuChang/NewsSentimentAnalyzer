package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
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
	"github.com/jackc/pgx/v5"
)

type APIRepo struct {
	Version     string
	Service     service.Service
	View        view.View
	TokenMaker  tokenmaker.TokenMaker
	CookieMaker *cookiemaker.CookieMaker
	FormDecoder *form.Decoder
}

func NewAPIRepo(ver string, srvc service.Service,
	view view.View, tokenmaker tokenmaker.TokenMaker,
	cookiemaker *cookiemaker.CookieMaker) APIRepo {
	return APIRepo{
		Version:     ver,
		Service:     srvc,
		View:        view,
		TokenMaker:  tokenmaker,
		CookieMaker: cookiemaker,
		FormDecoder: form.NewDecoder(),
	}
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
			HeadConent: view.NewHeadContent(),
			Title:      "Welcome",
		},
		Name:             payload.GetUsername(),
		Role:             payload.GetRole().String(),
		PageEndpoint:     strings.TrimLeft(global.AppVar.Server.RoutePattern.Pages["endpoints"], "/"),
		PageChangePWD:    strings.TrimLeft(global.AppVar.Server.RoutePattern.Pages["change_password"], "/"),
		PageManageAPIKey: strings.TrimLeft(global.AppVar.Server.RoutePattern.Pages["apikey"], "/"),
		PageSeeResult:    "#",
		PageAdmin:        strings.TrimLeft(global.AppVar.Server.RoutePattern.Pages["admin"], "/"),
		PageLogout:       global.AppVar.Server.RoutePattern.Pages["logout"],
	}

	_ = repo.View.ExecuteTemplate(w, "welcome.gotmpl", pageData)
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
			HeadConent: view.NewHeadContent(),
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
			Icon: "/static/image/logo_API_Default.svg",
		}

		if a.Image.Valid {
			obj.Icon = a.Icon.String
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
		fmt.Println(err)
		return
	}

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
			HeadConent: view.NewHeadContent(),
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
		http.Redirect(w, req, "/v1/forbidden", http.StatusSeeOther)
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
