package googlecse

import (
	"encoding/json"
	"fmt"
	"net/http"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
		Status string `json:"status"`
	} `json:"error"`
}

func (er ErrorResponse) ToError() error {
	var ecCode ec.ErrorCode
	switch er.Error.Code {
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

	ecErr, _ := ec.GetEcErr(ecCode)
	ecErr.WithMessage(er.Error.Message)
	ecErr.WithDetails(fmt.Sprintf("status: %s", er.Error.Status))
	for _, e := range er.Error.Errors {
		ecErr.WithDetails(e.Message)
	}
	return ecErr
}

func (er ErrorResponse) String() string {
	b, _ := json.Marshal(er)
	return string(b)
}
