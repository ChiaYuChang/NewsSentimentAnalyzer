package router

import (
	"html/template"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/api/v1"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/auth"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(srvc service.Service, tmpl *template.Template, filesystem http.FileSystem,
	tokenmaker tokenmaker.TokenMaker, cookiemaker *cookiemaker.CookieMaker) *chi.Mux {
	auth := auth.NewAuthRepo("v1", srvc, tmpl, tokenmaker, cookiemaker)

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

	bearerTokenMaker := middleware.BearerTokenMaker{
		AllowFromHTTPCookie: true,
		TokenMaker:          tokenmaker,
	}

	api := api.NewAPIRepo("v1", srvc, tmpl, tokenmaker, cookiemaker)
	r.Route("/v1", func(r chi.Router) {
		r.Use(bearerTokenMaker.BearerAuthenticator)
		r.Get("/welcome", api.GetWelcome)
		r.Get("/apikey", api.GetAPIKey)
		r.Post("/apikey", api.PostAPIKey)
		r.Get("/change_password", auth.GetChangePassword)
		r.Post("/change_password", auth.PostChangPassword)
		r.Get("/endpoints", api.GetEndpoints)
	})

	return r
}
