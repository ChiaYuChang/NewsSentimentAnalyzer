package otp

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const SCHEME = "otpauth"

type param string

const (
	ParamSecret    param = "secret"
	ParamIssuer    param = "issuer"
	ParamAlgorithm param = "algorithm"
	ParamDigits    param = "digits"
	ParamCounter   param = "counter"
	ParamPeriod    param = "period"
)

type label string

const (
	LabelIssuer      label = "issuer"
	LabelAccountName label = "accountName"
)

type otpGenerator struct {
	secret    []byte
	digit     int
	hash      crypto.Hash
	mod       int
	formatter string
	params    map[param]string
	labels    map[label]string
}

func must(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}

func NewOTP(secret []byte, digit int, hash crypto.Hash) *otpGenerator {
	otp := &otpGenerator{
		params: map[param]string{},
		labels: map[label]string{},
	}

	otp = otp.WithAlgorithm(hash).
		WithSecret(secret).
		WithDigit(digit)

	return otp
}

func NewDefautlOTP(secret []byte) *otpGenerator {
	return NewOTP(secret, 6, crypto.SHA1)
}

func (otp *otpGenerator) WithIssuer(issuer string) *otpGenerator {
	otp.labels[LabelIssuer] = issuer
	otp.params[ParamIssuer] = issuer
	return otp
}

func (otp *otpGenerator) WithAccountName(name string) *otpGenerator {
	otp.labels[LabelAccountName] = name
	return otp
}

func (otp *otpGenerator) WithDigit(digit int) *otpGenerator {
	otp.digit = digit
	otp.params[ParamDigits] = strconv.Itoa(digit)

	mod := 1
	for i := 0; i < digit; i++ {
		mod *= 10
	}
	otp.mod = mod
	otp.formatter = fmt.Sprintf("%%0%dd", digit)
	return otp
}

func (otp *otpGenerator) WithAlgorithm(hash crypto.Hash) *otpGenerator {
	otp.params[ParamAlgorithm] = strings.Replace(hash.String(), "-", "", 1)
	otp.hash = hash
	return otp
}

func (otp *otpGenerator) WithSecret(secret []byte) *otpGenerator {
	base32secret := make([]byte, base32.StdEncoding.EncodedLen(len(secret)))
	base32.StdEncoding.Encode(base32secret, secret)
	otp.secret = base32secret
	otp.params[ParamSecret] = string(base32secret)
	return otp
}

func (otp *otpGenerator) ToHOTP() (*HOTP, error) {
	hotp := &HOTP{otp: otp}
	counter, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, fmt.Errorf("error while rand.Int, %w", err)
	}
	return hotp.WithCounter(counter.Int64()), nil
}

func (otp *otpGenerator) ToTOTP(period int64) (*TOTP, error) {
	if period < 0 {
		return nil, errors.New("interval should be greather than 0")
	}
	totp := &TOTP{otpGenerator: otp}
	return totp.WithPeriod(period), nil
}

func (otp otpGenerator) genNewOTP(internal int64) (string, error) {
	if internal < 0 {
		return "", errors.New("input should be greather than 0")
	}

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(internal))
	hasher := hmac.New(otp.hash.New, otp.secret)
	hasher.Write(bs)
	hashsum := hasher.Sum(nil)
	offset := (hashsum[19] & 0b00001111)

	var header uint32
	r := bytes.NewReader(hashsum[offset : offset+4])
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return "", fmt.Errorf("error while binary.Read: %w", err)
	}

	// retain only the last 31 bits of the resulting value
	otpInt := (int(header) & 0x7FFFFFFF) % otp.mod
	return fmt.Sprintf(otp.formatter, otpInt), nil
}

func (otp otpGenerator) getEncodedSecret() []byte {
	return otp.secret
}

func (otp otpGenerator) getDecodedSecret() []byte {
	secret := make([]byte, base32.StdEncoding.DecodedLen(len(otp.secret)))
	_, _ = base32.StdEncoding.Decode(secret, otp.secret)
	return secret
}

func (otp otpGenerator) toUrl(otpType string) (string, error) {
	params := url.Values{}
	for k, v := range otp.params {
		params.Add(string(k), v)
	}
	switch otpType {
	case "totp":
		params.Add("period", otp.params[ParamPeriod])
	case "hotp":
		params.Add("counter", otp.params[ParamCounter])
	default:
		return "", fmt.Errorf("known otp type: %s", otpType)
	}

	label := ""
	if issuer, ok := otp.labels[LabelIssuer]; !ok {
		label = otp.labels[LabelAccountName]
	} else {
		label = fmt.Sprintf("%s:%s", issuer, otp.labels[LabelAccountName])

	}
	return fmt.Sprintf("%s://%s/%s?%s",
		SCHEME, otpType, url.PathEscape(label), params.Encode()), nil
}

type HOTP struct {
	counter int64
	otp     *otpGenerator
}

func NewHOTP(secret []byte, digit int, hash crypto.Hash) (*HOTP, error) {
	return NewOTP(secret, digit, hash).ToHOTP()
}

func (hotp *HOTP) WithCounter(counter int64) *HOTP {
	hotp.counter = counter
	hotp.otp.params[ParamCounter] = strconv.Itoa(int(counter))
	return hotp
}

func (hotp HOTP) Generate() (string, error) {
	return hotp.otp.genNewOTP(hotp.counter)
}

func (hotp HOTP) Validate(otpStr string) bool {
	state := hotp.counter
	hotp.counter++
	if hotp.counter < 1 {
		hotp.counter = 1
	}
	return otpStr == must(hotp.otp.genNewOTP(state))
}

func (hotp HOTP) Type() string {
	return "hotp"
}

func (hotp HOTP) ToUrl() (string, error) {
	return hotp.otp.toUrl(hotp.Type())
}

type TOTP struct {
	*otpGenerator
	period int64
}

func NewTOTP(secret []byte, period int64, digit int, hash crypto.Hash) (*TOTP, error) {
	if period < 0 {
		return nil, errors.New("interval should be greather than 0")
	}
	return NewOTP(secret, digit, hash).ToTOTP(period)
}

func (totp *TOTP) WithPeriod(period int64) *TOTP {
	totp.period = period
	totp.params[ParamPeriod] = strconv.Itoa(int(period))
	return totp
}

func (totp TOTP) GenerateAt(t time.Time) (string, error) {
	return totp.genNewOTP(t.Unix() / totp.period)
}

func (totp TOTP) GenerateNow() (string, error) {
	return totp.GenerateAt(time.Now())
}

func (totp TOTP) ValidateAt(otpStr string, t time.Time) bool {
	return otpStr == must(totp.GenerateAt(t))
}

func (totp TOTP) ValidateNow(otpStr string) bool {
	return totp.ValidateAt(otpStr, time.Now())
}

func (totp TOTP) ValidateNInterval(otpStr string, nInterval int) bool {
	isValid := false

	t := time.Now()
	for i := 0; !isValid && i < nInterval; i++ {
		isValid = isValid || totp.ValidateAt(otpStr, t.Add(time.Duration(-i*int(totp.period))))
	}
	return isValid
}

func (totp TOTP) Type() string {
	return "totp"
}

func (totp TOTP) ToUrl() (string, error) {
	return totp.toUrl(totp.Type())
}
