package client

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
)

var ErrTypeAssertionFailure = errors.New("type assertion failure")
var ErrPageFormHandlerHasBeenRegistered = errors.New("the PageFormHandler has already been registered")
var ErrHandlerNotFound = errors.New("unregistered handler")

var PageFormHandlerRepo = NewPageFormHandlerRepo()

func RegisterPageForm(pf pageform.PageForm, handler PageFormHandler) {
	PageFormHandlerRepo.RegisterPageForm(pf, handler)
}

func NewQueryFromPageFrom(apikey string, pf pageform.PageForm) (Query, error) {
	return PageFormHandlerRepo.NewQueryFromPageFrom(apikey, pf)
}

type repoMapKey [2]string

func newRepoMapKey(apiName string, endpointName string) repoMapKey {
	return repoMapKey{apiName, endpointName}
}

func (k repoMapKey) APIName() string {
	return k[0]
}

func (k repoMapKey) EndpointName() string {
	return k[1]
}

func (k repoMapKey) String() string {
	return fmt.Sprintf("%s-%s", k[0], k[1])
}

type PageFormHandler func(apikey string, pageForm pageform.PageForm) (Query, error)

type pageFormHandlerRepo map[repoMapKey]PageFormHandler

func NewPageFormHandlerRepo() pageFormHandlerRepo {
	return pageFormHandlerRepo(make(map[repoMapKey]PageFormHandler))
}

func (repo pageFormHandlerRepo) RegisterPageForm(pf pageform.PageForm, handler PageFormHandler) error {
	key := newRepoMapKey(pf.API(), pf.Endpoint())
	if _, ok := repo[key]; ok {
		return ErrPageFormHandlerHasBeenRegistered
	}
	repo[key] = handler
	return nil
}

func (repo pageFormHandlerRepo) NewQueryFromPageFrom(apikey string, pf pageform.PageForm) (Query, error) {
	key := newRepoMapKey(pf.API(), pf.Endpoint())
	if handler, ok := repo[key]; !ok {
		return nil, ErrHandlerNotFound
	} else {
		return handler(apikey, pf)
	}
}

type SelectOpts [2]string

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
		val.Add(k, strings.Join(v, ","))
	}
	return val
}

func (p Params) ToQueryString() string {
	return p.ToUrlVals().Encode()
}

func (p Params) ToBeautifulJSON(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(p, prefix, indent)
}

func (p Params) ToJSON() ([]byte, error) {
	return json.Marshal(map[string][]string(p))
}

type Query interface {
	ToRequestURL(u *url.URL) string
	ToQueryString() string
	ToBeautifulJSON(prefix, indent string) ([]byte, error)
	ToJSON() ([]byte, error)
}

type Response interface {
	fmt.Stringer
	GetStatus()
	Len() int
	ToNews(c chan<- News)
}

type News struct {
	MD5Hash     string
	Title       string
	Url         string
	Description string
	Content     string
	Source      string
	PublishAt   time.Time
}

var re = regexp.MustCompile(`[\p{P}\p{Zs}[:punct:]]`)

func MD5Hash(title string, publishedAt time.Time) string {
	text := re.ReplaceAllString(title, "")
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s@%s", text, publishedAt.UTC().Format(time.DateTime))))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
