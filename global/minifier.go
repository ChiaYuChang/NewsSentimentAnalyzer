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

type MinifierOpt struct {
	Mimetype string
	Minifier minify.MinifierFunc
}

func SetMinifier(opts ...MinifierOpt) *minify.M {
	minifier.Do(func() {
		minifier.M = minify.New()
		for _, opt := range opts {
			minifier.M.AddFunc(opt.Mimetype, opt.Minifier)
		}
	})
	return minifier.M
}
