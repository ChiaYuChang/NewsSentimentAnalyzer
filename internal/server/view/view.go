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
}

func NewView(tmplFuncs template.FuncMap, patterns ...string) (View, error) {
	v := View{Template: template.New("empty").Funcs(tmplFuncs)}
	err := v.ParseTemplates(patterns...)
	return v, err
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
		return strings.Title(text)
	}
	fm["lower"] = func(text string) string {
		return strings.ToLower(text)
	}
	fm["trim"] = func(text string) string {
		return strings.TrimSpace(text)
	}
	return NewView(fm, pattern...)
}

func (v View) ParseTemplates(patterns ...string) error {
	ecErr := ec.MustGetEcErr(ec.ECServerError)
	for _, pattern := range patterns {
		var err error
		v.Template, err = v.Template.ParseGlob(pattern)
		if err != nil {
			ecErr.WithDetails(err.Error())
		}
	}
	return nil
}

func (v View) AddTemplate(name, text string) error {
	if t, err := v.Template.New(name).Parse(text); err == nil {
		v.Template = t
		return nil
	} else {
		return fmt.Errorf("error while parse template %s: %w", name, err)
	}
}
