package newsapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/newsapi"
)

const (
	qKeyword        = "q"
	qSearchIn       = "searchIn"
	qDomains        = "domains"
	qExcludeDomains = "excludeDomains"
	qCountry        = "country"
	qCategory       = "category"
	qLanguage       = "language"
	qSources        = "sources"
	qSortBy         = "sortBy"
	qPageSize       = "pageSize"
	qPage           = "page"
	qFromTime       = "from"
	qToTime         = "to"
)

type Query struct {
	Apikey   string
	Endpoint string
	params   api.Params
}

func newQuery(apikey string) *Query {
	return &Query{
		Apikey: apikey,
		params: api.NewParams(),
	}
}

// Keywords or phrases to search for in the article title and body.
func (q *Query) WithKeywords(keyword string) *Query {
	q.params.Set(qKeyword, keyword)
	return q
}

// The domains to restrict the search to.
func (q *Query) WithDomains(domain ...string) *Query {
	for _, d := range domain {
		q.params.Add(qDomains, d)
	}
	return q
}

// The domains to remove from the results.
func (q *Query) WithExcludeDomains(domain ...string) *Query {
	for _, d := range domain {
		q.params.Add(qExcludeDomains, d)
	}
	return q
}

func (q *Query) WithSearchIn(searchIn string) *Query {
	q.params.Set(qSearchIn, searchIn)
	return q
}

// The fields to restrict your keywords search to. (Default: title,description,content)
func (q *Query) SearchInTitle() *Query {
	q.params.Add(qSearchIn, "title")
	return q
}

// The fields to restrict your keywords search to. (Default: title,description,content)
func (q *Query) SearchInDescription() *Query {
	q.params.Add(qSearchIn, "description")
	return q
}

// The fields to restrict your keywords search to. (Default: title,description,content)
func (q *Query) SearchInContent() *Query {
	q.params.Add(qSearchIn, "content")
	return q
}

// Set up the order to sort the articles in. (default: by published at )
func (q *Query) SortByRelevancy() *Query {
	q.params.Set(qSortBy, "relevancy")
	return q
}

// Set up the order to sort the articles in. (default: by published at )
func (q *Query) SortByPopularity() *Query {
	q.params.Set(qSortBy, "popularity")
	return q
}

// Set up the order to sort the articles in. (default: by published at )
func (q *Query) SortByPublishedAt() *Query {
	q.params.Set(qSortBy, "publishedAt")
	return q
}

// The number of results to return per page. (Default: 100, Max: 100)
func (q *Query) WithPageSize(ps int) *Query {
	if ps <= 0 || ps > API_MAX_PAGE_SIZE {
		ps = API_MAX_PAGE_SIZE
	}
	q.params.Set(qPageSize, strconv.Itoa(ps))
	return q
}

// The 2-letter ISO 3166-1 code of the country you want to get headlines for.
func (q *Query) WithCountry(country ...string) *Query {
	for _, c := range country {
		q.params.Add(qCountry, c)
	}
	return q
}

// Find sources that display news in a specific language
func (q *Query) WithLanguage(lang ...string) *Query {
	for _, l := range lang {
		q.params.Add(qLanguage, l)
	}
	return q
}

// Find sources that display news of this category.
func (q *Query) WithCategory(category ...string) *Query {
	for _, c := range category {
		q.params.Add(qCategory, c)
	}
	return q
}

// Identifiers (maximum 20) for the news sources or blogs you want headlines from.
func (q *Query) WithSources(src ...string) *Query {
	for _, s := range src {
		q.params.Add(qSources, s)
	}
	return q
}

func (q *Query) withTime(t time.Time, format string, key string) *Query {
	if t.After(API_MIN_TIME) {
		q.params.Set(key, t.Format(format))
	}
	return q
}

// A date and optional time for the oldest article allowed.
func (q *Query) WithFrom(t time.Time) *Query {
	q.withTime(t, API_TIME_FORMAT, qFromTime)
	return q
}

// A date and optional time for the newest article allowed.
func (q *Query) WithTo(t time.Time) *Query {
	q.withTime(t, API_TIME_FORMAT, qToTime)
	return q
}

func (q *Query) SetEndpoint(ep string) (*Query, error) {
	switch ep {
	case newsapi.EPEverything:
		q.Endpoint = EPEverything
	case newsapi.EPTopHeadlines:
		q.Endpoint = EPTopHeadlines
	case newsapi.EPSources:
		return nil, client.ErrNotSupportedEndpoint
	default:
		return nil, client.ErrUnknownEndpoint
	}
	return q, nil
}

func (q *Query) Params() api.Params {
	return q.params
}

func (q *Query) String() string {
	return q.params.ToQueryString()
}

func (q *Query) ToRequest() (*http.Request, error) {
	return api.ToRequest(API_URL, API_METHOD, q.Apikey, q.Endpoint, nil, q)
}
