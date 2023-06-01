package randangenerator_test

import (
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	val "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestCharSet(t *testing.T) {
	n := 100
	rs := make([]rune, n)
	for i := 0; i < n; i++ {
		rs[i] = rune(i + 1)
	}
	t.Logf("rs: %s", string(rs))
}

func TestGenRdnEmail(t *testing.T) {
	validate := val.New()

	for i := 0; i < 50000; i++ {
		email, err := rg.GenRdmEmail(rg.AlphaNum, rg.Alphabet)
		require.NoError(t, err)
		err = validate.Var(email, "email")
		require.NoError(t, err)
	}
}

func TestGenRdnPassword(t *testing.T) {
	validate := val.New()
	pwdVal := validator.NewDefaultPasswordValidator()
	validate.RegisterValidation(pwdVal.Tag(), pwdVal.ValFun())

	for i := 0; i < 50000; i++ {
		pwd, err := rg.GenRdmPwd(8, 30, 1, 1, 1, 1)
		require.NoError(t, err)
		err = validate.Var(string(pwd), pwdVal.Tag())
		if err != nil {
			t.Log(string(pwd))
		}
		require.NoError(t, err)
	}
}
