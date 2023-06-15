package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
		Name: payload.GetUsername(),
		Role: payload.GetRole().String(),
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
		obj := &object.APIKey{
			ID:   a.ApiID,
			Name: a.Name,
			Key:  a.Key,
			Icon: "/static/image/logo_API_Default.svg",
		}

		if a.Image.Valid {
			obj.Icon = a.Image.String
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

func (repo APIRepo) PostAPIKey(w http.ResponseWriter, req *http.Request) {
	// userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	// if !ok {
	// 	ecErr := ec.MustGetEcErr(ec.ECServerError)
	// 	ecErr.WithDetails("user information not found")
	// 	w.WriteHeader(ecErr.HttpStatusCode)
	// 	w.Write(ecErr.MustToJson())
	// 	return
	// }

	body, _ := io.ReadAll(req.Body)
	m := make(map[string]string)
	_ = json.Unmarshal(body, &m)
	fmt.Println("map :", m)
	fmt.Println("Body:", string(body))

	bs, _ := json.Marshal(m)
	// w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	// repo.Service.APIKey().Get(req.Context(), &service.APIKeyGetRequest{
	// 	Owner: userInfo.GetUserID(), A
	// })

	// repo.Service.APIKey().Create()

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
	err = repo.View.ExecuteTemplate(w, "endpoint.gotmpl", pageData)
	if err != nil {
		fmt.Println(err)
	}
	return
}
