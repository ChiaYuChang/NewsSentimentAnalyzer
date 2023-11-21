package newsapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"
	"github.com/google/uuid"
)

const (
	Keyword        api.Key = "q"
	SearchIn       api.Key = "searchIn"
	Domains        api.Key = "domains"
	ExcludeDomains api.Key = "excludeDomains"
	Country        api.Key = "country"
	Category       api.Key = "category"
	Language       api.Key = "language"
	Sources        api.Key = "sources"
	SortBy         api.Key = "sortBy"
	PageSize       api.Key = "pageSize"
	Page           api.Key = "page"
	FromTime       api.Key = "from"
	ToTime         api.Key = "to"
	APIKey         api.Key = "apikey"
)

const (
	AuthorizationHeader = "X-Api-Key"
)

type SearchInField string

const (
	SearchInTitle       SearchInField = "title"
	SearchInDescription SearchInField = "description"
	SearchInContent     SearchInField = "content"
)

func (f *SearchInField) Parse(s pageform.SearchIn) {
	switch s.String() {
	default:
		return
	case "in-title":
		*f = SearchInTitle
	case "in-description":
		*f = SearchInDescription
	case "in-content":
		*f = SearchInContent
	}
}

type SoryByField string

const (
	SortByRelevancy   SearchInField = "relevancy"
	SortByPopularity  SearchInField = "popularity"
	SortByPublishedAt SearchInField = "publishedAt"
)

type Request struct {
	*api.RequestProto
	Page int
}

func NewRequest(apikey string) *Request {
	r := api.NewRequestProtoType(srv.API_NAME, ",")
	r.SetApiKey(apikey)

	return &Request{RequestProto: r}
}

// Keywords or phrases to search for in the article title and body.
func (r *Request) WithKeywords(keyword string) *Request {
	r.Set(Keyword, keyword)
	return r
}

// The domains to restrict the search to.
func (r *Request) WithDomains(domain ...string) *Request {
	for _, d := range domain {
		r.Add(Domains, d)
	}
	return r
}

// The domains to remove from the results.
func (r *Request) WithExcludeDomains(domain ...string) *Request {
	for _, d := range domain {
		r.Add(ExcludeDomains, d)
	}
	return r
}

func (r *Request) WithSearchIn(searchIn SearchInField) *Request {
	r.Set(SearchIn, string(searchIn))
	return r
}

// The fields to restrict your keywords search to. (Default: title,description,content)
func (r *Request) SearchInTitle() *Request {
	r.Add(SearchIn, string(SearchInTitle))
	return r
}

// The fields to restrict your keywords search to. (Default: title,description,content)
func (r *Request) SearchInDescription() *Request {
	r.Add(SearchIn, string(SearchInDescription))
	return r
}

// The fields to restrict your keywords search to. (Default: title,description,content)
func (r *Request) SearchInContent() *Request {
	r.Add(SearchIn, string(SearchInContent))
	return r
}

func (r *Request) SortBy(which SoryByField) *Request {
	r.Set(SortBy, string(which))
	return r
}

// Set up the order to sort the articles in. (default: by published at )
func (r *Request) SortByRelevancy() *Request {
	r.Set(SortBy, string(SortByRelevancy))
	return r
}

// Set up the order to sort the articles in. (default: by published at )
func (r *Request) SortByPopularity() *Request {
	r.Set(SortBy, string(SortByPopularity))
	return r
}

// Set up the order to sort the articles in. (default: by published at )
func (r *Request) SortByPublishedAt() *Request {
	r.Set(SortBy, string(SortByPublishedAt))
	return r
}

// The number of results to return per page. (Default: 100, Max: 100)
func (r *Request) WithPageSize(ps int) *Request {
	if ps <= 0 || ps > API_MAX_PAGE_SIZE {
		ps = API_MAX_PAGE_SIZE
	}
	r.Set(PageSize, strconv.Itoa(ps))
	return r
}

func (r *Request) WithPage(page int) *Request {
	r.Page = page
	return r
}

// The 2-letter ISO 3166-1 code of the country you want to get headlines for.
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

// Find sources that display news of this category.
func (r *Request) WithCategory(category ...string) *Request {
	for _, c := range category {
		r.Add(Category, c)
	}
	return r
}

// Identifiers (maximum 20) for the news sources or blogs you want headlines from.
func (r *Request) WithSources(src ...string) *Request {
	for _, s := range src {
		r.Add(Sources, s)
	}
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
	switch ep {
	case srv.EPEverything, EPEverything:
		r.RequestProto.SetEndpoint(EPEverything)
	case srv.EPTopHeadlines, EPTopHeadlines:
		r.RequestProto.SetEndpoint(EPTopHeadlines)
	case srv.EPSources:
		return nil, client.ErrNotSupportedEndpoint
	default:
		return nil, client.ErrUnknownEndpoint
	}
	return r, nil
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

	httpReq.Header.Set(AuthorizationHeader, req.APIKey())
	if req.Page > 1 {
		// default: page=1
		p.Set(Page, strconv.Itoa(req.Page))
	}
	httpReq.URL.RawQuery = p.Encode()
	return httpReq, nil
}

func (req Request) ToPreviewCache(uid uuid.UUID) (cKey string, c *api.PreviewCache) {
	return req.RequestProto.ToPreviewCache(uid, api.IntNextPageToken(1), nil)
}

func RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	if cq.NextPage.Equal(api.IntLastPageToken) {
		return nil, api.ErrEndOfQuery
	}

	var err error
	req := NewRequest(cq.API.Key)
	_, err = req.SetEndpoint(cq.API.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error while set endpoint: %w", err)
	}

	req.Values, err = url.ParseQuery(cq.RawQuery)
	if err != nil {
		return nil, err
	}

	i, ok := cq.NextPage.(api.IntNextPageToken)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}
	req = req.WithPage(int(i))
	return req, nil
}
