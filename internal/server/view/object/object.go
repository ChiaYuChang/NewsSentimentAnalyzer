package object

import (
	"bytes"
	"encoding/json"
	"html/template"
	"sort"
)

type HeadConent struct {
	Meta    *HTMLElementList `json:"meta"`
	Link    *HTMLElementList `json:"link"`
	Script  *HTMLElementList `json:"script"`
	hasExec bool             `json:"-"`
	content template.HTML    `json:"-"`
}

func (hc *HeadConent) FromJson(data []byte, tmpl *template.Template) error {
	err := json.Unmarshal(data, hc)
	if err != nil {
		return err
	}
	return hc.Execute(tmpl)
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

func (hc1 HeadConent) Copy() HeadConent {
	hc2 := HeadConent{}
	hc2.Meta = hc1.Meta.Copy()
	hc2.Link = hc1.Link.Copy()
	hc2.Script = hc1.Script.Copy()
	hc2.hasExec = hc1.hasExec
	hc2.content = hc1.content
	return hc2
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
	revMap := make(map[string]string, len(sopt.OptMap))
	for k, v := range sopt.OptMap {
		revMap[v] = k
	}

	revKeys := make([]string, 0, len(revMap))
	for rkey := range revMap {
		revKeys = append(revKeys, rkey)
	}
	sort.Sort(sort.StringSlice(revKeys))

	keys := make([]string, 0, len(revMap))
	for _, rkey := range revKeys {
		keys = append(keys, revMap[rkey])
	}
	return keys
}
