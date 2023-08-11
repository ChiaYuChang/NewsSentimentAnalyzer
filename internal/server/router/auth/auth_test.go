package auth_test

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	mock_model "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model/mockdb"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model/testtool"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/auth"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	mock_tokenMaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker/mockTokenMaker"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	// "github.com/go-playground/validator/v10"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var ErrRedirect = errors.New("redirect")
var ErrMakerMakeToken = errors.New("token maker error")
var ErrDatabaseConn = errors.New("database connection error")

const VIEWS_PATH = "../../../../views"

const AUTH_TOKEN = "[[::AUTH_TOKEN::]]"

func init() {
	cm := cookiemaker.NewDefaultCookieMacker("localhost")
	cookiemaker.SetDefaultCookieMaker(cm)
}

func TestAuthLogin(t *testing.T) {
	vw, err := view.NewView(nil, VIEWS_PATH+"/template/*.gotmpl")
	require.NoError(t, err)
	user, _ := testtool.GenRdmUser()

	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrRedirect
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	t.Run(
		"Get login page",
		func(t *testing.T) {
			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Any()).
				Times(0)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			srvc := service.NewService(store, validator.Validate)

			authReup := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Get("/login", authReup.GetSignIn)

			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL
			resp, err := cli.Get(srvcUrl + "/login")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, validator.Validate.Var(string(body), "html"))
		},
	)

	t.Run(
		"OK",
		func(t *testing.T) {
			cpassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
			require.NoError(t, err)
			row := &model.GetUserAuthRow{
				ID:       user.ID,
				Email:    user.Email,
				Password: cpassword,
				Role:     user.Role,
			}

			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(row, nil)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(
					gomock.Eq(user.Email),
					gomock.Eq(user.ID),
					gomock.Eq(tokenmaker.ParseRole(string(user.Role)))).
				Times(1).
				Return(AUTH_TOKEN, nil)

			srvc := service.NewService(store, validator.Validate)

			authReup := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/login", authReup.PostSignIn)
			mux.Get("/v1/welcome", func(w http.ResponseWriter, req *http.Request) {
				bearer := "TOKEN"
				if err != nil {
					fmt.Fprint(w, "hi, your bearer token could not be found")
					return
				}
				fmt.Fprintf(w, "hi, your bearer token is: %s", bearer)
				w.WriteHeader(http.StatusOK)
			})

			srv := httptest.NewTLSServer(mux)

			srvcUrl := srv.URL

			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", string(user.Password))
			resp, err := cli.PostForm(srvcUrl+"/login", values)
			require.ErrorIs(t, err, ErrRedirect)
			require.Equal(t, http.StatusSeeOther, resp.StatusCode)
			for _, cookie := range resp.Cookies() {
				switch cookie.Name {
				case cookiemaker.AUTH_COOKIE_KEY:
					require.Equal(t, AUTH_TOKEN, cookie.Value)
				}
			}
		},
	)

	t.Run(
		"Wrong Password",
		func(t *testing.T) {
			cpassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
			require.NoError(t, err)

			row := &model.GetUserAuthRow{
				ID:       user.ID,
				Email:    user.Email,
				Password: cpassword,
				Role:     user.Role,
			}

			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(row, nil)
			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			srvc := service.NewService(store, validator.Validate)

			authReup := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/login", authReup.PostSignIn)
			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL

			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", "[[::Wrong Password::]]")
			resp, err := cli.PostForm(srvcUrl+"/login", values)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, validator.Validate.Var(string(body), "html"))
			require.Contains(t, string(body), "<span>Wrong password. Please try again.</span>")
		},
	)

	t.Run(
		"Username (email) not found",
		func(t *testing.T) {
			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(new(model.GetUserAuthRow), pgx.ErrNoRows)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			srvc := service.NewService(store, validator.Validate)

			authReup := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/login", authReup.PostSignIn)
			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL

			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", string(user.Password))
			resp, err := cli.PostForm(srvcUrl+"/login", values)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, validator.Validate.Var(string(body), "html"))
			require.Contains(t, string(body), "<span>Couldn’t find your Account</span>")
		},
	)

	t.Run(
		"Database error",
		func(t *testing.T) {
			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(new(model.GetUserAuthRow), ErrDatabaseConn)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			srvc := service.NewService(store, validator.Validate)

			authReup := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/login", authReup.PostSignIn)
			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL

			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", string(user.Password))
			resp, err := cli.PostForm(srvcUrl+"/login", values)
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, validator.Validate.Var(string(body), "json"))
		},
	)

	t.Run(
		"Token maker error",
		func(t *testing.T) {
			cpassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
			require.NoError(t, err)
			row := &model.GetUserAuthRow{
				ID:       user.ID,
				Email:    user.Email,
				Password: cpassword,
				Role:     user.Role,
			}

			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(row, nil)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(
					gomock.Eq(user.Email),
					gomock.Eq(user.ID),
					gomock.Eq(tokenmaker.ParseRole(string(user.Role)))).
				Times(1).
				Return("", ErrMakerMakeToken)

			srvc := service.NewService(store, validator.Validate)

			authReup := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/login", authReup.PostSignIn)

			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL

			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", string(user.Password))
			resp, err := cli.PostForm(srvcUrl+"/login", values)
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Contains(
				t, string(body),
				"error while making signature: token maker error")
			require.NoError(t, validator.Validate.Var(string(body), "json"))
		},
	)
}

