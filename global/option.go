package global

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Option struct {
	TokenMaker        JWTOption               `json:"tokenmaker"`
	PasswordValidator PasswordValidatorOption `json:"password_validator"`
	Server            ServerOption            `json:"server"`
}

func (o Option) String() string {
	sb := strings.Builder{}
	sb.WriteString("Option:\n")
	sb.WriteString("- " + o.TokenMaker.String())
	sb.WriteString("- " + o.PasswordValidator.String())
	sb.WriteString("- " + o.Server.String())
	return sb.String()
}

func ReadOption(path string) (*Option, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error while opening option.json: %w", err)
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error while reading token maker option: %w", err)
	}

	var option Option
	err = json.Unmarshal(bs, &option)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshal option: %w", err)
	}

	return &option, nil
}

type JWTOption struct {
	Secret          []byte            `json:"-"`
	SecretLength    int               `json:"secret_len"`
	ExpireAfterHour int               `json:"expire_after"`
	ValidAfterHour  int               `json:"valid_after"`
	SignMethod      jwt.SigningMethod `json:"-"`
}

func (jwtOpt JWTOption) String() string {
	sb := strings.Builder{}

	sb.WriteString("JWT Opions:\n")
	if jwtOpt.Secret == nil || len(jwtOpt.Secret) == 0 {
		sb.WriteString("\t- Secret           : [[:EMPTY:]]\n")
	} else {
		hexSrct := jwtOpt.GetHexSecretString()
		sb.WriteString(fmt.Sprintf("\t- Secret           : %s...\n", hexSrct[:80]))
	}

	sb.WriteString(fmt.Sprintf("\t- Expire After     : %s\n", jwtOpt.ExpireAfter()))
	sb.WriteString(fmt.Sprintf("\t- Valid After      : %s\n", jwtOpt.ValidAfter()))

	if jwtOpt.SignMethod == nil {
		sb.WriteString("\t- Sign Method      : [[:EMPTY:]]\n")
	} else {
		sb.WriteString(fmt.Sprintf("\t- Sign Method      : %s\n", jwtOpt.SignMethod.Alg()))
	}
	return sb.String()
}

func (jwtOpt JWTOption) ExpireAfter() time.Duration {
	return time.Duration(jwtOpt.ExpireAfterHour) * time.Hour
}

func (jwtOpt JWTOption) ValidAfter() time.Duration {
	return time.Duration(jwtOpt.ValidAfterHour) * time.Hour
}

func (jwtOpt *JWTOption) UpdateSecret() error {
	secret := make([]byte, jwtOpt.SecretLength)
	_, err := rand.Read(secret)
	if err != nil {
		return err
	}
	jwtOpt.Secret = secret
	return nil
}

func (jwtOpt JWTOption) GetSecret() []byte {
	if len(jwtOpt.Secret) == 0 {
		jwtOpt.UpdateSecret()
	}
	return jwtOpt.Secret
}

func (jwtOpt JWTOption) GetHexSecretString() string {
	return hex.EncodeToString(jwtOpt.GetSecret())
}

type PasswordValidatorOption struct {
	AcceptASCIIOnly bool `json:"ascii_only"`
	MinLength       int  `json:"min_length"`
	MaxLength       int  `json:"max_length"`
	MinDigit        int  `json:"min_digit"`
	MinUpper        int  `json:"min_upper"`
	MinLower        int  `json:"min_lower"`
	MinSpecial      int  `json:"min_special"`
}

func (pwdOpt PasswordValidatorOption) String() string {
	sb := strings.Builder{}
	sb.WriteString("Password Validator:\n")
	sb.WriteString(fmt.Sprintf("\t- Only ASCII       : %v\n", pwdOpt.AcceptASCIIOnly))
	sb.WriteString(fmt.Sprintf("\t- Password Len     : (%d, %d)\n",
		pwdOpt.MinLength, pwdOpt.MaxLength))
	sb.WriteString(fmt.Sprintf("\t- Min # of digit   : %d\n", pwdOpt.MinDigit))
	sb.WriteString(fmt.Sprintf("\t- Min # of upper   : %d\n", pwdOpt.MinUpper))
	sb.WriteString(fmt.Sprintf("\t- Min # of lower   : %d\n", pwdOpt.MinLower))
	sb.WriteString(fmt.Sprintf("\t- Min # of special : %d\n", pwdOpt.MinSpecial))
	return sb.String()
}

type ServerOption struct {
	TemplatePath   []string `json:"template_path"`
	StaticFilePath string   `json:"static_file_path"`
}

func (srvOpt ServerOption) String() string {
	sb := strings.Builder{}
	sb.WriteString("Server Options:\n")
	sb.WriteString(fmt.Sprintf("\t- Templates :\n"))
	for _, tmpl := range srvOpt.TemplatePath {
		sb.WriteString(fmt.Sprintf("\t  - %s\n", tmpl))
	}
	sb.WriteString(fmt.Sprintf("\t- Static Files : %v\n", srvOpt.StaticFilePath))
	return sb.String()
}
