package newsdata

import (
	"errors"
	"net/url"
	"time"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
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
	// qWithImage       = "image"        // currently not support
	// qWithVideo       = "video"        // currently not support
	// qWithFullContent = "full_content" // currently not support
)

type Query struct {
	Apikey   string
	Endpoint string
	cli.Params
	Page     string
	NextPage string
}

func newQuery(apikey string) *Query {
	return &Query{
		Apikey: apikey,
		Params: cli.NewParams(),
	}
}

func HandleLatestNewsQuery(apikey string, pf pageform.PageForm) (cli.Query, error) {
	data, ok := pf.(newsdata.NEWSDATAIOLatestNews)
	if !ok {
		return nil, cli.ErrTypeAssertionFailure
	}

	q, err := newQuery(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Keyword).
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	return q, nil
}

func HandleNewsArchive(apikey string, pf pageform.PageForm) (cli.Query, error) {
	data, ok := pf.(newsdata.NEWSDATAIONewsArchive)
	if !ok {
		return nil, cli.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTP()

	q, err := newQuery(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Keyword).
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To)

	return q, nil
}

func HandleNewsSources(apikey string, pf pageform.PageForm) (cli.Query, error) {
	data := pf.(newsdata.NEWSDATAIONewsSources)

	q, err := newQuery(apikey).SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	return q, nil
}

func (q *Query) WithKeywords(keyword string) *Query {
	q.Params.Set(qKeyword, keyword)
	return q
}

func (q *Query) WithKeywordsInTitle(keyword string) *Query {
	q.Params.Add(qKeywordInTitle, keyword)
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

func (q *Query) WithDomain(domain string) *Query {
	q.Params.Add(qDomain, domain)
	return q
}

func (q *Query) WithCategory(category ...string) *Query {
	for _, c := range category {
		q.Params.Add(qCategory, c)
	}
	return q
}

func (q *Query) WithNextPage(next string) *Query {
	q.NextPage = next
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

func (q *Query) ToBeautifulJSON(prefix, indent string) ([]byte, error) {
	return q.Params.ToBeautifulJSON(prefix, indent)
}

func (q *Query) ToJSON() ([]byte, error) {
	return q.Params.ToJSON()
}

func (q *Query) ToQueryString() string {
	return q.Params.ToQueryString()
}
