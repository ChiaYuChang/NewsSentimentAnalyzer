package randangenerator

import (
	"errors"
	"fmt"
	"strings"

	mrand "math/rand"
)

var ErrContainsMultibyteChar = errors.New("the given charset contains one or more multibyte characters")

func Must[T string | []byte | []rune](result T, err error) T {
	return result
}

func GenRdmEmail(cUserName CharSet, cDomainName CharSet) (string, error) {
	usernameLen := mrand.Intn(20) + 10
	domainSegs := mrand.Intn(3) + 2

	username, err := cUserName.GenRdmString(usernameLen)
	if err != nil {
		return "", err
	}

	domain := make([]string, domainSegs)
	for i := 0; i < domainSegs; i++ {
		domain[i], err = cDomainName.GenRdmString(mrand.Intn(10) + 3)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s@%s", username, strings.Join(domain, ".")), nil
}

func GenRdmPwd(minLength, maxLength, minDigit,
	minUpper, minLower, minSpecial int) ([]byte, error) {
	if minUpper+minLower+minDigit+minSpecial > minLength {
		minLength = minUpper + minLower + minDigit + minSpecial
	}

	pwdLen := mrand.Intn(maxLength-minLength) + minLength
	pwd := make([]byte, 0, pwdLen)

	bsU, err := AlphabetUpper.GenRdmBytes(minUpper)
	if err != nil {
		return nil, err
	}
	pwd = append(pwd, bsU...)

	bsL, err := AlphabetLower.GenRdmBytes(minLower)
	if err != nil {
		return nil, err
	}
	pwd = append(pwd, bsL...)

	bsD, err := Digit.GenRdmBytes(minLower)
	if err != nil {
		return nil, err
	}
	pwd = append(pwd, bsD...)

	bsS, err := Special.GenRdmBytes(minLower)
	if err != nil {
		return nil, err
	}
	pwd = append(pwd, bsS...)

	bsR, err := Password.GenRdmBytes(pwdLen - len(pwd))
	if err != nil {
		return nil, err
	}
	pwd = append(pwd, bsR...)
	mrand.Shuffle(pwdLen, func(i, j int) { pwd[i], pwd[j] = pwd[j], pwd[i] })
	return pwd, nil
}
