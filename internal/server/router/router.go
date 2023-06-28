package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/api/v1"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/auth"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(srvc service.Service, vw view.View, filesystem http.FileSystem,
	tmaker tokenmaker.TokenMaker, cmaker *cookiemaker.CookieMaker) *chi.Mux {
	auth := auth.NewAuthRepo("v1", srvc, vw, tmaker, cmaker)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Get("/favicon.ico", http.NotFound)
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesystem)))
	r.Get("/login", auth.GetLogin)
	r.Post("/login", auth.PostLogin)

	r.Get("/sign-up", auth.GetSignUp)
	r.Post("/sign-up", auth.PostSignUp)

	r.Get("/logout", auth.Logout)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	r.Get("/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   cookiemaker.AUTH_COOKIE_KEY,
			MaxAge: -1,
		})
		w.WriteHeader(http.StatusUnauthorized)
		_ = vw.ExecuteTemplate(w, "errorpage.gotmpl", view.ErrorPage401)
	})

	bearerTokenMaker := middleware.BearerTokenMaker{
		AllowFromHTTPCookie: true,
		TokenMaker:          tmaker,
	}

	api := api.NewAPIRepo("v1", srvc, vw, tmaker, cmaker)

	eps := make(chan *model.ListAllEndpointRow)
	go func(chan *model.ListAllEndpointRow) {
		api.Service.Endpoint().ListAll(context.Background(), 100, eps)
	}(eps)

	r.Route(fmt.Sprintf("/%s", api.Version), func(r chi.Router) {
		r.Use(bearerTokenMaker.BearerAuthenticator)
		r.Get(global.AppVar.Server.RoutePattern.Pages["welcome"], api.GetWelcome)
		r.Get(global.AppVar.Server.RoutePattern.Pages["apikey"], api.GetAPIKey)
		r.Post(global.AppVar.Server.RoutePattern.Pages["apikey"], api.PostAPIKey)
		r.Delete(global.AppVar.Server.RoutePattern.Pages["apikey"]+"/{id}", api.DeleteAPIKey)
		r.Get(global.AppVar.Server.RoutePattern.Pages["change_password"], auth.GetChangePassword)
		r.Post(global.AppVar.Server.RoutePattern.Pages["change_password"], auth.PostChangPassword)
		r.Get(global.AppVar.Server.RoutePattern.Pages["admin"], api.GetAdmin)

		r.Route(
			global.AppVar.Server.RoutePattern.Pages["endpoints"],
			func(r chi.Router) {
				r.Get("/", api.GetEndpoints)
				for ep := range eps {
					apiName, endpointName, templateName := ep.ApiName, ep.EndpointName, ep.TemplateName

					r.Get(
						strings.TrimSuffix("/"+templateName, ".gotmpl"),
						func(w http.ResponseWriter, r *http.Request) {
							pageData := object.APIEndpointPage{
								Page: object.Page{
									HeadConent: view.NewHeadContent(),
									Title:      endpointName,
								},
								API:      apiName,
								Version:  global.AppVar.Server.APIVersion,
								Endpoint: endpointName,
							}
							w.WriteHeader(http.StatusOK)
							err := api.View.ExecuteTemplate(w, templateName, pageData)
							if err != nil {
								fmt.Println(err)
							}
						},
					)
				}

			})
		r.Get("/forbidden", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			api.View.ExecuteTemplate(w, "errorpage.gotmpl", view.ErrorPage403)
		})

		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			api.View.ExecuteTemplate(w, "errorpage.gotmpl", view.ErrorPage404)
		})
	})

	return r
}
