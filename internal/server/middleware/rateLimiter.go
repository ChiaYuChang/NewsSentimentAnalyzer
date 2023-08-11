package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/go-redis/redis_rate/v10"
	rrate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	ratelimit.Limiter
}

func NewRateLimiter(rate int, opts ratelimit.Option) RateLimiter {
	return RateLimiter{Limiter: ratelimit.New(rate, opts)}
}

func (rl RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rl.Take()
		next.ServeHTTP(w, req)
	})
}

type RedisRateLimiter struct {
	Limit int
	*rrate.Limiter
}

func NewRedisRateLimiter(rdb *redis.Client, n int, per time.Duration) RedisRateLimiter {
	limit := int(float64(n) / per.Seconds())
	if limit < 1 {
		limit = 1
	}

	return RedisRateLimiter{
		Limit:   limit,
		Limiter: rrate.NewLimiter(rdb),
	}
}

func (rl RedisRateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		payload, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
		if !ok {
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			ecErr.WithDetails("user information not found")
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}

		res, err := rl.Allow(req.Context(),
			payload.GetSessionID(), redis_rate.PerSecond(rl.Limit))
		if err != nil {
			ecErr := ec.MustGetEcErr(ec.ECServerError)
			ecErr.WithDetails(err.Error())
			w.WriteHeader(ecErr.HttpStatusCode)
			w.Write(ecErr.MustToJson())
			return
		}

		w.Header().Set("X-Rate-Limit-Remaining", strconv.Itoa(res.Remaining))
		if res.Allowed == 0 {
			// reach rate limit
			seconds := int(res.RetryAfter / time.Second)
			w.Header().Set("X-Rate-Limit-Retry-After", strconv.Itoa(seconds))
			// ecErr := ec.MustGetEcErr(ec.ECTooManyRequests)
			// ecErr.WithDetails(fmt.Sprintf("Too may request, please retry after %s second", strconv.Itoa(res.Remaining)))
			// w.WriteHeader(ecErr.HttpStatusCode)
			// w.Write(ecErr.MustToJson())
			http.Redirect(w, req, global.AppVar.App.RoutePattern.ErrorPage["too-many-request"], http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, req)
	})
}
