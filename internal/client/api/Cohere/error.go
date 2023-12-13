package cohere

import ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"

type ErrorResponseBody struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (erb ErrorResponseBody) ToEcError() *ec.Error {
	ecErr := ec.MustGetEcErr(ec.ErrorCode(erb.Code))
	ecErr.Message = erb.Message
	return ecErr
}
