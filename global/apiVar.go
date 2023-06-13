package global

import ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"

var AppVar = appVar{}

type appVar struct {
	*Secret
	*Option
	Endpoints []Endpoint
}

func ReadAppVar(secret string, option string, endpoint string) error {
	var err error
	ecErr := ec.MustGetEcErr(ec.ECServerError)

	AppVar.Secret, err = ReadSecret(secret)
	if err != nil {
		ecErr.WithDetails(err.Error())
	}

	AppVar.Option, err = ReadOption(option)
	if err != nil {
		ecErr.WithDetails(err.Error())
	}

	AppVar.Endpoints, err = ReadEndpoints(endpoint)
	if err != nil {
		ecErr.WithDetails(err.Error())
	}

	if len(ecErr.Details) > 0 {
		return ecErr
	}
	return nil
}
