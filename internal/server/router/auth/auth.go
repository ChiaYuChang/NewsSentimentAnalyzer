package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/cookieMaker"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/go-playground/form"
	"github.com/go-playground/mold/v4"
	val "github.com/go-playground/validator/v10"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	APIVersion   string
	Service      service.Service
	View         view.View
	TokenMaker   tokenmaker.TokenMaker
	Validator    *val.Validate
	FormDecoder  *form.Decoder
	FormModifier *mold.Transformer
}

func NewAuthRepo(
	version string, srvc service.Service, view view.View,
	tokenmaker tokenmaker.TokenMaker, validator *val.Validate,
	decoder *form.Decoder, modifier *mold.Transformer) AuthRepo {

	return AuthRepo{
		APIVersion:   version,
		Service:      srvc,
		View:         view,
		TokenMaker:   tokenmaker,
		Validator:    validator,
		FormDecoder:  decoder,
		FormModifier: modifier,
	}
}

func (repo AuthRepo) GetSignIn(w http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie(cookiemaker.AUTH_COOKIE_KEY); err == nil {
		http.Redirect(w, req,
			fmt.Sprintf("/%s/%s", repo.APIVersion,
				global.AppVar.App.RoutePattern.Page["welcome"]),
			http.StatusSeeOther)
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

func (repo AuthRepo) PostSignIn(w http.ResponseWriter, req *http.Request) {
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
		global.Logger.Info().Msg("login error")
		data := object.LoginPage{
			Page: object.Page{
				HeadConent: view.NewHeadContent(),
				Title:      "Sign-In",
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

	global.Logger.Debug().Msg("make token")
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

	global.Logger.Debug().Msg("set cookie")
	http.SetCookie(w, cookiemaker.NewCookie(cookiemaker.AUTH_COOKIE_KEY, bearer))

	global.Logger.Debug().Msg("redirect")
	http.Redirect(w, req,
		fmt.Sprintf("/%s/%s", repo.APIVersion, global.AppVar.App.RoutePattern.Page["welcome"]),
		http.StatusSeeOther)
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
	if err := req.ParseForm(); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("client error: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	var signUpInfo pageform.SignUpInfo
	if err := repo.FormDecoder.Decode(&signUpInfo, req.PostForm); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("error while get auth info: %s", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}

	if _, err := repo.Service.User().
		GetAuthInfo(context.Background(), signUpInfo.Email); err == nil {
		w.WriteHeader(http.StatusOK)
		repo.View.ExecuteTemplate(w, "signup.gotmpl", object.SignUpPage{
			Page: object.Page{
				HeadConent: view.NewHeadContent(),
				Title:      "Sign up",
			},
			ShowUsernameHasUsedAlert: true,
		})
		return
	} else {
		if !errors.Is(err, pgx.ErrNoRows) {
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			if err != nil {
				ecErr.WithDetails(fmt.Sprintf("error while get auth info: %s", err))
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
		}
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
		cookie := cookiemaker.NewCookie(cookiemaker.AUTH_COOKIE_KEY, bearer)
		http.SetCookie(w, cookie)
		// }
		http.Redirect(w, req,
			fmt.Sprintf("/%s/%s", repo.APIVersion, global.AppVar.App.RoutePattern.Page["welcome"]),
			http.StatusSeeOther)
	}
}

func (repo AuthRepo) GetSignOut(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, cookiemaker.DeleteCookie(cookiemaker.AUTH_COOKIE_KEY))
	http.Redirect(w, req,
		global.AppVar.App.RoutePattern.Page["sign-in"],
		http.StatusSeeOther)
}

func (repo AuthRepo) GetChangePassword(w http.ResponseWriter, req *http.Request) {
	pageData := object.ChangePasswordPage{
		Page: object.Page{
			HeadConent: view.NewHeadContent(),
			Title:      "API Key",
		},
		OldPassword: object.PasswordInput{
			IdPrefix:              "old",
			Name:                  "old-password",
			PlaceHolder:           "Old Password",
			PasswordStrengthCheck: false,
			PasswordCreteria:      nil,
			AlertMessage:          "Your current password is missing or incorrect.",
		},
		NewPassword: object.PasswordInput{
			IdPrefix:              "new",
			Name:                  "new-password",
			PlaceHolder:           "New Password",
			PasswordStrengthCheck: true,
			PasswordCreteria:      object.GetDefaultPasswordCreteria(),
			AlertMessage:          "Your new password cannot not be the same as your current password.",
		},
	}
	w.WriteHeader(http.StatusOK)
	_ = repo.View.ExecuteTemplate(w, "change_password.gotmpl", pageData)
	return
}

type ChangePasswordResponse struct {
	Status             int      `json:"status"`
	Message            string   `json:"message"`
	PasswordNotMatched bool     `json:"password_not_matched"`
	PasswordNotChanged bool     `json:"password_not_changed"`
	Detail             []string `json:"detail"`
}

type PatchChangePasswordBody struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func (repo AuthRepo) PatchChangePassword(w http.ResponseWriter, req *http.Request) {
	jsn := ChangePasswordResponse{}

	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		jsn.Status = ecErr.HttpStatusCode
		jsn.Message = ecErr.Message
		jsn.Detail = append(jsn.Detail,
			"user information not found",
			"check your login status",
		)
		b, _ := json.Marshal(jsn)

		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		jsn.Status = ecErr.HttpStatusCode
		jsn.Message = ecErr.Message
		jsn.Detail = append(jsn.Detail, fmt.Sprintf("body reading error: %s", err))
		b, _ := json.Marshal(jsn)

		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	var changePasswordInfo PatchChangePasswordBody
	json.Unmarshal(body, &changePasswordInfo)

	global.Logger.Info().
		Str("old", changePasswordInfo.Old).
		Str("new", changePasswordInfo.New).
		Msg("input")

	if changePasswordInfo.Old == "" || changePasswordInfo.New == "" {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		jsn.Status = ecErr.HttpStatusCode
		jsn.Message = ecErr.Message
		jsn.Detail = append(jsn.Detail, "parameter(s) is missing")
		b, _ := json.MarshalIndent(jsn, "", "    ")

		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	if changePasswordInfo.Old == changePasswordInfo.New {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		jsn.Status = ecErr.HttpStatusCode
		jsn.Message = ecErr.Message
		jsn.Detail = append(jsn.Detail, "old password and new password cannot be the same")
		jsn.PasswordNotChanged = true
		b, _ := json.MarshalIndent(jsn, "", "    ")

		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}

	if err, _, _ := repo.Service.User().Login(
		req.Context(), userInfo.GetUsername(),
		changePasswordInfo.Old); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECUnauthorized)
		jsn.Status = ecErr.HttpStatusCode
		jsn.Message = ecErr.Message
		jsn.PasswordNotMatched = true
		jsn.Detail = append(jsn.Detail, "password not matched")
		b, _ := json.MarshalIndent(jsn, "", "    ")

		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	if _, err := repo.Service.User().UpdatePassword(
		req.Context(), &service.UserUpdatePasswordRequest{
			ID:       userInfo.GetUserID(),
			Password: changePasswordInfo.New,
		},
	); err != nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		jsn.Status = ecErr.HttpStatusCode
		jsn.Message = ecErr.Message
		jsn.Detail = append(jsn.Detail, err.Error())
		b, _ := json.MarshalIndent(jsn, "", "    ")

		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	b, _ := json.MarshalIndent(jsn, "", "    ")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}
