package grpc

import errorcode "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"

type PbError interface {
	GetCode() int64
	GetMessage() string
	GetDetails() []string
}

func NewEcError(pbErr PbError) *errorcode.Error {
	return &errorcode.Error{
		ErrorCode: errorcode.ErrorCode(int(pbErr.GetCode())),
		Message:   pbErr.GetMessage(),
		Details:   pbErr.GetDetails(),
	}
}
