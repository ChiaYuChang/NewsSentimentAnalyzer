package gnews

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"
	"github.com/google/uuid"
)

// query parameters
const (
	Keyword  api.Key = "q"
	Language api.Key = "lang"
	Country  api.Key = "country"
	Category api.Key = "category"
	Max      api.Key = "max"
	In       api.Key = "in"
	Nullable api.Key = "nullable"
	FromTime api.Key = "from"
	ToTime   api.Key = "to"
	SortBy   api.Key = "sortby"
	APIKey   api.Key = "apikey"
)

// These parameters will only work for paid account
const (
	Page   api.Key = "page"
	Expand api.Key = "expand"
)

type Request struct {
	*api.RequestProto
	Page api.IntNextPageToken
}

func NewRequest(apikey string) *Request {
	r := api.NewRequestProtoType(srv.API_NAME, ",")
	r.SetApiKey(apikey)

	return &Request{RequestProto: r}
}

// Append keywords to the query object
func (r *Request) WithKeywords(keyword string) *Request {
	r.RequestProto.Set(Keyword, keyword)
	return r
}

// Append categories to the query object
func (r *Request) WithCategory(category ...string) *Request {
	for _, c := range category {
		r.RequestProto.Add(Category, c)
	}
	return r
}

// Append countries to the query object
func (r *Request) WithCountry(country ...string) *Request {
	for _, c := range country {
		r.RequestProto.Add(Country, c)
	}
	return r
}

// Append languages to the query object
func (r *Request) WithLanguage(lang ...string) *Request {
	for _, l := range lang {
		r.RequestProto.Add(Language, l)
	}
	return r
}

// set the maximum number of articles by a single query (max = 100)
func (r *Request) WithMaxArticles(i int) *Request {
	if i > API_MAX_NUM_ARTICLE {
		i = API_MAX_NUM_ARTICLE
	}
	if i < 0 {
		i = 10
	}
	r.Set(Max, strconv.Itoa(i))
	return r
}

// the field to search
func (r *Request) In(where ...string) *Request {
	for _, w := range where {
		r.Add(In, w)
	}
	return r
}

// allow null values in certain fields
func (r *Request) NullableIn(where ...string) *Request {
	for _, w := range where {
		r.Add(Nullable, w)
	}
	return r
}

// sort the query result by time
func (r *Request) SortByTime() *Request {
	r.Set(SortBy, "publishedAt")
	return r
}

// sort the query result by relevance
func (r *Request) SortByRelevance() *Request {
	r.Set(SortBy, "relevance")
	return r
}

// a helper function to set time parameter
func (r *Request) withTime(t time.Time, format string, key api.Key) *Request {
	if t.After(API_MIN_TIME) {
		r.Set(key, t.Format(format))
	}
	return r
}

func (r *Request) WithFrom(t time.Time) *Request {
	r.withTime(t, API_TIME_FORMAT, FromTime)
	return r
}

func (r *Request) WithTo(t time.Time) *Request {
	r.withTime(t, API_TIME_FORMAT, ToTime)
	return r
}

// Only for paid user. Set the number of articles in a single page.
func (r *Request) WithPage(n int) *Request {
	if n > 1 {
		r.Page = api.IntNextPageToken(n)
	}
	return r
}

// Only for paid user. Should the content contains full text.
func (r *Request) WithExpand() *Request {
	r.Set(Expand, "content")
	return r
}

// Set endpoints
func (r *Request) SetEndpoint(ep string) (*Request, error) {
	switch ep {
	case srv.EPTopHeadlines, EPTopHeadlines:
		r.RequestProto.SetEndpoint(EPTopHeadlines)
	case srv.EPSearch, EPSearch:
		r.RequestProto.SetEndpoint(EPSearch)
	default:
		return nil, client.ErrUnknownEndpoint
	}
	return r, nil
}

// generate a http.Request
func (r *Request) ToHttpRequest() (*http.Request, error) {
	httpReq, err := r.RequestProto.ToHTTPRequest(API_URL, API_METHOD, nil)
	if err != nil {
		return nil, err
	}

	p, err := r.Params.Clone()
	if err != nil {
		return nil, err
	}

	p.Set(APIKey, r.APIKey())
	if r.Page > 1 {
		p.Set(Page, strconv.Itoa(int(r.Page)))
	}
	httpReq.URL.RawQuery = p.Encode()
	return httpReq, nil
}

func (r Request) ToPreviewCache(uid uuid.UUID) (cKey string, c *api.PreviewCache) {
	if r.Page == 0 {
		return r.RequestProto.ToPreviewCache(uid, api.IntNextPageToken(1), nil)
	}
	return r.RequestProto.ToPreviewCache(uid, r.Page, nil)
}

func RequestFromPreviewCache(c *api.PreviewCache) (api.Request, error) {
	if c.Query.NextPage.Equal(api.IntLastPageToken) {
		return nil, api.ErrNotNextPage
	}

	var err error
	req := NewRequest(c.Query.API.Key)
	_, err = req.SetEndpoint(c.Query.API.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error while set endpoint: %w", err)
	}

	req.Values, err = url.ParseQuery(c.Query.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("error while parsing raw query: %w", err)
	}

	token, ok := c.Query.NextPage.(api.IntNextPageToken)
	if !ok {
		return nil, api.ErrNextTokenAssertionFailure
	}

	req = req.WithPage(int(token))
	return req, nil
}
