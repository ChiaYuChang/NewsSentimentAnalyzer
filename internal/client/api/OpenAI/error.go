package openai

import (
	"fmt"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type ErrorResponseBody struct {
	StatusCode int `json:"-"`
	Error      struct {
		Code       any            `json:"code,omitempty"`
		Message    string         `json:"message"`
		Param      *string        `json:"param,omitempty"`
		Type       string         `json:"type"`
		InnerError map[string]any `json:"innererror,omitempty"`
	} `json:"error"`
}

func (erb ErrorResponseBody) ToEcError(statusCode int) *ec.Error {
	ecErr, ok := ec.GetEcErr(ec.ErrorCode(erb.StatusCode))
	if !ok {
		ecErr = ec.MustGetEcErr(ec.ECBadRequest)
		ecErr.WithDetails(fmt.Sprintf("original error code: %d", statusCode))
	}

	if code, ok := erb.Error.Code.(string); ok {
		ecErr.Message = code
	}

	if erb.Error.Param != nil {
		ecErr.WithDetails(*erb.Error.Param)
	}

	ecErr.WithDetails(fmt.Sprintf("type: %s", erb.Error.Type))
	ecErr.WithDetails(erb.Error.Message)
	for key, val := range erb.Error.InnerError {
		ecErr.WithDetails(fmt.Sprintf("inner error %s: %v", key, val))
	}
	return ecErr
}
