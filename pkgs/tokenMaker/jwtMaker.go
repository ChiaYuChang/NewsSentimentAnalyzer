package tokenmaker

import (
	"fmt"
	"strings"
	"time"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/golang-jwt/jwt/v5"
)

var DEFAULT_JWT_SIGN_METHOD = jwt.SigningMethodHS512

type JWTClaims struct {
	UserInfo
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserName string `json:"username"`
	Role     Role   `json:"role"`
}

func (c JWTClaims) GetUserInfo() UserInfo {
	return c.UserInfo
}

func (c JWTClaims) GetRole() Role {
	return c.Role
}

func (c JWTClaims) GetUsername() string {
	return c.UserName
}

func (c JWTClaims) String() string {
	sb := strings.Builder{}
	sb.WriteString("JWT Claims:\n")
	sb.WriteString(fmt.Sprintf("\t - Username  : %s\n", c.UserName))
	sb.WriteString(fmt.Sprintf("\t - Role      : %v\n", c.Role))
	sb.WriteString(fmt.Sprintf("\t - Issuer    : %s\n", c.Issuer))
	sb.WriteString(fmt.Sprintf("\t - Subject   : %s\n", c.Subject))
	sb.WriteString(fmt.Sprintf("\t - Audience  : %v\n", c.Audience))
	sb.WriteString(fmt.Sprintf("\t - ExpiresAt : %v\n", c.ExpiresAt))
	sb.WriteString(fmt.Sprintf("\t - NotBefore : %v\n", c.NotBefore))
	sb.WriteString(fmt.Sprintf("\t - IssuedAt  : %v\n", c.IssuedAt))
	sb.WriteString(fmt.Sprintf("\t - JWT ID    : %v\n", c.ID))
	return sb.String()
}

type JWTClaimsOpt struct {
	OptionName    string
	AddField      func(claims *JWTClaims) error
	ValidateField func(claims *JWTClaims) error
}

func (opt JWTClaimsOpt) SkipAddition() bool {
	return opt.AddField == nil
}

func (opt JWTClaimsOpt) SkipValidation() bool {
	return opt.ValidateField == nil
}

func NewJWTClaimsOpt(optName string, addfield, validatefield func(*JWTClaims) error) JWTClaimsOpt {
	return JWTClaimsOpt{optName, addfield, validatefield}
}

func WithIssuer(issuer string) JWTClaimsOpt {
	return NewJWTClaimsOpt(
		"Issuer",
		func(j *JWTClaims) error {
			j.Issuer = issuer
			return nil
		},
		func(j *JWTClaims) error {
			if j.Issuer != issuer {
				return ec.MustGetErr(JWTErrIssuer)
			}
			return nil
		},
	)
}

func WithJWTID(jID string) JWTClaimsOpt {
	return NewJWTClaimsOpt(
		"ID",
		func(j *JWTClaims) error {
			j.ID = jID
			return nil
		},
		func(j *JWTClaims) error {
			if j.ID != jID {
				return ec.MustGetErr(JWTErrIssuer)
			}
			return nil
		},
	)
}

type JWTMaker struct {
	secret      []byte
	ExpireAfter time.Duration
	ValidAfter  time.Duration
	SignMethod  jwt.SigningMethod
	Options     map[string]JWTClaimsOpt
}

func NewJWTMaker(secret []byte, signMethod jwt.SigningMethod, expireAfter, validAfter time.Duration) JWTMaker {
	return JWTMaker{
		secret, expireAfter,
		validAfter, signMethod,
		map[string]JWTClaimsOpt{}}
}

func NewJWTMakerWithDefaultVal() JWTMaker {
	return NewJWTMaker(
		DEFAULT_SECRET,
		DEFAULT_JWT_SIGN_METHOD,
		DEFAULT_EXPIRE_AFTER,
		DEFAULT_VALID_AFTER,
	)
}

