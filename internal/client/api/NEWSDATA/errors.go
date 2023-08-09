package newsdata

import (
	"fmt"
	"net/http"
	"strings"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type APIError struct {
	Status string            `json:"status"`
	Result map[string]string `json:"results"`
}

// see https://newsdata.io/documentation/#http-response-codes
func (apiErr APIError) ToError(code int) error {
	var ecCode ec.ErrorCode
	switch code {
	case http.StatusOK:
		return ec.MustGetErr(ec.Success)
	case http.StatusUnauthorized:
		ecCode = ec.ECUnauthorized
	case http.StatusForbidden:
		ecCode = ec.ECForbidden
	case http.StatusConflict:
		ecCode = ec.ECConflict
	case http.StatusUnsupportedMediaType:
		ecCode = ec.ECUnsupportedMediaType
	case http.StatusUnprocessableEntity:
		ecCode = ec.ECUnprocessableContent
	case http.StatusTooManyRequests:
		ecCode = ec.ECTooManyRequests
	case http.StatusInternalServerError:
		ecCode = ec.ECServerError
	}

	if apiErr.Result["message"] != "" {
		return ec.MustGetErr(ecCode).(*ec.Error).
			WithDetails(apiErr.Result["message"])
	}
	return ec.MustGetErr(ecCode)
}

func (apiErr APIError) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	if len(apiErr.Result) > 0 {
		sb.WriteString("\t- Result:\n")
		for key, val := range apiErr.Result {
			sb.WriteString(fmt.Sprintf("\t  - %s: %s\n", key, val))
		}
	}
	return sb.String()
}
