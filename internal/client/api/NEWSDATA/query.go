package newsdata

import (
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
)

const (
	qKeyword        = "q"
	qKeywordInTitle = "qInTitle"
	qCountry        = "country"
	qCategory       = "category"
	qLanguage       = "language"
	qDomain         = "domain"
	qFromTime       = "from"
	qToTime         = "to"
	qPage           = "page"
	// qWithImage       = "image"        // currently not support
	// qWithVideo       = "video"        // currently not support
	// qWithFullContent = "full_content" // currently not support
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

// Keywords or phrases to search for in the article body.
func (q *Query) WithKeywords(keyword string) *Query {
	q.params.Set(qKeyword, keyword)
	return q
}

// Keywords or phrases to search for in the article title.
func (q *Query) WithKeywordsInTitle(keyword string) *Query {
	q.params.Add(qKeywordInTitle, keyword)
	return q
}

// The 2-letter code of the country you want to get headlines for.
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

// The domains to restrict the search to.
func (q *Query) WithDomain(domain string) *Query {
	q.params.Add(qDomain, domain)
	return q
}

// Find sources that display news of this category.
func (q *Query) WithCategory(category ...string) *Query {
	for _, c := range category {
		q.params.Add(qCategory, c)
	}
	return q
}

// Token for retrieving next page.
func (q *Query) WithNextPage(page string) *Query {
	q.params.Set(qPage, page)
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
	case newsdata.EPLatestNews:
		q.Endpoint = EPLatestNews
	case newsdata.EPNewsArchive:
		q.Endpoint = EPNewsArchive
	case newsdata.EPNewsSources:
		q.Endpoint = EPNewsSources
	case newsdata.EPCrypto:
		q.Endpoint = EPCrypto
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