func (jm JWTMaker) GetSecret() []byte {
	srct := make([]byte, len(jm.secret))
	copy(srct, jm.secret)
	return srct
}

func (jm *JWTMaker) UpdateSecret(secret []byte) *JWTMaker {
	jm.secret = make([]byte, len(secret))
	copy(jm.secret, secret)
	return jm
}

func (jm *JWTMaker) WithOptions(options ...JWTClaimsOpt) *JWTMaker {
	for _, opt := range options {
		jm.Options[opt.OptionName] = opt
	}
	return jm
}

func (jm JWTMaker) MakeToken(username string, role Role) (string, error) {
	currTime := time.Now()
	valideAt := currTime.Add(jm.ValidAfter)
	expireAt := valideAt.Add(jm.ExpireAfter)

	claims := JWTClaims{
		UserInfo: UserInfo{
			UserName: username,
			Role:     role,
		},
	}

	claims.IssuedAt = jwt.NewNumericDate(currTime)
	if jm.ValidAfter > 0 {
		claims.NotBefore = jwt.NewNumericDate(valideAt)
	}
	claims.ExpiresAt = jwt.NewNumericDate(expireAt)
	for _, opt := range jm.Options {
		if !opt.SkipAddition() {
			if err := opt.AddField(&claims); err != nil {
				return "", fmt.Errorf(
					"error while adding %s field: %w", opt.OptionName, err)
			}
		}
	}

	token := jwt.NewWithClaims(jm.SignMethod, claims)
	return token.SignedString(jm.secret)
}

func (jm JWTMaker) ValidateToken(tokenStr string) (Payload, error) {
	vertifiedToken, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jm.SignMethod.Alg() {
			return nil, ec.MustGetErr(JWTSignatureInvalid)
		}

		currTime := time.Now()
		if t, _ := t.Claims.GetIssuedAt(); t == nil {
			return nil, ec.MustGetErr(JWTErrIssueAt)
		} else {
			if t.After(currTime) {
				return nil, ec.MustGetErr(JWTErrUsedBeforeIssued)
			}
		}

		if t, _ := t.Claims.GetNotBefore(); t != nil {
			if t.After(currTime) {
				return nil, ec.MustGetErr(JWTErrNotValidYet)
			}
		}

		if t, _ := t.Claims.GetExpirationTime(); t == nil {
			return nil, ec.MustGetErr(JWTClaimsInvalid)
		} else {
			if t.Before(currTime) {
				return nil, ec.MustGetErr(JWTErrExpired)
			}
		}

		if claims, ok := t.Claims.(*JWTClaims); ok {
			for _, opt := range jm.Options {
				if !opt.SkipValidation() {
					if err := opt.ValidateField(claims); err != nil {
						return nil, err
					}
				}
			}
		} else {
			err := ec.MustGetErr(JWTUnverifiable).(*ec.Error)
			err.WithDetails("Claim interface assertion error")
			return nil, err
		}
		return jm.secret, nil
	})

	if err != nil {
		return *new(JWTClaims), err
	}

	return *vertifiedToken.Claims.(*JWTClaims), nil
}

func (jm JWTMaker) String() string {
	sb := strings.Builder{}
	sb.WriteString("JWT Token Maker:\n")
	sb.WriteString(fmt.Sprintf("\t - Secret: %s\n", string(jm.secret)))
	sb.WriteString(fmt.Sprintf("\t - Token valid after    : %5.2f hours\n", jm.ValidAfter.Hours()))
	sb.WriteString(fmt.Sprintf("\t - Token expire after   : %5.2f hours\n", jm.ExpireAfter.Hours()))
	sb.WriteString(fmt.Sprintf("\t - Token signing method : %s\n", jm.SignMethod.Alg()))
	if len(jm.Options) > 0 {
		sb.WriteString("\t - Option Names:\n")
		for optName := range jm.Options {
			sb.WriteString(fmt.Sprintf("\t\t - %s\n", optName))
		}
	}
	return sb.String()
}
