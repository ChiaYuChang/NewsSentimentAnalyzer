package view

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type View struct {
	*template.Template
	Errors []error
}

func (v View) HasError() bool {
	return len(v.Errors) > 0
}

func NewView(tmplFuncs template.FuncMap, patterns ...string) (View, error) {
	v := View{Template: template.New("empty").Funcs(tmplFuncs)}
	v.ParseTemplates(patterns...)
	if v.HasError() {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		for _, e := range v.Errors {
			ecErr.WithDetails(e.Error())
		}
		return v, ecErr
	}
	return v, nil
}

func NewViewWithDefaultTemplateFuncs(pattern ...string) (View, error) {
	fm := template.FuncMap{}
	fm["now"] = func(layout string) string {
		return time.Now().Format(layout)
	}
	fm["tomorrow"] = func(layout string) string {
		return time.Now().Add(24 * time.Hour).Format(layout)
	}
	fm["title"] = func(text string) string {
		return strings.ToTitle(text)
	}
	fm["lower"] = func(text string) string {
		return strings.ToLower(text)
	}
	fm["trim"] = func(text string) string {
		return strings.TrimSpace(text)
	}

	fm["div_ceiling"] = func(x, y int) int {
		n := x / y
		if x%y > 0 {
			n++
		}
		return n
	}

	fm["div_floor"] = func(x, y int) int {
		return x / y
	}

	return NewView(fm, pattern...)
}

func (v View) ParseTemplates(patterns ...string) {
	for _, pattern := range patterns {
		var err error
		v.Template, err = v.Template.ParseGlob(pattern)
		if err != nil {
			v.Errors = append(v.Errors, fmt.Errorf("error while ParseGlob(%s): %w", pattern, err))
		}
	}
	return
}

func (v View) AddTemplate(name, text string) error {
	if t, err := v.Template.New(name).Parse(text); err == nil {
		v.Template = t
		return nil
	} else {
		return fmt.Errorf("error while parse template %s: %w", name, err)
	}
}
