package client

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
)

var ErrTypeAssertionFailure = errors.New("type assertion failure")

type SelectOpts [2]string

type Response interface {
	fmt.Stringer
	Status()
	NTotalResutl() int
}
type Params map[string][]string

func NewParams() Params {
	return make(Params)
}

func (p Params) Add(key, val string) {
	if val == "" {
		return
	}
	p[key] = append(p[key], val)
}

func (p Params) AddList(key string, val ...string) {
	nonEmptyVal := make([]string, 0, len(val))
	for _, v := range val {
		if v != "" {
			nonEmptyVal = append(nonEmptyVal, v)
		}
	}
	p[key] = append(p[key], nonEmptyVal...)
}

func (p Params) Set(key, val string) {
	if val == "" {
		delete(p, key)
		return
	}
	p[key] = []string{val}
}

func (p Params) SetList(key string, val ...string) {
	if val == nil {
		delete(p, key)
		return
	}
	p[key] = val
}

func (p Params) ToUrlVals() url.Values {
	val := url.Values{}
	for k, v := range p {
		fmt.Printf("Add %s\n", strings.Join(v, ","))
		val.Add(k, strings.Join(v, ","))
	}
	return val
}

type Query interface {
	ToRequestURL(u *url.URL) string
}

type pageFormRepoKey struct {
	Name     string
	Endpoint string
}

type PageFormHandler func(apikey string, pageForm pageform.PageForm) (Query, error)

type PageFormRepo map[pageFormRepoKey]PageFormHandler

func NewPageFormRepo() PageFormRepo {
	return PageFormRepo(make(map[pageFormRepoKey]PageFormHandler))
}

func (qb PageFormRepo) RegisterPageForm(pf pageform.PageForm, handler PageFormHandler) {
	qb[pageFormRepoKey{Name: pf.API(), Endpoint: pf.Endpoint()}] = handler
}

func (qb PageFormRepo) Build(apikey string, pf pageform.PageForm) (Query, error) {
	return qb[pageFormRepoKey{Name: pf.API(), Endpoint: pf.Endpoint()}](apikey, pf)
}
