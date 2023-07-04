package router

import (
	"context"
	"fmt"
	"net/http"
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
)

func NewRouter(srvc service.Service, vw view.View, filesystem http.FileSystem,
	tmaker tokenmaker.TokenMaker, cmaker *cookiemaker.CookieMaker) *chi.Mux {
	errHandlerRepo, _ := errorhandler.NewErrorHandlerRepo(vw.Template.Lookup("errorpage.gotmpl"))

	formDecoder := form.NewDecoder()
	formDecoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", vals[0])
	}, time.Time{})

	auth := auth.NewAuthRepo("v1", srvc, vw, tmaker, cmaker, formDecoder)
	apiRepo := api.NewAPIRepo("v1", srvc, vw, tmaker, cmaker, formDecoder)

	epRepo := apiRepo.EndpointRepo()
	epChan := make(chan *model.ListAllEndpointRow)
	go func(chan *model.ListAllEndpointRow) {
		apiRepo.Service.Endpoint().ListAll(context.Background(), 100, epChan)
	}(epChan)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(epRepo api.EndpointRepo, epChan chan *model.ListAllEndpointRow, wg *sync.WaitGroup) {
		for ep := range epChan {
			apiName, endpointName, templateName := ep.ApiName, ep.EndpointName, ep.TemplateName
			_ = epRepo.RegisterEndpointsPageView(apiName, endpointName, templateName)
		}
		wg.Done()
	}(epRepo, epChan, wg)

	bearerTokenMaker := middleware.BearerTokenMaker{
		AllowFromHTTPCookie: true,
		TokenMaker:          tmaker,
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/favicon.ico", http.NotFound)
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesystem)))

	r.Get(global.AppVar.Server.RoutePattern.Pages["login"], auth.GetLogin)
	r.Post(global.AppVar.Server.RoutePattern.Pages["login"], auth.PostLogin)

	r.Get(global.AppVar.Server.RoutePattern.Pages["sign-up"], auth.GetSignUp)
	r.Post(global.AppVar.Server.RoutePattern.Pages["sign-up"], auth.PostSignUp)

	r.Get(global.AppVar.Server.RoutePattern.Pages["logout"], auth.Logout)
	r.Get(global.AppVar.Server.RoutePattern.ErrorPages["unauthorized"], errHandlerRepo.Unauthorized)
	r.Get(global.AppVar.Server.RoutePattern.ErrorPages["badrequest"], errHandlerRepo.BadRequest)
	r.Get("/*", errHandlerRepo.SeeOther("/login"))

	r.Route(fmt.Sprintf("/%s", apiRepo.Version), func(r chi.Router) {
		r.Use(bearerTokenMaker.BearerAuthenticator)
		r.Get(global.AppVar.Server.RoutePattern.Pages["welcome"], apiRepo.GetWelcome)

		r.Get(global.AppVar.Server.RoutePattern.Pages["apikey"], apiRepo.GetAPIKey)
		r.Post(global.AppVar.Server.RoutePattern.Pages["apikey"], apiRepo.PostAPIKey)
		r.Delete(global.AppVar.Server.RoutePattern.Pages["apikey"]+"/{id}", apiRepo.DeleteAPIKey)

		r.Get(global.AppVar.Server.RoutePattern.Pages["change_password"], auth.GetChangePassword)
		r.Post(global.AppVar.Server.RoutePattern.Pages["change_password"], auth.PostChangPassword)

		r.Get(global.AppVar.Server.RoutePattern.Pages["admin"], apiRepo.GetAdmin)

		r.Route(
			global.AppVar.Server.RoutePattern.Pages["endpoints"],
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
		r.Get(global.AppVar.Server.RoutePattern.ErrorPages["forbidden"], errHandlerRepo.Forbidden)
		r.Get("/*", errHandlerRepo.NotFound)
	})

	return r
}
