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
	CtxUserInfo CtxKey = "UserInfo"
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

func (bm BearerTokenMaker) BearerAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validation
		bearer := r.Header.Get("Authorization")
		bearer = strings.TrimSpace(strings.TrimLeft(bearer, "Bearer"))
		payload, err := bm.TokenMaker.ValidateToken(bearer)
		if err != nil {
			ecErr, ok := err.(*ec.Error)
			if !ok {
				ecErr = ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
				ecErr.WithDetails(err.Error())
			}
			body, _ := ecErr.ToJson()
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(body)
			return
		}

		// sign new token
		userInfo := payload.GetUserInfo()
		bearer, _ = bm.TokenMaker.MakeToken(userInfo.UserName, userInfo.Role)
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
		r = r.WithContext(context.WithValue(r.Context(), CtxUserInfo, userInfo))
		next.ServeHTTP(w, r)
		return
	})
}
