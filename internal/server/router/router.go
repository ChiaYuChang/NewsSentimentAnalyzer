package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/api/v1"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/auth"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	errorhandler "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/errorHandler"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/GNews"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/newsapi"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/form"
	"github.com/spf13/viper"
)

func NewRouter(srvc service.Service, vw view.View, filesystem http.FileSystem,
	tmaker tokenmaker.TokenMaker, cmaker *cookiemaker.CookieMaker) *chi.Mux {
	errHandlerRepo, err := errorhandler.NewErrorHandlerRepo(vw.Template.Lookup("errorpage.gotmpl"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while errorhandler.NewErrorHandlerRepo: %s", err)
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()
	formDecoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse(time.DateOnly, vals[0])
	}, time.Time{})

	auth := auth.NewAuthRepo(viper.GetString("APP_API_VERSION"), srvc, vw, tmaker, cmaker, formDecoder)
	apiRepo := api.NewAPIRepo(viper.GetString("APP_API_VERSION"), srvc, vw, tmaker, cmaker, formDecoder)

	epRepo := apiRepo.EndpointRepo()
	epChan := make(chan *model.ListAllEndpointRow)
	go func(chan *model.ListAllEndpointRow) {
		apiRepo.Service.Endpoint().ListAll(context.Background(), 100, epChan)
	}(epChan)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(epRepo api.EndpointRepo, epChan chan *model.ListAllEndpointRow, wg *sync.WaitGroup) {
		for ep := range epChan {
			apiName, apiID, endpointName, templateName := ep.ApiName, ep.ApiID, ep.EndpointName, ep.TemplateName
			_ = epRepo.RegisterEndpointsPageView(apiName, apiID, endpointName, templateName)
		}
		wg.Done()
	}(epRepo, epChan, wg)

	bearerTokenMaker := middleware.BearerTokenMaker{
		AllowFromHTTPCookie: true,
		TokenMaker:          tmaker,
	}

	rp := global.AppVar.App.RoutePattern
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/favicon.ico", http.NotFound)
	r.Handle(rp.StaticPage, http.StripPrefix("/static", http.FileServer(filesystem)))

	r.Get(rp.Page["sign-in"], auth.GetSignIn)
	r.Post(rp.Page["sign-in"], auth.PostSignIn)

	r.Get(rp.Page["sign-up"], auth.GetSignUp)
	r.Post(rp.Page["sign-up"], auth.PostSignUp)

	r.Get(rp.Page["sign-out"], auth.GetSignOut)
	r.Get(rp.ErrorPage["unauthorized"], errHandlerRepo.Unauthorized)
	r.Get(rp.ErrorPage["bad-request"], errHandlerRepo.BadRequest)
	r.Get("/*", errHandlerRepo.SeeOther("/sign-in"))

	r.Route(fmt.Sprintf("/%s", apiRepo.Version), func(r chi.Router) {
		r.Use(bearerTokenMaker.BearerAuthenticator)
		r.Get(rp.Page["welcome"], apiRepo.GetWelcome)

		r.Get(rp.Page["apikey"], apiRepo.GetAPIKey)
		r.Post(rp.Page["apikey"], apiRepo.PostAPIKey)
		r.Delete(rp.Page["apikey"]+"/{id}", apiRepo.DeleteAPIKey)

		r.Get(rp.Page["change-password"], auth.GetChangePassword)
		r.Post(rp.Page["change-password"], auth.PostChangPassword)

		r.Get(rp.Page["admin"], apiRepo.GetAdmin)

		r.Route(
			rp.Page["endpoints"],
			func(r chi.Router) {
				r.Get("/", apiRepo.GetEndpoints)
				wg.Wait()
				for key := range epRepo.PageView {
					if epGetHandlerFun, err := epRepo.GetAPIEndpoints(key); err == nil {
						r.Get("/"+key.String(), epGetHandlerFun)
					} else {
						fmt.Println(err)
					}

					if epPostHandlerFun, err := epRepo.PostAPIEndpoints(key); err == nil {
						r.Post("/"+key.String(), epPostHandlerFun)
					} else {
						fmt.Println(err)
					}
				}
				r.Get("/*", errHandlerRepo.NotFound)
			})
		r.Get(rp.ErrorPage["forbidden"], errHandlerRepo.Forbidden)
		r.Get("/*", errHandlerRepo.NotFound)
	})

	return r
}
