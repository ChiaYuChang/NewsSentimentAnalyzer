package cookieMaker_test

import (
	"net/http"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
)

func TestCookieMaker(t *testing.T) {
	cm := cookieMaker.NewCookieMaker(
		"/",
		"localhost",
		120,
		true,
		true,
		http.SameSiteLaxMode,
	)

	c := cm.NewCookie("key", "val")

	t.Log(c)
}
