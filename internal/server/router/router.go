package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	r.Route("/v1", func(r chi.Router) {

	})

	return r
}
