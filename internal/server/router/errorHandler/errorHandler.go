package errorhandler

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/cache"
	"github.com/redis/go-redis/v9"
)

type ErrorPage string

type ErrorHandlerRepo struct {
	code2key           map[int]string
	code2page          map[int]object.ErrorPage
	cache              *cache.RedsiStore
	tmpl               *template.Template
	DefaultClientError int
	DefaultServerError int
}

func NewErrorHandlerRepo(tmpl *template.Template, store *cache.RedsiStore) (ErrorHandlerRepo, error) {
	if tmpl == nil {
		return ErrorHandlerRepo{}, errors.New("a nil templated are provided")
	}

	repo := ErrorHandlerRepo{
		code2key:           map[int]string{},
		tmpl:               tmpl,
		cache:              store,
		DefaultClientError: http.StatusBadRequest,
		DefaultServerError: http.StatusInternalServerError,
	}

	repo.code2page = map[int]object.ErrorPage{
		http.StatusBadRequest:          view.ErrorPage400,
		http.StatusUnauthorized:        view.ErrorPage401,
		http.StatusForbidden:           view.ErrorPage403,
		http.StatusNotFound:            view.ErrorPage404,
		http.StatusTooManyRequests:     view.ErrorPage429,
		http.StatusInternalServerError: view.ErrorPage500,
	}

	for code, page := range repo.code2page {
		if err := repo.RegisterErrorPage(code, page); err != nil {
			return repo, nil
		}
	}
	return repo, nil
}

func (repo *ErrorHandlerRepo) RegisterErrorPage(statusCode int, epage object.ErrorPage) error {
	buf := bytes.NewBuffer(nil)
	gz, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)

	if err := repo.tmpl.Execute(gz, epage); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}
	ctext := buf.Bytes()

	hasher := sha1.New()
	hasher.Write(ctext)
	key := fmt.Sprintf(
		"error-page-%d-%s",
		statusCode,
		base64.StdEncoding.EncodeToString(hasher.Sum(nil)))

	cmd := repo.cache.Set(context.TODO(), key, ctext, global.CacheExpireLong)
	if cmd.Err() != nil {
		return fmt.Errorf("error while caching error page: %w", cmd.Err())
	}
	repo.code2key[statusCode] = key

	global.Logger.
		Info().
		Str("redis-key", key).
		Time("expiration", time.Now().Add(global.CacheExpireLong)).
		Msg("add page cache to redis")

	return nil
}

func (repo ErrorHandlerRepo) fetchErrorPage(httpEC int, w http.ResponseWriter, r *http.Request) {
	key, ok := repo.code2key[httpEC]
	if !ok {
		if httpEC >= 400 && httpEC < 500 {
			// 400 Bad Request
			key = repo.code2key[repo.DefaultClientError]
		} else {
			// 500 Internal Server Error
			key = repo.code2key[repo.DefaultServerError]
		}
	}

	cmd := repo.cache.Get(context.TODO(), key)

	b, err := cmd.Bytes()
	if err == redis.Nil {
		repo.RegisterErrorPage(httpEC, repo.code2page[httpEC])
	}
	_ = repo.cache.ExpireGT(context.TODO(), key, global.CacheExpireLong)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(httpEC)
	w.Write(b)
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

// 426 error
func (repo ErrorHandlerRepo) TooManyRequests(w http.ResponseWriter, req *http.Request) {
	repo.fetchErrorPage(http.StatusTooManyRequests, w, req)
}

// 500 error
func (repo ErrorHandlerRepo) InternalServerError(w http.ResponseWriter, req *http.Request) {
	repo.fetchErrorPage(http.StatusInternalServerError, w, req)
}