func TestAuthSignUp(t *testing.T) {
	vw, err := view.NewView(nil, VIEWS_PATH+"/template/*.gotmpl")
	require.NoError(t, err)
	user, _ := testtool.GenRdmUser()
	user.Role = model.RoleUser

	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrRedirect
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	t.Run(
		"Get sign up page",
		func(t *testing.T) {
			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Any()).
				Times(0)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			srvc := service.NewService(store, validator.Validate)

			authReup := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Get("/signup", authReup.GetSignUp)

			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL
			resp, err := cli.Get(srvcUrl + "/signup")
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, validator.Validate.Var(string(body), "html"))
			require.Contains(t, string(body), "<title>Sign up</title>")
		},
	)

	t.Run(
		"OK",
		func(t *testing.T) {
			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(new(model.GetUserAuthRow), pgx.ErrNoRows)

			matcher, err := testtool.NewUserCreateReqMatcher(
				&service.UserCreateRequest{
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
					Password:  string(user.Password),
					Role:      string(user.Role),
				})
			require.NoError(t, err)
			store.
				EXPECT().
				CreateUser(gomock.Any(), matcher).
				Times(1).
				Return(user.ID, nil)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(
					gomock.Eq(user.Email),
					gomock.Eq(user.ID),
					gomock.Eq(tokenmaker.ParseRole(string(user.Role)))).
				Times(1).
				Return(AUTH_TOKEN, nil)

			srvc := service.NewService(store, validator.Validate)

			authRepo := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/signup", authRepo.PostSignUp)

			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL
			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", string(user.Password))
			values.Add("first-name", user.FirstName)
			values.Add("last-name", user.LastName)

			resp, err := cli.PostForm(srvcUrl+"/signup", values)
			require.ErrorIs(t, err, ErrRedirect)
			require.Equal(t, http.StatusSeeOther, resp.StatusCode)

			for _, cookie := range resp.Cookies() {
				switch cookie.Name {
				case cookiemaker.AUTH_COOKIE_KEY:
					require.Equal(t, AUTH_TOKEN, cookie.Value)
				}
			}
		},
	)

	t.Run(
		"Form error",
		func(t *testing.T) {
			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(new(model.GetUserAuthRow), pgx.ErrNoRows)

			store.
				EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(0)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			srvc := service.NewService(store, validator.Validate)

			authRepo := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/signup", authRepo.PostSignUp)

			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL
			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", string(user.Password))
			values.Add("first_name", user.FirstName)
			values.Add("last_name", user.LastName)

			resp, err := cli.PostForm(srvcUrl+"/signup", values)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Contains(t, string(body), "failed on the 'required' tag")
		},
	)

	t.Run(
		"Username(email) has used",
		func(t *testing.T) {
			ctl := gomock.NewController(t)
			store := mock_model.NewMockStore(ctl)
			store.
				EXPECT().
				GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
				Times(1).
				Return(new(model.GetUserAuthRow), nil)

			store.
				EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(0)

			tm := mock_tokenMaker.NewMockTokenMaker(ctl)
			tm.
				EXPECT().
				MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			srvc := service.NewService(store, validator.Validate)

			authRepo := auth.AuthRepo{
				Service:     srvc,
				View:        vw,
				TokenMaker:  tm,
				FormDecoder: form.NewDecoder(),
			}

			mux := chi.NewMux()
			mux.Post("/signup", authRepo.PostSignUp)

			srv := httptest.NewTLSServer(mux)
			srvcUrl := srv.URL
			values := url.Values{}
			values.Add("email", user.Email)
			values.Add("password", string(user.Password))
			values.Add("first-name", user.FirstName)
			values.Add("last-name", user.LastName)

			resp, err := cli.PostForm(srvcUrl+"/signup", values)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, validator.Validate.Var(string(body), "html"))
			require.Contains(t, string(body), " <span>Username already used</span>")
		},
	)
}
