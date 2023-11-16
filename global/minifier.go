package global

import (
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

var minifier minifierSingleton

type minifierSingleton struct {
	*minify.M
	sync.Once
}

func Minifier() *minify.M {
	minifier.Do(func() {
		minifier.M = minify.New()
		minifier.M.AddFunc("application/javascript", js.Minify)
	})
	return minifier.M
}
