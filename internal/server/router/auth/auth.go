package auth

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	Service     service.Service
	Template    *template.Template
	TokenMaker  tokenmaker.TokenMaker
	CookieMaker *cookiemaker.CookieMaker
	FormDecoder *form.Decoder
}

func (repo AuthRepo) ShowLoginPage(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	repo.Template.ExecuteTemplate(w, "login.gotmpl", object.LoginPage{
		Page: object.Page{
			HeadConent: view.NewHeadContent(),
			Title:      "Login",
		},
		ShowUsernameNotFountAlert: false,
		ShowPasswordMismatchAlert: false,
	})
	return
}

func (repo AuthRepo) Login(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("client error: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	var auth pageform.AuthInfo
	repo.FormDecoder.Decode(&auth, req.PostForm)
	err, uid, role := repo.Service.User().
		Login(context.Background(), auth.Email, auth.Password)

	if err != nil {
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
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}
		repo.Template.ExecuteTemplate(w, "login.gotmpl", data)
		return
	}

	bearer, err := repo.TokenMaker.
		MakeToken(auth.Email, uid, tokenmaker.ParseRole(role))
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails(fmt.Sprintf("error while making signature: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	http.SetCookie(w,
		repo.CookieMaker.NewCookie(
			cookiemaker.AUTH_COOKIE_KEY,
			bearer))
	http.Redirect(w, req, "/v1/welcome", http.StatusSeeOther)
	return
}

func (repo AuthRepo) ShowSignUpPage(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	repo.Template.ExecuteTemplate(w, "signup.gotmpl", object.SignUpPage{
		Page: object.Page{
			HeadConent: view.NewHeadContent(),
			Title:      "Sign up",
		},
		ShowUsernameHasUsedAlert: false,
	})
	return
}

func (repo AuthRepo) SignUp(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("client error: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	var signUpInfo pageform.SignUpInfo
	err = repo.FormDecoder.Decode(&signUpInfo, req.PostForm)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("error while get auth info: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	_, err = repo.Service.User().
		GetAuthInfo(context.Background(), signUpInfo.Email)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		repo.Template.ExecuteTemplate(w, "signup.gotmpl", object.SignUpPage{
			Page: object.Page{
				HeadConent: view.NewHeadContent(),
				Title:      "Sign up",
			},
			ShowUsernameHasUsedAlert: true,
		})
		return
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		if err != nil {
			ecErr.WithDetails(fmt.Sprintf("error while get auth info: %s", err))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
	}

	r := &service.UserCreateRequest{
		FirstName: signUpInfo.FirstName,
		LastName:  signUpInfo.LastName,
		Password:  signUpInfo.Password,
		Email:     signUpInfo.Email,
		Role:      "user",
	}

	if uid, err := repo.Service.
		User().Create(context.Background(), r); err != nil {
		var ecErr *ec.Error
		switch e := err.(type) {
		case val.ValidationErrors:
			ecErr = ec.MustGetEcErr(ec.ECBadRequest)
		case *ec.Error:
			ecErr = e
		default:
			ecErr = ec.MustGetEcErr(ec.ECServerError)
		}
		ecErr.WithDetails(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
	} else {
		bearer, err := repo.TokenMaker.MakeToken(r.Email, uid, tokenmaker.ParseRole(r.Role))
		if err != nil {
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			ecErr.WithDetails(fmt.Sprintf("error while making signature: %s", err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
		}

		if req.Header.Get("User-Agent") != "" {
			cookie := repo.CookieMaker.NewCookie(
				cookiemaker.AUTH_COOKIE_KEY,
				bearer)
			http.SetCookie(w, cookie)
		}
		http.Redirect(w, req, "/v1/welcome", http.StatusSeeOther)
	}
}
