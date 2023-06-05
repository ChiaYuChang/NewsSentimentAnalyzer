package router

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Router struct {
	*chi.Mux
	*template.Template
	service.Service
}

func NewRouter(srvc service.Service) *chi.Mux {
	var err error
	r := Router{}
	r.Template, err = template.ParseGlob(global.AppVar.Server.TemplatePath)
	if err != nil {
		panic("error while parsing template")
	}

	r.Mux = chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/login", func(w http.ResponseWriter, req *http.Request) {
		r.Template.ExecuteTemplate(w, "login.gotmpl", object.LoginPage{
			Page: object.Page{
				HeadConent: view.NewHeadContent(),
				Title:      "Login",
			},
			ShowUsernameNotFountAlert: false,
			ShowPasswordMismatchAlert: false,
		})
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/login", func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			w.Write([]byte("Client error"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pform := req.PostForm
		email, password := pform.Get("email"), pform.Get("password")
		if err := r.Service.User().
			Login(context.Background(), email, password); err != nil {
			data := object.LoginPage{
				Page: object.Page{
					HeadConent: view.NewHeadContent(),
					Title:      "Login",
				}}

			if errors.Is(err, pgx.ErrNoRows) {
				data.ShowUsernameNotFountAlert = true
			} else if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				data.ShowPasswordMismatchAlert = true
			} else {
				w.Write([]byte(fmt.Sprintf("Unknown error: %s", err)))
				w.WriteHeader(http.StatusBadRequest)
			}
			r.Template.ExecuteTemplate(w, "login.gotmpl", data)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Redirect(w, req, "/v1/welcome", http.StatusSeeOther)
		}
		return
	})

	r.Route("/v1", func(r chi.Router) {

	})

	return r
}
