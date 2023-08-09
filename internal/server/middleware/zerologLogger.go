package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type logFormatter struct {
	zerolog.Level
}

type logEntry struct {
	Level      zerolog.Level
	UserAgent  string
	HttpMethod string
	Addr       string
	Path       string
}

func (lf logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	le := logEntry{
		UserAgent:  r.Header.Get("User-Agent"),
		HttpMethod: r.Method,
		Addr:       r.RemoteAddr,
		Path:       r.URL.Path,
	}
	return le
}

func (le logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	if strings.HasPrefix(le.Path, global.AppVar.App.RoutePattern.HealthCheck) {
		return
	}

	global.Logger.
		WithLevel(le.Level).
		Str("ua", le.UserAgent).
		Str("method", le.HttpMethod).
		Str("addr", le.Addr).
		Str("path", le.Path).
		Int("status", status).
		Int("in", bytes).
		Dur("dur", elapsed).
		Send()
}

func (l logEntry) Panic(v interface{}, stack []byte) {
	global.Logger.Panic().Msg(fmt.Sprintf("%v", v))
}

func NewZerologLogger(level zerolog.Level) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.RequestLogger(logFormatter{Level: level})(next)
	}
}
