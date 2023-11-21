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
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/auth"
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
	user, err := testtool.GenRdmUser()
	require.NoError(t, err)
	require.NotNil(t, user)

	cpassword, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	require.NoError(t, err)

	row := &model.GetUserAuthRow{
		ID:       user.ID,
		Email:    user.Email,
		Password: cpassword,
		Role:     user.Role,
	}

	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrRedirect
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	type mocks struct {
		Store      *mock_model.MockStore
		TokenMaker *mock_tokenMaker.MockTokenMaker
	}

	type testCase struct {
		Name     string
		NewMocks func(t *testing.T, row *model.GetUserAuthRow) mocks
		TestFunc func(t *testing.T, srvcUrl string)
	}
	tcs := []testCase{
		{
			Name: "Get login page",
			NewMocks: func(t *testing.T, row *model.GetUserAuthRow) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Any()).
					Times(0)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
				resp, err := cli.Get(srvcUrl + "/login")
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.StatusCode)
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()
				require.NoError(t, validator.Validate.Var(string(body), "html"))
			},
		},
		{
			Name: "Post login page - OK",
			NewMocks: func(t *testing.T, row *model.GetUserAuthRow) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(row, nil)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(
						gomock.Eq(user.Email),
						gomock.Eq(user.ID),
						gomock.Eq(tokenmaker.ParseRole(string(user.Role)))).
					Times(1).
					Return(AUTH_TOKEN, nil)

				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
				values := url.Values{}
				values.Add("email", user.Email)
				values.Add("password", string(user.Password))
				resp, err := cli.PostForm(srvcUrl+"/login", values)
				require.NotNil(t, resp)
				require.Equal(t, http.StatusSeeOther, resp.StatusCode)
				require.ErrorIs(t, err, ErrRedirect)

				for _, cookie := range resp.Cookies() {
					switch cookie.Name {
					case cookiemaker.AUTH_COOKIE_KEY:
						require.Equal(t, AUTH_TOKEN, cookie.Value)
					}
				}
			},
		},
		{
			Name: "Wrong Password",
			NewMocks: func(t *testing.T, row *model.GetUserAuthRow) mocks {
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
				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
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
				require.Contains(t, string(body), "Wrong password. Please try again")
			},
		},
		{
			Name: "Username (email) not found",
			NewMocks: func(t *testing.T, row *model.GetUserAuthRow) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(new(model.GetUserAuthRow), pgx.ErrNoRows)
				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
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
				require.Contains(t, string(body), "Couldn’t find your Account")
			},
		},
		{
			Name: "Database error",
			NewMocks: func(t *testing.T, row *model.GetUserAuthRow) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(new(model.GetUserAuthRow), ErrDatabaseConn)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)

				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
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
		},
		{
			Name: "Token maker error",
			NewMocks: func(t *testing.T, row *model.GetUserAuthRow) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(row, nil)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(
						gomock.Eq(user.Email),
						gomock.Eq(user.ID),
						gomock.Eq(tokenmaker.ParseRole(string(user.Role)))).
					Times(1).
					Return("", ErrMakerMakeToken)
				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
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
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.Name,
			func(t *testing.T) {
				t.Parallel()
				mocks := tc.NewMocks(t, row)
				srvc := service.NewService(mocks.Store, validator.Validate)
				authReup := auth.AuthRepo{
					Service:     srvc,
					View:        vw,
					TokenMaker:  mocks.TokenMaker,
					FormDecoder: form.NewDecoder(),
				}

				mux := chi.NewMux()
				mux.Get("/login", authReup.GetSignIn)
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

				srvr := httptest.NewTLSServer(mux)
				tc.TestFunc(t, srvr.URL)
				srvr.Close()
			},
		)
	}
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

	type mocks struct {
		Store      *mock_model.MockStore
		TokenMaker *mock_tokenMaker.MockTokenMaker
	}

	type testCase struct {
		Name     string
		NewMocks func(t *testing.T) mocks
		TestFunc func(t *testing.T, srvcUrl string)
	}

	tcs := []testCase{
		{
			Name: "Get sign up page",
			NewMocks: func(t *testing.T) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Any()).
					Times(0)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
				resp, err := cli.Get(srvcUrl + "/signup")
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()
				require.NoError(t, validator.Validate.Var(string(body), "html"))
				require.Contains(t, string(body), "<title>Sign up</title>")
			},
		}, {
			Name: "OK",
			NewMocks: func(t *testing.T) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
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
				store.EXPECT().
					CreateUser(gomock.Any(), matcher).
					Times(1).
					Return(user.ID, nil)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(
						gomock.Eq(user.Email),
						gomock.Eq(user.ID),
						gomock.Eq(tokenmaker.ParseRole(string(user.Role)))).
					Times(1).
					Return(AUTH_TOKEN, nil)
				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
				values := url.Values{}
				values.Add("email", user.Email)
				values.Add("password", string(user.Password))
				values.Add("first-name", user.FirstName)
				values.Add("last-name", user.LastName)

				resp, err := cli.PostForm(srvcUrl+"/signup", values)
				require.ErrorIs(t, err, ErrRedirect)
				require.NotNil(t, resp)
				require.Equal(t, http.StatusSeeOther, resp.StatusCode)

				for _, cookie := range resp.Cookies() {
					switch cookie.Name {
					case cookiemaker.AUTH_COOKIE_KEY:
						require.Equal(t, AUTH_TOKEN, cookie.Value)
					}
				}
			},
		},
		{
			Name: "Form error",
			NewMocks: func(t *testing.T) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(new(model.GetUserAuthRow), pgx.ErrNoRows)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				return mocks{Store: store, TokenMaker: tm}
			},
			TestFunc: func(t *testing.T, srvcUrl string) {
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
		},
		{
			Name: "Username(email) has used",
			NewMocks: func(t *testing.T) mocks {
				ctl := gomock.NewController(t)
				store := mock_model.NewMockStore(ctl)
				store.EXPECT().
					GetUserAuth(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(new(model.GetUserAuthRow), nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				tm := mock_tokenMaker.NewMockTokenMaker(ctl)
				tm.EXPECT().
					MakeToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
				return mocks{Store: store, TokenMaker: tm}

			},
			TestFunc: func(t *testing.T, srvcUrl string) {
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
				require.Contains(t, string(body), "Username already used")
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.Name,
			func(t *testing.T) {
				t.Parallel()
				mocks := tc.NewMocks(t)
				srvc := service.NewService(mocks.Store, validator.Validate)

				authRepo := auth.AuthRepo{
					Service:     srvc,
					View:        vw,
					TokenMaker:  mocks.TokenMaker,
					FormDecoder: form.NewDecoder(),
				}

				mux := chi.NewMux()
				mux.Get("/signup", authRepo.GetSignUp)
				mux.Post("/signup", authRepo.PostSignUp)

				srv := httptest.NewTLSServer(mux)
				tc.TestFunc(t, srv.URL)
				srv.Close()
			},
		)
	}
}
