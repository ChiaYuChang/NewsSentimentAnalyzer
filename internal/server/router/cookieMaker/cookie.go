package cookieMaker

import (
	"net/http"
	"time"
)

// const AUTH_COOKIE_KEY string = "__Secure-JWT-Token"
const AUTH_COOKIE_KEY string = "JWT-Token"

var maker *CookieMaker

func init() {
	if maker == nil {
		maker = NewTestCookieMaker()
	}
}

type CookieMaker struct {
	Path     string
	Domain   string
	Expires  time.Time
	MaxAge   int //sec
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

func NewTestCookieMaker() *CookieMaker {
	return NewCookieMaker("/", "localhost", 10*60, false, true, http.SameSiteLaxMode)
}

func NewDefaultCookieMacker(domain string) *CookieMaker {
	return NewCookieMaker("/", domain, 10*60, true, true, http.SameSiteLaxMode)
}

func NewCookieMaker(path, domain string, maxAge int,
	secure bool, httpOnly bool, sameSite http.SameSite) *CookieMaker {
	return &CookieMaker{
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge, // 10 mins
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteLaxMode,
	}
}

func (cm CookieMaker) DeleteCookie(key string) *http.Cookie {
	return &http.Cookie{
		Name:   key,
		MaxAge: -1,
	}
}

func (cm CookieMaker) NewCookie(key, val string) *http.Cookie {
	cookie := &http.Cookie{
		Name:  key,
		Value: val,
	}

	if cm.Path != "" {
		cookie.Path = cm.Path
	}

	if cm.Domain != "" {
		cookie.Domain = cm.Domain
	}

	if !cm.Expires.IsZero() {
		cookie.Expires = cm.Expires
	}

	cookie.MaxAge = cm.MaxAge
	cookie.Secure = cm.Secure
	cookie.HttpOnly = cm.HttpOnly
	cookie.SameSite = cm.SameSite
	return cookie
}

func NewCookie(key, val string) *http.Cookie {
	return maker.NewCookie(key, val)
}

func DeleteCookie(key string) *http.Cookie {
	return maker.DeleteCookie(key)
}
