package otp_test

import (
	"crypto"
	crand "crypto/rand"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/otp"
	"github.com/stretchr/testify/require"
)

func TestHOPT(t *testing.T) {
	type testCase struct {
		SecretLen int
		NDigit    int
	}

	tcs := []testCase{
		{32, 6},
		{64, 10},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-(%d, %d)", i+1, tc.SecretLen, tc.NDigit),
			func(t *testing.T) {
				src := make([]byte, tc.SecretLen)
				_, _ = crand.Read(src)

				acc := "admin++@example.com"
				iss := "admin123"

				hotp, err := otp.NewOTP(src, tc.NDigit, crypto.SHA256).
					WithAccountName(acc).
					WithIssuer(iss).
					ToHOTP()
				require.NoError(t, err)
				require.NotNil(t, hotp)

				otpStr, err := hotp.Generate()
				require.NoError(t, err)
				require.Equal(t, tc.NDigit, len(otpStr))
				require.True(t, hotp.Validate(otpStr))

				otpurl, err := hotp.ToUrl()
				require.NoError(t, err)
				require.Contains(t, otpurl, fmt.Sprintf("%s://", otp.SCHEME))
				require.Contains(t, otpurl, url.PathEscape(fmt.Sprintf("%s:%s", iss, acc)))
				require.Contains(t, otpurl, fmt.Sprintf("issuer=%s", iss))
				require.Contains(t, otpurl, fmt.Sprintf("digits=%d", tc.NDigit))
				require.Contains(t, otpurl, "secret=")
				require.Contains(t, otpurl, "counter=")
				require.NotContains(t, otpurl, "period=")
			},
		)
	}
}

func TestTOPT(t *testing.T) {
	type testCase struct {
		SecretLen int
		NDigit    int
		Period    int64
		NInterval int
	}

	tcs := []testCase{
		{
			SecretLen: 128,
			NDigit:    6,
			Period:    30,
			NInterval: 1,
		},
		{
			SecretLen: 256,
			NDigit:    10,
			Period:    10,
			NInterval: 2,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-(%d, %d)", i+1, tc.SecretLen, tc.NDigit),
			func(t *testing.T) {
				src := make([]byte, tc.SecretLen)
				_, _ = crand.Read(src)

				acc := "admin@exmaple.com"
				totp, err := otp.NewOTP(src, tc.NDigit, crypto.SHA1).
					WithAccountName(acc).
					ToTOTP(tc.Period)
				require.NoError(t, err)
				require.NotNil(t, totp)

				tnow := time.Now()
				tnow = time.Unix(tnow.Unix()-(tnow.Unix()%tc.Period), 0)
				t.Logf("%s %d %d\n", tnow.Format("04:05"), tnow.Unix(), tnow.Unix()/tc.Period)
				otpStr, err := totp.GenerateAt(tnow)
				require.NoError(t, err)
				require.Equal(t, tc.NDigit, len(otpStr))
				require.True(t, totp.ValidateAt(otpStr, tnow))

				t.Run(
					fmt.Sprintf("now + interval -1"),
					func(t *testing.T) {
						t1 := tnow.Add(time.Duration(tc.Period-1) * time.Second)
						require.True(t, totp.ValidateAt(otpStr, t1))
					},
				)

				t.Run(
					fmt.Sprintf("now + interval +1"),
					func(t *testing.T) {
						t2 := tnow.Add(time.Duration(tc.Period+1) * time.Second)
						require.False(t, totp.ValidateAt(otpStr, t2))
					},
				)

				t.Run(
					"not yet valid",
					func(t *testing.T) {
						t3 := tnow.Add(-1 * time.Second)
						require.False(t, totp.ValidateAt(otpStr, t3))
					},
				)

				otpurl, err := totp.ToUrl()
				require.NoError(t, err)
				require.Contains(t, otpurl, fmt.Sprintf("%s://", otp.SCHEME))
				require.Contains(t, otpurl, url.PathEscape(acc))
				require.NotContains(t, otpurl, "issuer=")
				require.Contains(t, otpurl, fmt.Sprintf("digits=%d", tc.NDigit))
				require.Contains(t, otpurl, "secret=")
				require.NotContains(t, otpurl, "counter=")
				require.Contains(t, otpurl, "period=")
			},
		)
	}
}
