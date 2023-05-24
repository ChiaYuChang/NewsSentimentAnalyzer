package newsdata_test

import (
	"encoding/json"
	"net/http"

	newsapi "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/newsAPI"
)

const (
	API_KEY_OK        string = "[[:OK]]"
	API_KEY_DISABLED         = "[[:DISABLE:]]"
	API_KEY_INVALID          = "[[:INVALID:]]"
	API_KEY_EXHAUSTED        = "[[:EXHAUSTED:]]"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var apikey string
		apikey = req.Header.Get("X-ACCESS-KEY")

		if apikey != API_KEY_OK {
			errResp := newsapi.APIError{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
			}
			switch apikey {
			case "":
				errResp.Message = "apiKeyMissing"
			case API_KEY_DISABLED:
				errResp.Message = "apiKeyDisable"
			case API_KEY_EXHAUSTED:
				errResp.Message = "apiKeyExhausted"
			case API_KEY_INVALID:
				errResp.Message = "apiKeyInvalid"
			}
			jsonObj, _ := json.Marshal(errResp)
			w.WriteHeader(errResp.Code)
			w.Write(jsonObj)
			return
		}
		next.ServeHTTP(w, req)
	})
}
