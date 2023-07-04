package gnews

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
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
	client.Params
}

func newQuery(apikey string) *Query {
	return &Query{
		Apikey: apikey,
		Params: client.NewParams(),
	}
}

func HandleHeadlines(apikey string, pf pageform.PageForm) (*Query, error) {
	data, ok := pf.(srv.GNewsHeadlines)
	if !ok {
		return nil, client.ErrTypeAssertionFailure
	}

	q, err := newQuery(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Keyword).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To)
	return q, nil
}

func HandleSearch(apikey string, pf pageform.PageForm) (*Query, error) {
	data, ok := pf.(srv.GNewsSearch)
	if !ok {
		return nil, client.ErrTypeAssertionFailure
	}

	q, err := newQuery(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Keyword).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithFrom(data.Form).
		WithTo(data.To)

	return q, nil
}

func (q *Query) WithKeywords(keyword string) *Query {
	q.Params.Set(qKeyword, keyword)
	return q
}

func (q *Query) WithCategory(category ...string) *Query {
	for _, c := range category {
		q.Params.Add(qCategory, c)
	}
	return q

}

func (q *Query) WithCountry(country ...string) *Query {
	for _, c := range country {
		q.Params.Add(qCountry, c)
	}
	return q
}

func (q *Query) WithLanguage(lang ...string) *Query {
	for _, l := range lang {
		q.Params.Add(qLanguage, l)
	}
	return q
}

func (q *Query) WithMaxArticles(i int) *Query {
	if i > API_MAX_NUM_ARTICLE {
		i = API_MAX_NUM_ARTICLE
	}
	if i < 0 {
		i = 10
	}
	q.Params.Set(qMax, strconv.Itoa(i))
	return q
}

func (q *Query) In(where ...string) *Query {
	for _, w := range where {
		q.Params.Add(qIn, w)
	}
	return q
}

func (q *Query) NullableIn(where ...string) *Query {
	for _, w := range where {
		q.Params.Add(qNullable, w)
	}
	return q
}

func (q *Query) SortByTime() *Query {
	q.Params.Set(qSortBy, "publishedAt")
	return q
}

func (q *Query) SortByRelevance() *Query {
	q.Params.Set(qSortBy, "relevance")
	return q
}

func (q *Query) withTime(t time.Time, format string, key string) *Query {
	q.Params.Set(key, t.Format(format))
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

func (q *Query) WithPage(n int) *Query {
	q.Params.Set(qPage, strconv.Itoa(n))
	return q
}

func (q *Query) WithExpand() *Query {
	q.Params.Set(q.Endpoint, "content")
	return q
}

func (q *Query) SetEndpoint(ep string) (*Query, error) {
	switch ep {
	case srv.EPTopHeadlines:
		q.Endpoint = EPTopHeadlines
	case srv.EPSearch:
		q.Endpoint = EPSearch
	default:
		return nil, errors.New("unknown endpoint")
	}
	return q, nil
}

func (q *Query) ToRequestURL(u *url.URL) string {
	v := q.Params.ToUrlVals()
	if q.Apikey != "" {
		v.Add("apikey", q.Apikey)
	}
	u = u.JoinPath(string(q.Endpoint))
	u.RawQuery = v.Encode()
	return u.String()
}
