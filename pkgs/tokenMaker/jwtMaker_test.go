package tokenmaker_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tm "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	maker := tm.NewJWTMakerWithDefaultVal()
	require.Equal(t, tm.DEFAULT_SECRET, maker.GetSecret())

	newSrct := []byte("Secret-for-testing")

	t.Run(
		"Update Secret",
		func(t *testing.T) {
			maker.UpdateSecret(newSrct)
			require.Equal(t, newSrct, maker.GetSecret())
		},
	)

	t.Run(
		"Copy when update secret",
		func(t *testing.T) {
			newSrct[0] = 's'
			require.NotEqual(t, newSrct, maker.GetSecret())
		},
	)

	t.Run(
		"Copy when get secret",
		func(t *testing.T) {
			getSrct := maker.GetSecret()
			getSrct[0] = 's'
			require.NotEqual(t, newSrct, maker.GetSecret())
		},
	)
}

func TestJWTMakerWithOptions(t *testing.T) {
	type testCast struct {
		name    string
		options []tm.JWTClaimsOpt
		errCase bool
	}

	tcs := []testCast{
		{
			name: "with issuer",
			options: []tm.JWTClaimsOpt{
				tm.WithIssuer("[[:ISSUER:]]"),
			},
			errCase: false,
		},
		{
			name: "with jwt id",
			options: []tm.JWTClaimsOpt{
				tm.WithJWTID("[[:ID:]]"),
			},
			errCase: false,
		},
		{
			name: "inline new option function",
			options: []tm.JWTClaimsOpt{
				tm.NewJWTClaimsOpt(
					"Username",
					func(j *tm.JWTClaims) error {
						j.ID = "JID"
						return nil
					},
					func(j *tm.JWTClaims) error {
						if j.ID != "JID" {
							return errors.New("JID mismatch")
						}
						return nil
					},
				),
			},
			errCase: false,
		},
		{
			name: "error maker",
			options: []tm.JWTClaimsOpt{
				tm.NewJWTClaimsOpt(
					"Username",
					nil,
					func(j *tm.JWTClaims) error {
						return errors.New("Test Error")
					},
				),
			},
			errCase: true,
		},
	}

	username := "user"
	role := tm.RAdmin
	uid := int32(1)
	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d: %s", i+1, tc.name),
			func(t *testing.T) {
				maker := tm.NewJWTMakerWithDefaultVal()
				maker.WithOptions(tc.options...)
				tokenStr, err := maker.MakeToken(username, uid, role)
				require.NoError(t, err)

				_, err = maker.ValidateToken(tokenStr)
				if tc.errCase {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			},
		)
	}
}

func TestJWTMakerErrors(t *testing.T) {
	var err error
	maker := tm.NewJWTMaker(
		tm.DEFAULT_SECRET,
		tm.DEFAULT_JWT_SIGN_METHOD,
		tm.DEFAULT_JWT_SIGN_METHOD_SIZE,
		3*time.Second,
		1*time.Second,
	)

	username := "user"
	role := tm.RAdmin
	uid := int32(1)
	tokenStr, _ := maker.MakeToken(username, uid, role)
	_, err = maker.ValidateToken(tokenStr)
	require.ErrorContains(t, err, jwt.ErrTokenUnverifiable.Error())
	require.ErrorContains(t, err, ec.MustGetErr(tm.JWTErrNotValidYet).Error())

	time.Sleep(2 * time.Second)
	_, err = maker.ValidateToken(tokenStr)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	_, err = maker.ValidateToken(tokenStr)
	require.Error(t, err)
	require.ErrorContains(t, err, jwt.ErrTokenUnverifiable.Error())
	require.ErrorContains(t, err, ec.MustGetErr(tm.JWTErrExpired).Error())
}

func TestJWTMakerAsMaker(t *testing.T) {
	var maker tm.TokenMaker = tm.NewJWTMakerWithDefaultVal()

	username := "user"
	role := tm.RAdmin
	uid := int32(1)
	tokenStr, err := maker.MakeToken(username, uid, role)
	require.NoError(t, err)
	payload, err := maker.ValidateToken(tokenStr)
	require.NoError(t, err)

	require.Equal(t, username, payload.GetUsername())
	require.Equal(t, role, payload.GetRole())
	require.Equal(t, uid, payload.GetUserID())
}
