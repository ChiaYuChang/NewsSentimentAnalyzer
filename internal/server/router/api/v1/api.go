package api

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5"
)

type APIRepo struct {
	Version     string
	Service     service.Service
	Template    *template.Template
	TokenMaker  tokenmaker.TokenMaker
	CookieMaker *cookiemaker.CookieMaker
	FormDecoder *form.Decoder
}

func NewAPIRepo(ver string, srvc service.Service,
	tmpl *template.Template, tokenmaker tokenmaker.TokenMaker,
	cookiemaker *cookiemaker.CookieMaker) APIRepo {
	return APIRepo{
		Version:     ver,
		Service:     srvc,
		Template:    tmpl,
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
		Name: payload.GetUsername(),
		Role: payload.GetRole().String(),
	}

	_ = repo.Template.ExecuteTemplate(w, "welcome.gotmpl", pageData)
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
		NewsAPIs:     []*object.APIKey{},
		AnalyzerAPIs: []*object.APIKey{},
	}

	for _, a := range apikey {
		img := object.NewHTMLElement("img").
			AddPair("alt", a.Name)
		if true {
			img.AddPair("src", "/static/image/api_default.svg")
		}

		input := object.NewHTMLElement("input").
			AddPair("id", fmt.Sprintf("apikey-%s", strings.ToLower(a.Name))).
			AddPair("name", fmt.Sprintf("apikey[%d]", a.ApiID))

		if a.Key.Valid {
			input.AddPair("value", a.Key.String)
		}

		switch a.Type {
		case model.ApiTypeSource:
			pageData.NewsAPIs = append(pageData.NewsAPIs, &object.APIKey{Image: img, Input: input})
		case model.ApiTypeLanguageModel:
			pageData.AnalyzerAPIs = append(pageData.AnalyzerAPIs, &object.APIKey{Image: img, Input: input})
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = repo.Template.ExecuteTemplate(w, "apikey.gotmpl", pageData)
	return
}
