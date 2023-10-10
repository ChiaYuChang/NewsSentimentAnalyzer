package randangenerator_test

import (
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	val "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestMRand(t *testing.T) {
	an, err := rg.AlphaNum.Clone()
	require.NoError(t, err)

	an.SetRand(rg.NewMRand(1))
	rs1, err := an.GenRdmString(100)
	require.NoError(t, err)

	an.SetRand(rg.NewMRand(1))
	rs2, err := an.GenRdmString(100)
	require.NoError(t, err)

	require.Equal(t, rs1, rs2)
}

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

func TestGenRdnTime(t *testing.T) {
	tf, err := time.Parse(time.DateTime, "2000-01-01 00:00:00")
	require.NoError(t, err)
	tt := time.Now()

	n := 100
	for i := 0; i < n; i++ {
		rdnt := rg.GenRdnTime(tf, tt)
		require.True(t, tf.Before(rdnt))
		require.True(t, tt.After(rdnt))
	}

	rdnts := rg.GenRdnTimes(100, tf, tt)
	require.True(t, tf.Before(rdnts[0]))
	require.True(t, tt.After(rdnts[n-1]))
	for i := 1; i < n-1; i++ {
		require.True(t, rdnts[i-1].Before(rdnts[i]))
		require.True(t, rdnts[i+1].After(rdnts[i]))
	}
}
