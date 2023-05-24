package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
)

type CtxKey string

const (
	CtxUserInfo CtxKey = "user_info"
)

type BearerTokenMaker struct {
	tokenmaker.TokenMaker
}

func NewJWTTokenMaker(makerOpt global.JWTOption, claimOpt ...tokenmaker.JWTClaimsOpt) BearerTokenMaker {
	maker := tokenmaker.NewJWTMaker(
		makerOpt.Secret, makerOpt.SignMethod,
		makerOpt.ExpireAfter, makerOpt.ValidAfter)
	maker.WithOptions(claimOpt...)
	return BearerTokenMaker{maker}
}

func (bm BearerTokenMaker) BearerAuthenticator(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.FormValue("Authorization")
		bearer = strings.TrimSpace(strings.TrimLeft(bearer, "Bearer"))
		payload, err := bm.TokenMaker.ValidateToken(bearer)
		if err != nil {
			if ecErr, ok := err.(*ec.Error); ok {
				body, _ := ecErr.ToJson()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(ecErr.HttpStatusCode)
				w.Write(body)
			} else {
				ecErr := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
				ecErr.WithDetails(err.Error())
				w.WriteHeader(ecErr.HttpStatusCode)
				w.Header().Set("Content-Type", "application/json")
				body, _ := ecErr.ToJson()
				w.Write(body)
			}
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), CtxUserInfo, payload.GetUserInfo()))
		next.ServeHTTP(w, r)
	}
}

func (bm BearerTokenMaker) BearerSigner(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usrInfo, ok := r.Context().Value(CtxUserInfo).(tokenmaker.UserInfo)
		if !ok {
			ecErr := ec.MustGetErr(ec.ECServerError).(*ec.Error)
			ecErr.WithDetails("username not found", "user role not found")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Header().Set("Content-Type", "application/json")
			body, _ := ecErr.ToJson()
			w.Write(body)
			return
		}

		bearer, err := bm.TokenMaker.MakeToken(usrInfo.UserName, usrInfo.Role)
		if err != nil {
			if ecErr, ok := err.(*ec.Error); ok {
				body, _ := ecErr.ToJson()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(ecErr.HttpStatusCode)
				w.Write(body)
				return
			} else {
				ecErr := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
				ecErr.WithDetails(err.Error())
				w.WriteHeader(ecErr.HttpStatusCode)
				w.Header().Set("Content-Type", "application/json")
				body, _ := ecErr.ToJson()
				w.Write(body)
				return
			}
		}
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
		next.ServeHTTP(w, r)
	}
}
