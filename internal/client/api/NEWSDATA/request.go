package newsdata

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/NEWSDATA"
)

const (
	Keyword        api.Key = "q"
	KeywordInTitle api.Key = "qInTitle"
	Country        api.Key = "country"
	Category       api.Key = "category"
	Language       api.Key = "language"
	Domain         api.Key = "domain"
	FromTime       api.Key = "from"
	ToTime         api.Key = "to"
	Page           api.Key = "page"
	APIKey         api.Key = "apikey"
	// WithImage       api.Key = "image"        // currently not support
	// WithVideo       api.Key = "video"        // currently not support
	// WithFullContent api.Key = "full_content" // currently not support
)

type Request struct {
	*api.RequestProto
}

func newRequest(apikey string) *Request {
	r := api.NewRequestProtoType(",")
	r.SetApiKey(apikey)
	return &Request{RequestProto: r}
}

// Keywords or phrases to search for in the article body.
func (r *Request) WithKeywords(keyword string) *Request {
	r.Set(Keyword, keyword)
	return r
}

// Keywords or phrases to search for in the article title.
func (r *Request) WithKeywordsInTitle(keyword string) *Request {
	r.Add(KeywordInTitle, keyword)
	return r
}

// The 2-letter code of the country you want to get headlines for.
func (r *Request) WithCountry(country ...string) *Request {
	for _, c := range country {
		r.Add(Country, c)
	}
	return r
}

// Find sources that display news in a specific language
func (r *Request) WithLanguage(lang ...string) *Request {
	for _, l := range lang {
		r.Add(Language, l)
	}
	return r
}

// The domains to restrict the search to.
func (r *Request) WithDomain(domain string) *Request {
	r.Add(Domain, domain)
	return r
}

// Find sources that display news of this category.
func (r *Request) WithCategory(category ...string) *Request {
	for _, c := range category {
		r.Add(Category, c)
	}
	return r
}

// Token for retrieving next page.
func (r *Request) WithNextPage(page string) *Request {
	r.Set(Page, page)
	return r
}

func (r *Request) withTime(t time.Time, format string, key api.Key) *Request {
	if t.After(API_MIN_TIME) {
		r.Set(key, t.Format(format))
	}
	return r
}

// A date and optional time for the oldest article allowed.
func (r *Request) WithFrom(t time.Time) *Request {
	r.withTime(t, API_TIME_FORMAT, FromTime)
	return r
}

// A date and optional time for the newest article allowed.
func (r *Request) WithTo(t time.Time) *Request {
	r.withTime(t, API_TIME_FORMAT, ToTime)
	return r
}

func (r *Request) SetEndpoint(ep string) (*Request, error) {
	fmt.Println("Set endpoint to: ", ep)

	switch ep {
	case newsdata.EPLatestNews:
		r.RequestProto.SetEndpoint(EPLatestNews)
	case newsdata.EPNewsArchive:
		r.RequestProto.SetEndpoint(EPNewsArchive)
	case newsdata.EPNewsSources:
		r.RequestProto.SetEndpoint(EPNewsSources)
	case newsdata.EPCrypto:
		r.RequestProto.SetEndpoint(EPCrypto)
	default:
		return nil, client.ErrUnknownEndpoint
	}
	return r, nil
}

func (r *Request) ToHttpRequest() (*http.Request, error) {
	req, err := r.RequestProto.
		ToHTTPRequest(API_URL, API_METHOD, nil)
	if err != nil {
		return nil, err
	}

	req = r.AddAPIKeyToQuery(req, APIKey)
	return req, nil
}
