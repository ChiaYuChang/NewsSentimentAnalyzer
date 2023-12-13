package global

import "time"

type CtxKey string

const (
	CtxUserInfo CtxKey = "UserInfo"
)

const (
	CLSToken string = "[CLS]"
	SEPToken string = "[SEP]"
)

const (
	CacheExpireVeryLong time.Duration = 31 * 24 * time.Hour //  1 month
	CacheExpireLong     time.Duration = 10 * 24 * time.Hour // 10 days
	CacheExpireDefault  time.Duration = 3 * time.Hour       //  3 hours
	CacheExpireShort    time.Duration = 3 * time.Minute     //  3 minutes
	CacheExpireInstant  time.Duration = 10 * time.Second    // 10 seconds
)

const (
	PREVIEW_CACHE_KEY_PREFIX   = "PCID:"
	PREVIEW_CACHE_KEY_SUFFIX   = ""
	PREVIEW_CACHE_KEY_SALT_LEN = 32
)

const (
	JWT_SECRET_CACHE_PREFIX = "JWT-SECRET:"
	JWT_SECRET_CACHE_SUFFIX = ""
	JWT_SECRET_CACHE_NUM    = 3
)
