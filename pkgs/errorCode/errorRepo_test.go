package errorcode_test

import (
	"errors"
	"net/http"
	"testing"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/stretchr/testify/require"
)

func TestDefaultErrorRepo(t *testing.T) {
	var err error
	var ok bool
	var ecerr *ec.Error

	t.Run(
		"Register an used error code",
		func(t *testing.T) {
			err = ec.RegisterErr(ec.Success, http.StatusOK, "OK")
			require.NotNil(t, err)
			require.Error(t, err)
			require.Equal(t, err, ec.MustGetErr(ec.ECCodeHasBeenUsed))

			err, ok = ec.GetErr(ec.Success)
			require.True(t, ok)
			require.NotNil(t, err)

			ecerr, ok = err.(*ec.Error)
			require.True(t, ok)
			require.Equal(t, ecerr.ErrorCode, ec.Success)
		},
	)

	t.Run(
		"Register error code OK",
		func(t *testing.T) {
			ECTestError := ec.ErrorCode(-1)
			err = ec.RegisterErr(ECTestError, -1, "For Test")
			require.Nil(t, err)

			err, ok = ec.GetErr(ECTestError)
			require.True(t, ok)
			require.NotNil(t, err)
			ecerr, ok = err.(*ec.Error)
			require.True(t, ok)
			require.Equal(t, ecerr.ErrorCode, ECTestError)
		},
	)

	t.Run(
		"Register error code from an error OK",
		func(t *testing.T) {
			err = errors.New("TestFromError")
			ECTestFromError := ec.ErrorCode(-2)
			ec.RegisterFromErr(ECTestFromError, -1, err)
			err, ok = ec.GetErr(ECTestFromError)
			require.True(t, ok)
			require.NotNil(t, err)
			ecerr, ok = err.(*ec.Error)
			require.True(t, ok)
			require.Equal(t, ecerr.ErrorCode, ECTestFromError)
		},
	)

	t.Run(
		"Register error by WithOpts method",
		func(t *testing.T) {
			ECTestWithOpts := ec.ErrorCode(-3)
			opts := func(repo ec.ErrorRepo) error {
				return repo.RegisterErr(ECTestWithOpts, -1, "For Test")
			}
			err = ec.WithOptions(opts)
			require.Nil(t, err)
			err, ok = ec.GetErr(ECTestWithOpts)
			require.True(t, ok)
			require.NotNil(t, err)
			ecerr, ok = err.(*ec.Error)
			require.True(t, ok)
			require.Equal(t, ecerr.ErrorCode, ECTestWithOpts)
		},
	)
}
