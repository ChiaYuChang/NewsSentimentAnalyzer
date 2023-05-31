package validator

import (
	"unicode"

	val "github.com/go-playground/validator/v10"
)

type password struct {
	text      string
	nDigit    int
	nUpper    int
	nLower    int
	nSpecial  int
	nNotASCII int
}

func NewPassword(s string) password {
	pwd := password{text: s}

	for _, r := range s {
		if r > unicode.MaxASCII {
			pwd.nNotASCII++
			continue
		}

		switch {
		case unicode.IsNumber(r):
			pwd.nDigit++
		case unicode.IsUpper(r):
			pwd.nUpper++
		case unicode.IsLower(r):
			pwd.nLower++
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			pwd.nSpecial++
		}
	}
	return pwd
}

type PasswordValidator struct {
	AcceptASCIIOnly bool
	MinLength       int
	MaxLength       int
	MinDigit        int
	MinUpper        int
	MinLower        int
	MinSpecial      int
}

func NewPasswordValidator(ASCIIOnly bool, minLen, maxLen, minDigit,
	minUpper, minLower, minSpecial int) PasswordValidator {
	return PasswordValidator{
		AcceptASCIIOnly: ASCIIOnly,
		MinLength:       minLen,
		MaxLength:       maxLen,
		MinDigit:        minDigit,
		MinUpper:        minUpper,
		MinLower:        minLower,
		MinSpecial:      minSpecial,
	}
}

func (pv PasswordValidator) IsValid(p password) bool {
	if pv.AcceptASCIIOnly &&
		p.nNotASCII > 0 {
		return false
	}

	if len(p.text) < pv.MinLength ||
		len(p.text) > pv.MaxLength {
		return false
	}

	if p.nDigit <= pv.MinDigit ||
		p.nUpper <= pv.MinUpper ||
		p.nLower <= pv.MinLower ||
		p.nSpecial <= pv.MinSpecial {
		return false
	}
	return true
}

func (pv PasswordValidator) PasswordValFun() val.Func {
	return func(fl val.FieldLevel) bool {
		password, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return pv.IsValid(NewPassword(password))
	}
}

func (pv PasswordValidator) Tag() string {
	return "password"
}
