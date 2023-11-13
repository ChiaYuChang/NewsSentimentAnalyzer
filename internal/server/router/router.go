package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/api/v1"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/auth"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"

	cm "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/cookieMaker"
	eh "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/errorHandler"
	tm "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	// init server side
	pf "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GoogleCSE"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/NEWSDATA"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"

	// init client side
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/GNews"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/GoogleCSE"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/NEWSDATA"
	_ "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/newsapi"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewRouter(srvc service.Service, rds *redis.Client, vw view.View,
	tmaker tm.TokenMaker, cmaker *cm.CookieMaker) *chi.Mux {
	errHandlerRepo, err := eh.NewErrorHandlerRepo(vw.Template.Lookup("errorpage.gotmpl"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while errorhandler.NewErrorHandlerRepo: %s", err)
		os.Exit(1)
	}

	auth := auth.NewAuthRepo(
		viper.GetString("APP_API_VERSION"),
		srvc, vw, tmaker, validator.Validate, pf.Decoder, pf.Modifier)

	apiRepo := api.NewAPIRepo(
		viper.GetString("APP_API_VERSION"),
		srvc, vw, tmaker, validator.Validate, pf.Decoder, pf.Modifier)

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
			if err := epRepo.RegisterEndpointsPageView(
				apiName, apiID, endpointName, templateName); err != nil {
				global.Logger.
					Error().
					Str("api", apiName).
					Str("endpoint", endpointName).
					Str("status", "failed").
					Err(err).Msg("error while registering endpoint")
			} else {
				global.Logger.
					Info().
					Str("api", apiName).
					Str("endpoint", endpointName).
					Str("status", "ok").
					Msg("add endpoint")
			}
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
		r.Patch(rp.Page["change-password"], auth.PatchChangePassword)

		r.Get(rp.Page["admin"], apiRepo.GetAdmin)

		r.Get(rp.Page["job"], apiRepo.GetJob)
		r.Post(rp.Page["job"], apiRepo.PostJob)
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
						global.Logger.
							Error().
							Err(err).
							Msgf("error while .GetAPIEndpoints with key: %s", key.String())
					}

					if epPostHandlerFun, err := epRepo.PostAPIEndpoints(key); err == nil {
						r.Post("/"+key.String(), epPostHandlerFun)
					} else {
						global.Logger.
							Error().
							Err(err).
							Msgf("error while .PostAPIEndpoints with key: %s", key.String())
					}
				}
				r.Get("/*", errHandlerRepo.NotFound)
			})

		r.Get("/result", apiRepo.GetResultSelector)

		r.Get(rp.ErrorPage["forbidden"], errHandlerRepo.Forbidden)
		r.Get("/*", errHandlerRepo.NotFound)
	})

	return r
}
