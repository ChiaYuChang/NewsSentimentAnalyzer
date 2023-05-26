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

type TokenMakerOption struct {
	Secret          []byte            `json:"-"`
	SecretLength    int               `json:"secret_len"`
	ExpireAfterHour int               `json:"expire_after"`
	ValidAfterHour  int               `json:"valid_after"`
	SignMethod      jwt.SigningMethod `json:"-"`
}

func (tmOpt TokenMakerOption) String() string {
	sb := strings.Builder{}
	sb.WriteString("JWT Opions:\n")
	sb.WriteString(fmt.Sprintf("Secret      : %s\n", tmOpt.GetSecretString()))
	sb.WriteString(fmt.Sprintf("Expire After: %s\n", tmOpt.ExpireAfter()))
	sb.WriteString(fmt.Sprintf("Valid After : %s\n", tmOpt.ValidAfter()))
	sb.WriteString(fmt.Sprintf("Sign Method : %s\n", tmOpt.SignMethod.Alg()))
	return sb.String()
}

func (t TokenMakerOption) ExpireAfter() time.Duration {
	return time.Duration(t.ExpireAfterHour) * time.Hour
}

func (t TokenMakerOption) ValidAfter() time.Duration {
	return time.Duration(t.ValidAfterHour) * time.Hour
}

func (t *TokenMakerOption) UpdateSecret() {
	secret := make([]byte, t.SecretLength)
	_, _ = rand.Read(secret)
	t.Secret = secret
}

func (t TokenMakerOption) GetSecret() []byte {
	if len(t.Secret) == 0 {
		t.UpdateSecret()
	}
	return t.Secret
}

func (t TokenMakerOption) GetSecretString() string {
	return hex.EncodeToString(t.GetSecret())
}

func ReadTokenMakerOption(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error while opening option.json: %w", err)
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error while reading token maker option: %w", err)
	}

	var secret Secret
	err = json.Unmarshal(bs, &secret)
	if err != nil {
		return fmt.Errorf("error while unmarshal secrets: %w", err)
	}
	AppVar.Secret = secret
	return nil
}
