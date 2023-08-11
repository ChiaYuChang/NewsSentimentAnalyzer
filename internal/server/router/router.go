package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
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
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewRouter(srvc service.Service, rds *redis.Client, vw view.View,
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

	auth := auth.NewAuthRepo(viper.GetString("APP_API_VERSION"), srvc, vw, tmaker, formDecoder)
	apiRepo := api.NewAPIRepo(viper.GetString("APP_API_VERSION"), srvc, vw, tmaker, formDecoder)

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

	authRateLimiter := middleware.NewRateLimiter(
		global.AppVar.RateLimiter.Auth.N,
		global.AppVar.RateLimiter.Auth.RateLimitOption(),
	)

	apiRateLimiter := middleware.NewRateLimiter(
		global.AppVar.RateLimiter.API.N,
		global.AppVar.RateLimiter.API.RateLimitOption(),
	)

	qureyRateLimiter := middleware.NewRedisRateLimiter(
		rds,
		global.AppVar.RateLimiter.User.N,
		global.AppVar.RateLimiter.User.Per,
	)

	rp := global.AppVar.App.RoutePattern
	r := chi.NewRouter()
	// r.Use(chimiddleware.Logger)
	r.Use(middleware.NewZerologLogger(zerolog.InfoLevel))
	r.Use(chimiddleware.Recoverer)
	r.Use(authRateLimiter.RateLimit)

	r.Get("/favicon.ico", http.NotFound)
	r.Handle(rp.StaticPage, http.StripPrefix(
		"/static", http.FileServer(http.Dir(global.AppVar.App.StaticFile.Path))))

	r.Get(rp.Page["sign-in"], auth.GetSignIn)
	r.Post(rp.Page["sign-in"], auth.PostSignIn)

	r.Get(rp.Page["sign-up"], auth.GetSignUp)
	r.Post(rp.Page["sign-up"], auth.PostSignUp)

	r.Get(rp.Page["sign-out"], auth.GetSignOut)

	r.Get(rp.ErrorPage["unauthorized"], errHandlerRepo.Unauthorized)
	r.Get(rp.ErrorPage["bad-request"], errHandlerRepo.BadRequest)
	r.Get(rp.ErrorPage["too-many-request"], errHandlerRepo.TooManyRequests)
	r.Get("/*", errHandlerRepo.SeeOther("/sign-in"))

	r.Route(fmt.Sprintf("/%s", apiRepo.Version), func(r chi.Router) {
		r.Use(chimiddleware.Recoverer)
		r.Use(apiRateLimiter.RateLimit)
		r.Use(bearerTokenMaker.BearerAuthenticator)

		r.Get(rp.Page["welcome"], apiRepo.GetWelcome)

		r.Get(rp.Page["apikey"], apiRepo.GetAPIKey)
		r.Post(rp.Page["apikey"], apiRepo.PostAPIKey)
		r.Delete(rp.Page["apikey"]+"/{id}", apiRepo.DeleteAPIKey)

		r.Get(rp.Page["change-password"], auth.GetChangePassword)
		r.Post(rp.Page["change-password"], auth.PostChangPassword)

		r.Get(rp.Page["admin"], apiRepo.GetAdmin)

		r.Get(rp.Page["job"], apiRepo.GetJob)
		r.Get(rp.Page["job"]+"/{jId}", apiRepo.GetJobDetail)

		r.Route(
			rp.Page["endpoints"],
			func(r chi.Router) {
				r.Use(qureyRateLimiter.RateLimit)
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
