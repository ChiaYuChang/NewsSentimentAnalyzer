package tokenmaker

import (
	"fmt"
	"time"

	errorcode "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

func init() {
	errorcode.WithOptions(WithTokenMakerError())
}

var DEFAULT_SECRET = []byte("SHOULD-NEVER-USED-IN-PRODUCTION")

const DEFAULT_EXPIRE_AFTER = 3 * 24 * time.Hour
const DEFAULT_VALID_AFTER = 0 * time.Second

type TokenMaker interface {
	MakeToken(username string, role Role) (string, error)
	ValidateToken(tokenStr string) (Payload, error)
	fmt.Stringer
}

type Role int8

const (
	RUnknown Role = iota
	RAdmin
	RUser
)

func (r Role) String() string {
	var s string
	switch r {
	case RAdmin:
		s = "admin"
	case RUser:
		s = "user"
	default:
		s = "unknown"
	}
	return s
}

func (r *Role) UnmarshalJSON(data []byte) error {
	if len(data) > 1 {
		*r = RUnknown
	}
	*r = Role(data[0] - '0')
	return nil
}

type Payload interface {
	GetRole() Role
	GetUsername() string
	GetUserInfo() UserInfo
	fmt.Stringer
}
