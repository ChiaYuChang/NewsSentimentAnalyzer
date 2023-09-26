package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"

	cookiemaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
)

type BearerTokenMaker struct {
	AllowFromHTTPCookie bool
	tokenmaker.TokenMaker
}

func NewJWTTokenMaker(makerOpt global.TokenMakerOption, claimOpt ...tokenmaker.JWTClaimsOpt) BearerTokenMaker {
	maker := tokenmaker.NewJWTMaker(
		makerOpt.Secret(),
		makerOpt.SignMethod.Algorthm,
		makerOpt.SignMethod.Size,
		makerOpt.ExpireAfter,
		makerOpt.ValidAfter)
	maker.WithOptions(claimOpt...)
	return BearerTokenMaker{false, maker}
}

func (bm BearerTokenMaker) BearerAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// var isNotFromHeader bool
		var cookie *http.Cookie
		var err error

		// validation
		bearer := req.Header.Get("Authorization")
		bearer = strings.TrimSpace(strings.TrimLeft(bearer, "Bearer"))
		if bearer == "" && bm.AllowFromHTTPCookie {
			cookie, err = req.Cookie(cookiemaker.AUTH_COOKIE_KEY)
			if err != nil {
				// ecErr := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
				// ecErr.WithDetails(err.Error())
				// w.Header().Add("Content-Type", "application/json")
				// w.WriteHeader(ecErr.HttpStatusCode)
				// w.Write(ecErr.MustToJson())
				// ecErr, ok := err.(*ec.Error)
				http.Redirect(w, req, "/unauthorized", http.StatusSeeOther)
				return
			}
			bearer = cookie.Value
		}

		payload, err := bm.TokenMaker.ValidateToken(bearer)
		if err != nil {
			// ecErr, ok := err.(*ec.Error)
			// if !ok {
			// 	ecErr = ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
			// 	ecErr.WithDetails(err.Error())
			// }
			// w.Header().Add("Content-Type", "application/json")
			// w.WriteHeader(ecErr.HttpStatusCode)
			// w.Write(ecErr.MustToJson())
			global.Logger.Debug().Err(err).Msg("token validation error")
			http.Redirect(w, req, "/unauthorized", http.StatusSeeOther)
			return
		}

		// sign new token
		// userInfo := payload.GetUserInfo()
		w.Header().Add("Authorization", "Trailer-JWT")
		ctx := context.WithValue(
			req.Context(),
			global.CtxUserInfo,
			payload,
		)

		bearer, _ = bm.TokenMaker.MakeToken(
			payload.GetUsername(),
			payload.GetUserID(),
			payload.GetRole(),
		)

		http.SetCookie(w, cookiemaker.NewCookie(cookiemaker.AUTH_COOKIE_KEY, bearer))
		next.ServeHTTP(w, req.WithContext(ctx))
		// update JWT
		w.Header().Set("Trailer-JWT", fmt.Sprintf("Bearer %s", bearer))
		return
	})
}
