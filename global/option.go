package global

import (
	"encoding/json"
	"time"

	"go.uber.org/ratelimit"
)

type Option struct {
	Token       TokenMakerOption `mapstructure:"token"`
	RateLimiter RateLimiter      `mapstructure:"ratelimiter"`
	Password    PasswordOption   `mapstructure:"password"`
	App         AppOption        `mapstructure:"app"`
}

func (opt Option) String() string {
	o, _ := json.MarshalIndent(opt, "", "\t")
	return string(o)
}

type TokenMakerOption struct {
	secret      []byte          `mapstructure:"-"`
	Signature   string          `mapstructure:"signature"`
	SignMethod  TokenSignMethod `mapstructure:"signMethod"`
	ExpireAfter time.Duration   `mapstructure:"expireAfter"`
	ValidAfter  time.Duration   `mapstructure:"validAfter"`
	Issuer      string          `mapstructure:"issuer,omitempty"`
	Subject     string          `mapstructure:"subject,omitempty"`
	Audience    string          `mapstructure:"audience,omitempty"`
	JWTId       string          `mapstructure:"id,omitempty"`
}

type TokenSignMethod struct {
	Algorthm string `mapstructure:"algorithm"`
	Size     int    `mapstructure:"size"`
}

func (tknOpt TokenMakerOption) String() string {
	opt, err := json.MarshalIndent(tknOpt, "", "\t")
	if err != nil {
		return ""
	}
	return string(opt)
}

func (tknOpt TokenMakerOption) Secret() []byte {
	cpSrct := make([]byte, len(tknOpt.secret))
	copy(cpSrct, tknOpt.secret)
	return cpSrct
}

func (tknOpt *TokenMakerOption) SetSecret(secret []byte) {
	tknOpt.secret = secret
}

type RateLimiter struct {
	Auth RateLimit `mapstructure:"auth"`
	API  RateLimit `mapstructure:"api"`
	User RateLimit `mapstructure:"user"`
}

type RateLimit struct {
	N   int           `mapstructure:"n"`
	Per time.Duration `mapstructure:"per"`
}

func (rlOpt RateLimit) RateLimitOption() ratelimit.Option {
	return ratelimit.Per(rlOpt.Per)
}

type PasswordOption struct {
	ASCIIOnly     bool `mapstructure:"asciiOnly"`
	MinLength     int  `mapstructure:"minLength"`
	MaxLength     int  `mapstructure:"maxLength"`
	MinNumDigit   int  `mapstructure:"minNumDigit"`
	MinNumUpper   int  `mapstructure:"minNumUpper"`
	MinNumLower   int  `mapstructure:"minNumLower"`
	MinNumSpecial int  `mapstructure:"minNumSpecial"`
}

func (pwdOpt PasswordOption) String() string {
	opt, err := json.MarshalIndent(pwdOpt, "", "\t")
	if err != nil {
		return ""
	}
	return string(opt)
}

type AppOption struct {
	Template     []string     `mapstructure:"template"`
	StaticFile   StaticFile   `mapstructure:"staticFile"`
	RoutePattern RoutePattern `mapstructure:"routePattern"`
}

type RoutePattern struct {
	Page        map[string]string `mapstructure:"page"`
	ErrorPage   map[string]string `mapstructure:"errorPage"`
	StaticPage  string            `mapstructure:"staticPage"`
	HealthCheck string            `mapstructure:"healthCheck"`
}

type StaticFile struct {
	Path      string            `mapstructure:"path"`
	SubFolder map[string]string `mapstructure:"subFolder"`
}

func (appOpt AppOption) String() string {
	opt, err := json.MarshalIndent(appOpt, "", "\t")
	if err != nil {
		return ""
	}
	return string(opt)
}
