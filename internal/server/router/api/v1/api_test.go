package api_test

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	mock_model "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/mockdb"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model/testtool"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/api/v1"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/go-playground/form"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var ErrRedirect = errors.New("redirect")
var ErrMakerMakeToken = errors.New("token maker error")
var ErrDatabaseConn = errors.New("database connection error")

const VIEWS_PATH = "../../../../../views"

const AUTH_TOKEN = "[[::AUTH_TOKEN::]]"

var opt global.JWTOption

func init() {
	secretLen := 256
	secret := make([]byte, secretLen)
	_, _ = rand.Read(secret)
	opt = global.JWTOption{
		Secret:          secret,
		SecretLength:    secretLen,
		ExpireAfterHour: 3,
		ValidAfterHour:  0,
		SignMethod:      tokenmaker.DEFAULT_JWT_SIGN_METHOD,
	}
}

func TestGetWelcome(t *testing.T) {
	tmpl, err := view.ParseTemplates(VIEWS_PATH+"/template/*.gotmpl", nil)
	require.NoError(t, err)

	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrRedirect
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	tm := middleware.NewJWTTokenMaker(opt)
	tm.AllowFromHTTPCookie = true
	cm := cookiemaker.NewDefaultCookieMacker("localhost")

	user, _ := testtool.GenRdmUser()
	bearer, err := tm.TokenMaker.MakeToken(user.Email, user.ID, tokenmaker.ParseRole(user.Role))
	require.NoError(t, err)

	type testCase struct {
		Name         string
		SetupStore   func(t *testing.T) model.Store
		SetupServer  func(store model.Store) *chi.Mux
		SetupRequest func(t *testing.T, url string) *http.Request
		Check        func(t *testing.T, resp *http.Response)
	}

	tcs := []testCase{
		{
			Name: "Get Welcome Page",
			SetupStore: func(t *testing.T) model.Store {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				return store
			},
			SetupServer: func(store model.Store) *chi.Mux {
				srvc := service.NewService(store, validator.Validate)
				apiRepo := api.APIRepo{
					Version:     "v1",
					Service:     srvc,
					Template:    tmpl,
					TokenMaker:  tm,
					CookieMaker: cm,
					FormDecoder: form.NewDecoder(),
				}
				mux := chi.NewMux()
				mux.Use(tm.BearerAuthenticator)
				mux.Get(fmt.Sprintf("/%s/welcome", apiRepo.Version), apiRepo.GetWelcome)
				return mux
			},
			SetupRequest: func(t *testing.T, url string) *http.Request {
				req, err := http.NewRequest(http.MethodGet, url+"/v1/welcome", nil)
				require.NoError(t, err)
				req.AddCookie(&http.Cookie{
					Name:  cookiemaker.AUTH_COOKIE_KEY,
					Value: bearer,
					Path:  "/",
				})
				return req
			},
			Check: func(t *testing.T, resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Contains(t, string(body), fmt.Sprintf("<h1>Welcome %s</h1>", user.Email))
				require.Contains(t, string(body), "<title>Welcome</title>")
			},
		},
		{
			Name: "Get API Key Page",
			SetupStore: func(t *testing.T) model.Store {
				apikeyrow := []*model.ListAPIKeyRow{}
				for i := 0; i < 10; i++ {
					api, err := testtool.GenRdmAPI(int16(i))
					require.NoError(t, err)

					key, err := testtool.GenRdmAPIKey(user.ID, api.ID)
					require.NoError(t, err)

					apikeyrow = append(apikeyrow, &model.ListAPIKeyRow{
						ApiKeyID: pgtype.Int4{Int32: key.ID, Valid: true},
						Owner:    pgtype.Int4{Int32: user.ID, Valid: true},
						Key:      pgtype.Text{String: key.Key, Valid: true},
						ApiID:    int16(i),
						Type:     api.Type,
						Name:     api.Name,
					})
				}

				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.
					EXPECT().
					ListAPIKey(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(apikeyrow, nil)
				return store
			},
			SetupServer: func(store model.Store) *chi.Mux {
				srvc := service.NewService(store, validator.Validate)
				apiRepo := api.APIRepo{
					Version:     "v1",
					Service:     srvc,
					Template:    tmpl,
					TokenMaker:  tm,
					CookieMaker: cm,
					FormDecoder: form.NewDecoder(),
				}
				mux := chi.NewMux()
				mux.Use(tm.BearerAuthenticator)
				mux.Get(fmt.Sprintf("/%s/apikey", apiRepo.Version), apiRepo.GetAPIKey)
				return mux
			},
			SetupRequest: func(t *testing.T, url string) *http.Request {
				req, err := http.NewRequest(http.MethodGet, url+"/v1/apikey", nil)
				require.NoError(t, err)
				req.AddCookie(&http.Cookie{
					Name:  cookiemaker.AUTH_COOKIE_KEY,
					Value: bearer,
					Path:  "/",
				})
				return req
			},
			Check: func(t *testing.T, resp *http.Response) {
				body, err := io.ReadAll(resp.Body)
				t.Log(string(body))
				require.Equal(t, http.StatusOK, resp.StatusCode)
				require.NoError(t, err)
			},
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("Case %d-%s", i, tc.Name),
			func(t *testing.T) {
				mux := tc.SetupServer(tc.SetupStore(t))
				srv := httptest.NewTLSServer(mux)
				req := tc.SetupRequest(t, srv.URL)
				resp, err := cli.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()
				tc.Check(t, resp)
			},
		)
	}
}
