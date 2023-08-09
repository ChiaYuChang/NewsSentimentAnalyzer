package gnews

import (
	"fmt"
	"net/http"
	"strings"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type APIError struct {
	Error []string `json:"errors"`
}

func (apiErr APIError) ToError(code int) error {
	var ecCode ec.ErrorCode
	switch code {
	case http.StatusOK:
		return ec.MustGetErr(ec.Success)
	case http.StatusBadRequest:
		ecCode = ec.ECBadRequest
	case http.StatusUnauthorized:
		ecCode = ec.ECUnauthorized
	case http.StatusForbidden:
		ecCode = ec.ECForbidden
	case http.StatusTooManyRequests:
		ecCode = ec.ECTooManyRequests
	case http.StatusInternalServerError:
		ecCode = ec.ECServerError
	case http.StatusServiceUnavailable:
		ecCode = ec.ECServiceUnavailable
	default:
		ecCode = ec.ECBadRequest
	}

	if len(apiErr.Error) != 0 {
		return ec.MustGetErr(ecCode).(*ec.Error).
			WithDetails(apiErr.Error...)
	}
	return ec.MustGetErr(ecCode)
}

func (apiErr APIError) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	if len(apiErr.Error) > 0 {
		sb.WriteString("\t- Error:\n")
		for _, val := range apiErr.Error {
			sb.WriteString(fmt.Sprintf("\t  - %s\n", val))
		}
	}
	return sb.String()
}
