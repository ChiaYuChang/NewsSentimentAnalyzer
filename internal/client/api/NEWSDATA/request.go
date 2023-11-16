package newsdata

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/NEWSDATA"
	"github.com/google/uuid"
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
	Page string
}

func NewRequest(apikey string) *Request {
	req := api.NewRequestProtoType(srv.API_NAME, ",")
	req.SetApiKey(apikey)
	return &Request{RequestProto: req}
}

// Keywords or phrases to search for in the article body.
func (req *Request) WithKeywords(keyword string) *Request {
	req.Set(Keyword, keyword)
	return req
}

// Keywords or phrases to search for in the article title.
func (req *Request) WithKeywordsInTitle(keyword string) *Request {
	req.Add(KeywordInTitle, keyword)
	return req
}

// The 2-letter code of the country you want to get headlines for.
func (req *Request) WithCountry(country ...string) *Request {
	for _, c := range country {
		req.Add(Country, c)
	}
	return req
}

// Find sources that display news in a specific language
func (req *Request) WithLanguage(lang ...string) *Request {
	for _, l := range lang {
		req.Add(Language, l)
	}
	return req
}

// The domains to restrict the search to.
func (req *Request) WithDomain(domain string) *Request {
	req.Add(Domain, domain)
	return req
}

// Find sources that display news of this category.
func (req *Request) WithCategory(category ...string) *Request {
	for _, c := range category {
		req.Add(Category, c)
	}
	return req
}

// Token for retrieving next page.
func (req *Request) WithPage(page string) *Request {
	req.Page = page
	return req
}

func (req *Request) withTime(t time.Time, format string, key api.Key) *Request {
	if t.After(API_MIN_TIME) {
		req.Set(key, t.Format(format))
	}
	return req
}

// A date and optional time for the oldest article allowed.
func (req *Request) WithFrom(t time.Time) *Request {
	req.withTime(t, API_TIME_FORMAT, FromTime)
	return req
}

// A date and optional time for the newest article allowed.
func (r *Request) WithTo(t time.Time) *Request {
	r.withTime(t, API_TIME_FORMAT, ToTime)
	return r
}

func (req *Request) SetEndpoint(ep string) (*Request, error) {
	switch ep {
	case srv.EPLatestNews, EPLatestNews:
		req.RequestProto.SetEndpoint(EPLatestNews)
	case srv.EPNewsArchive, EPNewsArchive:
		req.RequestProto.SetEndpoint(EPNewsArchive)
	case srv.EPNewsSources, EPNewsSources:
		req.RequestProto.SetEndpoint(EPNewsSources)
	case srv.EPCrypto, EPCrypto:
		req.RequestProto.SetEndpoint(EPCrypto)
	default:
		return nil, client.ErrUnknownEndpoint
	}
	return req, nil
}

func (req *Request) ToHttpRequest() (*http.Request, error) {
	httpReq, err := req.RequestProto.ToHTTPRequest(API_URL, API_METHOD, nil)
	if err != nil {
		return nil, err
	}

	p, err := req.Params.Clone()
	if err != nil {
		return nil, err
	}

	if req.Page != "" {
		p.Set(Page, req.Page)
	}
	httpReq.URL.RawQuery = p.Encode()
	httpReq.Header.Set("X-ACCESS-KEY", req.APIKey())
	return httpReq, nil
}

func (req Request) ToPreviewCache(uid uuid.UUID) (cKey string, c *api.PreviewCache) {
	return req.RequestProto.ToPreviewCache(uid, api.StrNextPageToken(req.Page), nil)
}

func RequestFromPreviewCache(cq api.CacheQuery) (api.Request, error) {
	if cq.NextPage.Equal(api.StrLastPageToken) {
		return nil, api.ErrNotNextPage
	}

	var err error
	req := NewRequest(cq.API.Key)
	_, err = req.SetEndpoint(cq.API.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error while set endpoint: %w", err)
	}

	req.Values, err = url.ParseQuery(cq.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("error while parsing raw query: %w", err)
	}

	token, ok := cq.NextPage.(api.StrNextPageToken)
	if !ok {
		return nil, api.ErrNextTokenAssertionFailure
	}

	req = req.WithPage(string(token))
	return req, nil
}
