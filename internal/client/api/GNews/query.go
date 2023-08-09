package gnews

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/GNews"
)

const (
	qKeyword  = "q"
	qLanguage = "lang"
	qCountry  = "country"
	qCategory = "category"
	qMax      = "max"
	qIn       = "in"
	qNullable = "nullable"
	qFromTime = "from"
	qToTime   = "to"
	qSortBy   = "sortby"
)

// These parameters will only work for paid account
const (
	qPage   = "page"
	qExpand = "expand"
)

type Query struct {
	Apikey   string
	Endpoint string
	Page     int
	params   api.Params
}

func newQuery(apikey string) *Query {
	return &Query{
		Apikey: apikey,
		params: api.NewParams(),
	}
}

// Append keywords to the query object
func (q *Query) WithKeywords(keyword string) *Query {
	q.params.Set(qKeyword, keyword)
	return q
}

// Append categories to the query object
func (q *Query) WithCategory(category ...string) *Query {
	for _, c := range category {
		q.params.Add(qCategory, c)
	}
	return q

}

// Append countries to the query object
func (q *Query) WithCountry(country ...string) *Query {
	for _, c := range country {
		q.params.Add(qCountry, c)
	}
	return q
}

// Append languages to the query object
func (q *Query) WithLanguage(lang ...string) *Query {
	for _, l := range lang {
		q.params.Add(qLanguage, l)
	}
	return q
}

// set the maximum number of articles by a single query (max = 100)
func (q *Query) WithMaxArticles(i int) *Query {
	if i > API_MAX_NUM_ARTICLE {
		i = API_MAX_NUM_ARTICLE
	}
	if i < 0 {
		i = 10
	}
	q.params.Set(qMax, strconv.Itoa(i))
	return q
}

// the field to search
func (q *Query) In(where ...string) *Query {
	for _, w := range where {
		q.params.Add(qIn, w)
	}
	return q
}

// allow null values in certain fields
func (q *Query) NullableIn(where ...string) *Query {
	for _, w := range where {
		q.params.Add(qNullable, w)
	}
	return q
}

// sort the query result by time
func (q *Query) SortByTime() *Query {
	q.params.Set(qSortBy, "publishedAt")
	return q
}

// sort the query result by relevance
func (q *Query) SortByRelevance() *Query {
	q.params.Set(qSortBy, "relevance")
	return q
}

// a helper function to set time parameter
func (q *Query) withTime(t time.Time, format string, key string) *Query {
	if t.After(API_MIN_TIME) {
		q.params.Set(key, t.Format(format))
	}
	return q
}

func (q *Query) WithFrom(t time.Time) *Query {
	q.withTime(t, API_TIME_FORMAT, qFromTime)
	return q
}

func (q *Query) WithTo(t time.Time) *Query {
	q.withTime(t, API_TIME_FORMAT, qToTime)
	return q
}

// Only for paid user. Set the number of articles in a single page.
func (q *Query) WithPage(n int) *Query {
	q.Page = n
	return q
}

// Only for paid user. Should the content contains full text.
func (q *Query) WithExpand() *Query {
	q.params.Set(q.Endpoint, "content")
	return q
}

// Set endpoints
func (q *Query) SetEndpoint(ep string) (*Query, error) {
	switch ep {
	case srv.EPTopHeadlines:
		q.Endpoint = EPTopHeadlines
	case srv.EPSearch:
		q.Endpoint = EPSearch
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
