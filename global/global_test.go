package global_test

import (
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestReadOption(t *testing.T) {
	opt, err := global.ReadOption("../config/option.json")
	require.NoError(t, err)

	err = opt.TokenMaker.UpdateSecret()
	require.NoError(t, err)
	opt.TokenMaker.SignMethod = jwt.SigningMethodHS512

	oldSrc := opt.TokenMaker.GetHexSecretString()
	err = opt.TokenMaker.UpdateSecret()
	require.NoError(t, err)
	newSrc := opt.TokenMaker.GetHexSecretString()
	require.NotEqual(t, oldSrc, newSrc)

	t.Log(opt)
}
