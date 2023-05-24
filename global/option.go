package global

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTOption struct {
	Secret      []byte        `json:"secret"`
	ExpireAfter time.Duration `json:"expire_after"`
	ValidAfter  time.Duration `json:"valid_after"`
	SignMethod  jwt.SigningMethod
}
