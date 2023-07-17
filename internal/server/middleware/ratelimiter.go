package middleware

import (
	"go.uber.org/ratelimit"
)

func NewRateLimiter(rate int, opts ratelimit.Option) {
	rl := ratelimit.New(rate, opts)

}
