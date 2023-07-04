package object

import (
	"bytes"
	"html/template"
	"sort"
)

type HeadConent struct {
	Meta    *HTMLElementList
	Link    *HTMLElementList
	Script  *HTMLElementList
	hasExec bool
	content template.HTML
}

func (hc *HeadConent) Execute(tmpl *template.Template) error {
	if hc.hasExec {
		return nil
	}

	bf := bytes.NewBufferString("")
	if err := tmpl.Execute(bf, hc); err != nil {
		return err
	}
	hc.hasExec = true
	hc.content = template.HTML(bf.String())
	return nil
}

func (hc HeadConent) HasExec() bool {
	return hc.hasExec
}

func (hc HeadConent) Content() template.HTML {
	return hc.content
}

type Page struct {
	HeadConent
	Title string
}

type ErrorPage struct {
	Page
	ErrorCode          int
	ErrorMessage       string
	ErrorDetail        string
	ShouldAutoRedirect bool
	RedirectPageUrl    string
	RedirectPageName   string
	CountDownFrom      int // second
}

type EndPoint struct {
	HeadConent
	API      string
	EndPoint string
}

type SelectOpts struct {
	OptMap         map[string]string
	MaxDiv         int
	DefaultValue   string
	DefaultText    string
	InsertButtonId string
	DeleteButtonId string
	PositionId     string
	AlertMessage   string
}

func (sopt SelectOpts) SortedOptKey() []string {
	keys := make([]string, 0, len(sopt.OptMap))
	for key := range sopt.OptMap {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))
	return keys
}
