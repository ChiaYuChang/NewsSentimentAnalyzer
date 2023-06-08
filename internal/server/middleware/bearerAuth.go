package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
)

type CtxKey string

const (
	CtxUserInfo CtxKey = "UserInfo"
)

type BearerTokenMaker struct {
	AllowFromHTTPCookie bool
	tokenmaker.TokenMaker
}

func NewJWTTokenMaker(makerOpt global.JWTOption, claimOpt ...tokenmaker.JWTClaimsOpt) BearerTokenMaker {
	maker := tokenmaker.NewJWTMaker(
		makerOpt.Secret,
		makerOpt.SignMethod,
		makerOpt.ExpireAfter(),
		makerOpt.ValidAfter())
	maker.WithOptions(claimOpt...)
	return BearerTokenMaker{false, maker}
}

func (bm BearerTokenMaker) BearerAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var isNotFromHeader bool
		var cookie *http.Cookie
		var err error

		// validation
		bearer := r.Header.Get("Authorization")
		bearer = strings.TrimSpace(strings.TrimLeft(bearer, "Bearer"))
		if bearer == "" && bm.AllowFromHTTPCookie {
			isNotFromHeader = true
			cookie, err = r.Cookie(cookiemaker.AUTH_COOKIE_KEY)
			if err != nil {
				ecErr := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
				ecErr.WithDetails(err.Error())
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(ecErr.HttpStatusCode)
				w.Write(ecErr.MustToJson())
				return
			}
			bearer = cookie.Value
		}

		payload, err := bm.TokenMaker.ValidateToken(bearer)
		if err != nil {
			ecErr, ok := err.(*ec.Error)
			if !ok {
				ecErr = ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
				ecErr.WithDetails(err.Error())
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}

		// sign new token
		userInfo := payload.GetUserInfo()
		w.Header().Add("Authorization", "Trailer-JWT")
		r = r.WithContext(context.WithValue(r.Context(), CtxUserInfo, userInfo))

		bearer, _ = bm.TokenMaker.MakeToken(userInfo.UserName, userInfo.UID, userInfo.Role)
		if isNotFromHeader {
			fmt.Println("Set Cookie")
			http.SetCookie(w, &http.Cookie{
				Name:  cookiemaker.AUTH_COOKIE_KEY,
				Value: bearer,
				Path:  cookie.Path,
			})
		}

		next.ServeHTTP(w, r)
		// update JWT
		w.Header().Set("Trailer-JWT", fmt.Sprintf("Bearer %s", bearer))
		return
	})
}
