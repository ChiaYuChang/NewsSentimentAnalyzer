package middleware_test

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

var opt global.TokenMakerOption

func init() {
	secret := make([]byte, 256)
	_, _ = rand.Read(secret)
	opt = global.TokenMakerOption{
		Secret:      secret,
		ExpireAfter: tokenmaker.DEFAULT_EXPIRE_AFTER,
		ValidAfter:  tokenmaker.DEFAULT_VALID_AFTER,
		SignMethod:  tokenmaker.DEFAULT_JWT_SIGN_METHOD,
	}
}

func TestBearerAuthenticator(t *testing.T) {
	username := "username"
	password := "password"
	role := tokenmaker.RUser

	maker := middleware.NewJWTTokenMaker(opt, tokenmaker.WithIssuer("[[:Issuer:]]"))

	r := chi.NewRouter()
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("login page"))
		return
	})
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		usr := r.FormValue("username")
		pwd := r.FormValue("password")

		ecErr := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
		if usr == username && pwd == password {
			ecErr = ec.MustGetErr(ec.Success).(*ec.Error)
			bearer, _ := maker.MakeToken(username, role)
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		body, _ := ecErr.ToJson()
		w.Write(body)
		return
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(maker.BearerAuthenticator)
		r.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(middleware.CtxUserInfo).(tokenmaker.UserInfo)
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(fmt.Sprintf("welcome! %s (%s)", user.UserName, user.Role)))
			return
		})
	})

	srv := httptest.NewTLSServer(r)
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	t.Run(
		"Get login page",
		func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", srv.URL, "login"), nil)
			resp, err := cli.Do(req)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer resp.Body.Close()
			require.Equal(t, "login page", string(body))
		},
	)

	t.Run(
		"Post login page correct Auth",
		func(t *testing.T) {
			vals := url.Values{}
			vals.Add("username", username)
			vals.Add("password", password)
			req, _ := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/%s", srv.URL, "login"),
				strings.NewReader(vals.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := cli.Do(req)
			require.NoError(t, err)

			_, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.NotEmpty(t, resp.Header.Get("Authorization"))

			req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", srv.URL, "user/welcome"), nil)
			req.Header.Add("Authorization", resp.Header.Get("Authorization"))

			resp, err = cli.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		},
	)

	t.Run(
		"Post login page wrong Auth",
		func(t *testing.T) {
			vals := url.Values{}
			vals.Add("username", username)
			vals.Add("password", "123456")
			req, _ := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/%s", srv.URL, "login"),
				strings.NewReader(vals.Encode()))

			resp, err := cli.Do(req)
			require.NoError(t, err)

			_, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			require.Empty(t, resp.Header.Get("Authorization"))
		},
	)

	t.Run(
		"Get user page without bearer token",
		func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", srv.URL, "user/welcome"), nil)
			resp, err := cli.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		},
	)

	t.Run(
		"Get user page with error bearer token",
		func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", srv.URL, "user/welcome"), nil)
			req.Header.Add("Authorization", fmt.Sprintf(fmt.Sprintf("Bearer %s", "ERROR-TOKEN")))

			resp, err := cli.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		},
	)

	t.Run(
		"Get user page with ok bearer token",
		func(t *testing.T) {
			bearer, err := maker.MakeToken(username, role)
			require.NoError(t, err)

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", srv.URL, "user/welcome"), nil)
			req.Header.Add("Authorization", fmt.Sprintf(fmt.Sprintf("Bearer %s", bearer)))

			resp, err := cli.Do(req)

			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.NotEmpty(t, resp.Header.Get("Authorization"))
			require.Contains(t, string(body), "welcome!")
		},
	)
}
