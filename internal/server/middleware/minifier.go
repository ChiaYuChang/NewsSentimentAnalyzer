package middleware

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

func NewMinifier() *minify.M {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	// m.AddFunc("test")

	// m.Middleware()
	return m
}
