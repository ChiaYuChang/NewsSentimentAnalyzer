package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
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
	APIVersion  string
	Service     service.Service
	View        view.View
	TokenMaker  tokenmaker.TokenMaker
	CookieMaker *cookiemaker.CookieMaker
	FormDecoder *form.Decoder
}

func NewAuthRepo(version string, srvc service.Service, view view.View,
	tokenmaker tokenmaker.TokenMaker, cookiemaker *cookiemaker.CookieMaker,
	decoder *form.Decoder) AuthRepo {
	return AuthRepo{
		APIVersion:  version,
		Service:     srvc,
		View:        view,
		TokenMaker:  tokenmaker,
		CookieMaker: cookiemaker,
		FormDecoder: decoder,
	}
}

func (repo AuthRepo) GetLogin(w http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie(cookiemaker.AUTH_COOKIE_KEY); err == nil {
		http.Redirect(w, req, fmt.Sprintf("/%s/welcome", repo.APIVersion), http.StatusSeeOther)
	}

	w.WriteHeader(http.StatusOK)
	repo.View.ExecuteTemplate(w, "login.gotmpl", object.LoginPage{
		Page: object.Page{
			HeadConent: view.NewHeadContent(),
			Title:      "Login",
		},
		ShowUsernameNotFountAlert: false,
		ShowPasswordMismatchAlert: false,
	})
	return
}

func (repo AuthRepo) PostLogin(w http.ResponseWriter, req *http.Request) {
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
			},
			Username: auth.Email,
		}

		if errors.Is(err, pgx.ErrNoRows) {
			data.ShowUsernameNotFountAlert = true
		} else if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			data.ShowPasswordMismatchAlert = true
		} else {
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			ecErr.WithDetails(err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}
		repo.View.ExecuteTemplate(w, "login.gotmpl", data)
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
	http.Redirect(w, req, fmt.Sprintf("/%s/welcome", repo.APIVersion), http.StatusSeeOther)
	return
}

func (repo AuthRepo) GetSignUp(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	repo.View.ExecuteTemplate(w, "signup.gotmpl", object.SignUpPage{
		Page: object.Page{
			HeadConent: view.NewHeadContent(),
			Title:      "Sign up",
		},
		ShowUsernameHasUsedAlert: false,
	})
	return
}

func (repo AuthRepo) PostSignUp(w http.ResponseWriter, req *http.Request) {
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
		repo.View.ExecuteTemplate(w, "signup.gotmpl", object.SignUpPage{
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

		// if req.Header.Get("User-Agent") != "" {
		cookie := repo.CookieMaker.NewCookie(
			cookiemaker.AUTH_COOKIE_KEY, bearer)
		http.SetCookie(w, cookie)
		// }
		http.Redirect(w, req, fmt.Sprintf("/%s/welcome", repo.APIVersion), http.StatusSeeOther)
	}
}

func (repo AuthRepo) Logout(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, repo.CookieMaker.DeleteCookie(cookiemaker.AUTH_COOKIE_KEY))
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func (repo AuthRepo) GetChangePassword(w http.ResponseWriter, req *http.Request) {
	pageData := object.ChangePasswordPage{
		Page: object.Page{
			HeadConent: view.NewHeadContent(),
			Title:      "API Key",
		},
		ShowPasswordNotMatchAlert:         false,
		ShowShouldNotUsedOldPasswordAlert: false,
	}
	w.WriteHeader(http.StatusOK)
	_ = repo.View.ExecuteTemplate(w, "change_password.gotmpl", pageData)
	return
}

func (repo AuthRepo) PostChangPassword(w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	err := req.ParseForm()
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("client error: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	var changePasswordInfo pageform.ChangePassword
	if err := repo.FormDecoder.Decode(&changePasswordInfo, req.PostForm); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("error while get auth info: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if err, _, _ := repo.Service.User().Login(
		req.Context(), userInfo.GetUsername(),
		changePasswordInfo.OldPassword); err != nil {
		pageData := object.ChangePasswordPage{
			Page: object.Page{
				HeadConent: view.NewHeadContent(),
				Title:      "API Key",
			},
			ShowPasswordNotMatchAlert:         true,
			ShowShouldNotUsedOldPasswordAlert: false,
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_ = repo.View.ExecuteTemplate(w, "change_password.gotmpl", pageData)
		return
	}

	if changePasswordInfo.OldPassword == changePasswordInfo.NewPassword {
		pageData := object.ChangePasswordPage{
			Page: object.Page{
				HeadConent: view.NewHeadContent(),
				Title:      "API Key",
			},
			ShowPasswordNotMatchAlert:         false,
			ShowShouldNotUsedOldPasswordAlert: true,
		}
		w.WriteHeader(http.StatusOK)
		_ = repo.View.ExecuteTemplate(w, "change_password.gotmpl", pageData)
		return
	}

	if _, err := repo.Service.User().UpdatePassword(
		req.Context(), &service.UserUpdatePasswordRequest{
			ID:       userInfo.GetUserID(),
			Password: changePasswordInfo.NewPassword,
		},
	); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails(err.Error())
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	http.Redirect(w, req, "welcome", http.StatusSeeOther)
	return
}
