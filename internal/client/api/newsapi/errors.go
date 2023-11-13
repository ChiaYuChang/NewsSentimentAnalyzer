package newsapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"-"`
	ErrCode string `json:"code"`
	Message string `json:"message"`
}

// See https://newsapi.org/docs/errors
func (er ErrorResponse) ToError(code int) error {
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

	if er.Message != "" {
		return ec.MustGetErr(ecCode).(*ec.Error).
			WithDetails(er.Message)
	}
	return ec.MustGetErr(ecCode)
}

func (er ErrorResponse) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	sb.WriteString("\t- Status Code: " + strconv.Itoa(er.Code) + "\n")
	sb.WriteString(fmt.Sprintf("\t- Message    : (%s) %s\n", er.ErrCode, er.Message))
	return sb.String()
}
