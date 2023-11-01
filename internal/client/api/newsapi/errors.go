package newsapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type APIError struct {
	Status  string `json:"status"`
	Code    int    `json:"-"`
	ErrCode string `json:"code"`
	Message string `json:"message"`
}

// See https://newsapi.org/docs/errors
func (apiErr APIError) ToError(code int) error {
	var ecCode ec.ErrorCode
	switch code {
	case http.StatusOK:
		return ec.MustGetErr(ec.Success)
	case http.StatusUnauthorized:
		ecCode = ec.ECUnauthorized
	case http.StatusTooManyRequests:
		ecCode = ec.ECTooManyRequests
	case http.StatusInternalServerError:
		ecCode = ec.ECServerError
	default:
		ecCode = ec.ECBadRequest
	}

	if apiErr.Message != "" {
		return ec.MustGetErr(ecCode).(*ec.Error).
			WithDetails(apiErr.Message)
	}
	return ec.MustGetErr(ecCode)
}

func (apiErr APIError) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	sb.WriteString("\t- Status Code: " + strconv.Itoa(apiErr.Code) + "\n")
	sb.WriteString(fmt.Sprintf("\t- Message    : (%s) %s\n", apiErr.ErrCode, apiErr.Message))
	return sb.String()
}
