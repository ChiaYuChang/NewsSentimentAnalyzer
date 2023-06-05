package view

import (
	"fmt"
	"html/template"
)

func ParseTemplates(pattern string, tmplFuncs template.FuncMap) (*template.Template, error) {
	tmpl := template.New("empty").Funcs(tmplFuncs)
	return tmpl.ParseGlob(pattern)
}

func AddTemplate(tmpl *template.Template, name, text string) error {
	if t, err := tmpl.New(name).Parse(text); err == nil {
		tmpl = t
		return nil
	} else {
		return fmt.Errorf("error while parse template %s: %w", name, err)
	}
}
