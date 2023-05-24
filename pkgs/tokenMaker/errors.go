package tokenmaker

import (
	"net/http"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTMalformed ec.ErrorCode = 550 + iota
	JWTUnverifiable
	JWTSignatureInvalid
	JWTErrAudience
	JWTErrExpired
	JWTErrUsedBeforeIssued
	JWTErrIssueAt
	JWTErrIssuer
	JWTErrNotValidYet
	JWTErrId
	JWTClaimsInvalid
)

// Register Success to ErrorReop
func WithTokenMakerError() ec.ErrorRepoOption {
	return func(repo ec.ErrorRepo) error {

		for _, e := range []struct {
			code   ec.ErrorCode
			status int
			err    error
		}{
			{JWTMalformed, http.StatusUnauthorized, jwt.ErrTokenMalformed},
			{JWTUnverifiable, http.StatusUnauthorized, jwt.ErrTokenUnverifiable},
			{JWTSignatureInvalid, http.StatusUnauthorized, jwt.ErrSignatureInvalid},
			{JWTErrAudience, http.StatusUnauthorized, jwt.ErrTokenInvalidAudience},
			{JWTErrExpired, http.StatusUnauthorized, jwt.ErrTokenExpired},
			{JWTErrUsedBeforeIssued, http.StatusUnauthorized, jwt.ErrTokenUsedBeforeIssued},
			{JWTErrIssuer, http.StatusUnauthorized, jwt.ErrTokenInvalidIssuer},
			{JWTErrNotValidYet, http.StatusUnauthorized, jwt.ErrTokenNotValidYet},
			{JWTErrId, http.StatusUnauthorized, jwt.ErrTokenInvalidId},
			{JWTClaimsInvalid, http.StatusUnauthorized, jwt.ErrTokenInvalidClaims},
		} {
			err := repo.RegisterErr(e.code, e.status, e.err.Error())
			if err != nil {
				return err
			}
		}
		return nil
	}
}
