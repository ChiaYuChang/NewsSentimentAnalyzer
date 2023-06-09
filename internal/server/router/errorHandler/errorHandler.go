package errorhandler

import (
	"bytes"
	"html/template"
	"net/http"

	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
)

type ErrorPage string

type ErrorHandlerRepo struct {
	page               map[int][]byte
	tmpl               *template.Template
	DefaultClientError int
	DefaultServerError int
}

func NewErrorHandlerRepo(tmpl *template.Template) (ErrorHandlerRepo, error) {
	repo := ErrorHandlerRepo{
		page:               map[int][]byte{},
		tmpl:               tmpl,
		DefaultClientError: http.StatusBadRequest,
		DefaultServerError: http.StatusInternalServerError,
	}

	type epages struct {
		epage  object.ErrorPage
		status int
	}

	eps := []epages{
		{view.ErrorPage400, http.StatusBadRequest},
		{view.ErrorPage401, http.StatusUnauthorized},
		{view.ErrorPage403, http.StatusForbidden},
		{view.ErrorPage404, http.StatusNotFound},
		{view.ErrorPage500, http.StatusInternalServerError},
	}

	for _, ep := range eps {
		if err := repo.RegisterErrorPage(ep.status, ep.epage); err != nil {
			return repo, nil
		}
	}
	return repo, nil
}

func (repo *ErrorHandlerRepo) RegisterErrorPage(statusCode int, epage object.ErrorPage) error {
	buffer := bytes.NewBufferString("")
	if err := repo.tmpl.Execute(buffer, epage); err != nil {
		return err
	}
	repo.page[statusCode] = buffer.Bytes()
	return nil
}

func (repo ErrorHandlerRepo) fetchErrorPage(httpEC int, w http.ResponseWriter, r *http.Request) {
	page, ok := repo.page[httpEC]
	if !ok {
		if httpEC >= 400 && httpEC < 500 {
			// 400
			page = repo.page[repo.DefaultClientError]
		} else {
			// 500
			page = repo.page[repo.DefaultServerError]
		}
	}
	w.WriteHeader(httpEC)
	w.Write(page)
}

func (repo ErrorHandlerRepo) SeeOther(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// 400 error
func (repo ErrorHandlerRepo) BadRequest(w http.ResponseWriter, req *http.Request) {
	repo.fetchErrorPage(http.StatusBadRequest, w, req)
}

// 401 error
func (repo ErrorHandlerRepo) Unauthorized(w http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie(cookiemaker.AUTH_COOKIE_KEY); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   cookiemaker.AUTH_COOKIE_KEY,
			MaxAge: -1,
		})
	}
	repo.fetchErrorPage(http.StatusUnauthorized, w, req)
}

// 403 error
func (repo ErrorHandlerRepo) Forbidden(w http.ResponseWriter, req *http.Request) {
	repo.fetchErrorPage(http.StatusForbidden, w, req)
}

// 404 error
func (repo ErrorHandlerRepo) NotFound(w http.ResponseWriter, req *http.Request) {
	repo.fetchErrorPage(http.StatusNotFound, w, req)
}

// 500 error
func (repo ErrorHandlerRepo) InternalServerError(w http.ResponseWriter, req *http.Request) {
	repo.fetchErrorPage(http.StatusInternalServerError, w, req)
}
