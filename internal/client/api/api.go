package api

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
)

var ErrTypeAssertionFailure = errors.New("type assertion failure")
var ErrNotNextPage = errors.New("there are no more pages to query")

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
		val.Add(k, strings.Join(unique(v), ","))
	}
	return val
}

func unique(v []string) []string {
	set := map[string]struct{}{}
	for _, e := range v {
		set[e] = struct{}{}
	}
	u := make([]string, 0, len(set))
	for k := range set {
		u = append(u, k)
	}
	sort.Sort(sort.StringSlice(u))
	return u
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
	String() string
	Params() Params
	ToRequest() (*http.Request, error)
}

type Response interface {
	String() string
	GetStatus() string
	HasNext() bool
	NextPageRequest(body io.Reader) (*http.Request, error)
	Len() int
	ToNews(ctx context.Context, wg *sync.WaitGroup, c chan<- *model.CreateNewsParams)
}

var re = regexp.MustCompile(`[\p{P}\p{Zs}[:punct:]]`)

func MD5Hash(title string, publishedAt time.Time) string {
	text := re.ReplaceAllString(title, "")
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s@%s", text, publishedAt.UTC().Format(time.DateOnly))))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func ToRequest(apiURL, apiMethod, apiKey, apiEndpoint string, body io.Reader, q Query) (*http.Request, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	v := q.Params().ToUrlVals()
	if apiKey != "" {
		v.Add("apikey", apiKey)
	}

	u = u.JoinPath(apiEndpoint)
	u.RawQuery = v.Encode()

	return http.NewRequest(apiMethod, u.String(), body)
}

func ToBeautifulJSON(q Query, prefix, indent string) ([]byte, error) {
	return q.Params().ToBeautifulJSON(prefix, indent)
}

func ToJSON(q Query) ([]byte, error) {
	return q.Params().ToJSON()
}

func ToQueryString(q Query) string {
	return q.Params().ToQueryString()
}
