package global_test

import (
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestReadOption(t *testing.T) {
	err := global.ReadAppVar(
		"../secret.json",
		"../config/option.json",
		"../config/endpoint.json",
	)
	require.NoError(t, err)

	opt := global.AppVar
	err = opt.TokenMaker.UpdateSecret()
	require.NoError(t, err)
	opt.TokenMaker.SignMethod = jwt.SigningMethodHS512

	oldSrc := opt.TokenMaker.GetHexSecretString()
	err = opt.TokenMaker.UpdateSecret()
	require.NoError(t, err)
	newSrc := opt.TokenMaker.GetHexSecretString()
	require.NotEqual(t, oldSrc, newSrc)

	t.Log(global.AppVar)
}
