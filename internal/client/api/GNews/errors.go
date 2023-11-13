package gnews

import (
	"fmt"
	"net/http"
	"strings"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type ErrorResponse struct {
	Error []string `json:"errors"`
}

func (er ErrorResponse) ToError(code int) error {
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

	if len(er.Error) != 0 {
		return ec.MustGetErr(ecCode).(*ec.Error).
			WithDetails(er.Error...)
	}
	return ec.MustGetErr(ecCode)
}

func (er ErrorResponse) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	if len(er.Error) > 0 {
		sb.WriteString("\t- Error:\n")
		for _, val := range er.Error {
			sb.WriteString(fmt.Sprintf("\t  - %s\n", val))
		}
	}
	return sb.String()
}
