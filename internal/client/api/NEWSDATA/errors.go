package newsdata

import (
	"fmt"
	"net/http"
	"strings"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type ErrorResponse struct {
	Status string            `json:"status"`
	Result map[string]string `json:"results"`
}

// see https://newsdata.io/documentation/#http-response-codes
func (er ErrorResponse) ToEcError(code int) *ec.Error {
	var ecCode ec.ErrorCode
	switch code {
	case http.StatusOK:
		return ec.MustGetEcErr(ec.Success)
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

	if er.Result["message"] != "" {
		return ec.MustGetEcErr(ecCode).
			WithDetails(er.Result["message"])
	}
	return ec.MustGetEcErr(ecCode)
}

func (er ErrorResponse) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	if len(er.Result) > 0 {
		sb.WriteString("\t- Result:\n")
		for key, val := range er.Result {
			sb.WriteString(fmt.Sprintf("\t  - %s: %s\n", key, val))
		}
	}
	return sb.String()
}
