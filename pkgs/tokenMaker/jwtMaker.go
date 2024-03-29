package tokenmaker

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/cache"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// default method for generating jwt signature (default: HMAC)
var DEFAULT_JWT_SIGN_METHOD = "HMAC"

// default siz for generating jwt signature (default: 384)
var DEFAULT_JWT_SIGN_METHOD_SIZE = 384

// get jwt.SigningMethod by alg, size (default: HMAC384)
func GetJWTSignMethod(alg string, size int) jwt.SigningMethod {
	var method jwt.SigningMethod

	switch alg {
	case "ED25519", "ed25519":
		method = jwt.SigningMethodEdDSA
	case "ECDSA", "ecdsa":
		switch size {
		case 256:
			method = jwt.SigningMethodES256
		case 384:
			method = jwt.SigningMethodES384
		case 512:
			method = jwt.SigningMethodES512
		}
	case "HMAC", "hmac":
		switch size {
		case 256:
			method = jwt.SigningMethodHS256
		case 384:
			method = jwt.SigningMethodHS384
		case 512:
			method = jwt.SigningMethodHS512
		}
	case "RSA", "rsa":
		switch size {
		case 256:
			method = jwt.SigningMethodRS256
		case 384:
			method = jwt.SigningMethodRS384
		case 512:
			method = jwt.SigningMethodRS512
		}
	case "RSAPSS", "rsapass":
		switch size {
		case 256:
			method = jwt.SigningMethodPS256
		case 384:
			method = jwt.SigningMethodPS384
		case 512:
			method = jwt.SigningMethodPS512
		}
	default:
		method = GetJWTSignMethod(
			DEFAULT_JWT_SIGN_METHOD,
			DEFAULT_JWT_SIGN_METHOD_SIZE,
		)
	}
	return method
}

type JWTClaims struct {
	UserInfo
	jwt.RegisteredClaims
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

func (c JWTClaims) GetUserID() uuid.UUID {
	return c.UID
}

func (c JWTClaims) GetSecretID() uuid.UUID {
	return c.SID
}

func (c JWTClaims) GetSessionID() string {
	return fmt.Sprintf("%s-%s",
		strconv.FormatInt(c.Timestamp, 10),
		c.UID.String(),
	)
}

func (c JWTClaims) String() string {
	sb := strings.Builder{}
	sb.WriteString("JWT Claims:\n")
	sb.WriteString(fmt.Sprintf("\t - Username  : %s\n", c.UserName))
	sb.WriteString(fmt.Sprintf("\t - Role      : %v\n", c.Role))
	sb.WriteString(fmt.Sprintf("\t - UID       : %d\n", c.UID))
	sb.WriteString(fmt.Sprintf("\t - Issuer    : %s\n", c.Issuer))
	sb.WriteString(fmt.Sprintf("\t - Subject   : %s\n", c.Subject))
	sb.WriteString(fmt.Sprintf("\t - Audience  : %v\n", c.Audience))
	sb.WriteString(fmt.Sprintf("\t - ExpiresAt : %v\n", c.ExpiresAt))
	sb.WriteString(fmt.Sprintf("\t - NotBefore : %v\n", c.NotBefore))
	sb.WriteString(fmt.Sprintf("\t - IssuedAt  : %v\n", c.IssuedAt))
	sb.WriteString(fmt.Sprintf("\t - JWT ID    : %v\n", c.ID))
	return sb.String()
}

// options for jwt payload
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
	Cache       *struct {
		Store *cache.RedsiStore
		Key   []string
	}
	SignMethod jwt.SigningMethod
	Options    map[string]JWTClaimsOpt
}

func NewJWTMaker(secret []byte, signMethodAlg string, signMethodSize int, expireAfter, validAfter time.Duration) JWTMaker {
	return JWTMaker{
		secret:      secret,
		ExpireAfter: expireAfter,
		ValidAfter:  validAfter,
		SignMethod:  GetJWTSignMethod(signMethodAlg, signMethodSize),
		Options:     map[string]JWTClaimsOpt{}}
}

func NewJWTMakerWithDefaultVal() JWTMaker {
	return NewJWTMaker(
		DEFAULT_SECRET,
		DEFAULT_JWT_SIGN_METHOD,
		DEFAULT_JWT_SIGN_METHOD_SIZE,
		DEFAULT_EXPIRE_AFTER,
		DEFAULT_VALID_AFTER,
	)
}

func (jm JWTMaker) GetSecret() []byte {
	srct := make([]byte, len(jm.secret))
	copy(srct, jm.secret)
	return srct
}

func (jm *JWTMaker) SetCacheStore(store *cache.RedsiStore) {
	jm.Cache.Store = store
	jm.Cache.Key = make([]string, global.JWT_SECRET_CACHE_NUM)
}

func (jm *JWTMaker) GetSecretFromStore(ctx context.Context, key string) {
	jm.Cache.Store.Get(ctx, key)
}

func (jm *JWTMaker) UpdateSecret(secret []byte) *JWTMaker {
	jm.secret = make([]byte, len(secret))
	copy(jm.secret, secret)
	return jm
}

func (jm *JWTMaker) UpdateRandomSecret() ([]byte, *JWTMaker) {
	secret := make([]byte, len(jm.secret))
	_, _ = rand.Read(secret)
	return secret, jm.UpdateSecret(secret)
}

func (jm *JWTMaker) WithSigningMethod(alg string, size int) *JWTMaker {
	jm.SignMethod = GetJWTSignMethod(alg, size)
	return jm
}

func (jm *JWTMaker) WithOptions(options ...JWTClaimsOpt) *JWTMaker {
	for _, opt := range options {
		jm.Options[opt.OptionName] = opt
	}
	return jm
}

func (jm JWTMaker) MakeToken(username string, uid uuid.UUID, role Role) (string, error) {
	currTime := time.Now()
	valideAt := currTime.Add(jm.ValidAfter)
	expireAt := valideAt.Add(jm.ExpireAfter)

	claims := JWTClaims{
		UserInfo: UserInfo{
			UserName:  username,
			Role:      role,
			UID:       uid,
			Timestamp: time.Now().Unix(),
		},
	}
	if jm.Cache != nil {
		// TODO read secret from db
		claims.SID = uuid.UUID{}
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
