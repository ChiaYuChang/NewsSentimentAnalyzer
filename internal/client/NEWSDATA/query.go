package newsdata

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	val "github.com/go-playground/validator/v10"
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
	Params
	Page     string
	NextPage string
}
type Params map[string][]string

func newQuery(apikey string, ep string) *Query {
	return &Query{
		Apikey:   apikey,
		Endpoint: ep,
		Params:   make(Params),
		NextPage: "",
	}
}

type QueryBuilder struct {
	*Query
	val    *val.Validate
	errors []error
}

func NewQueryBuilder(apikey string, val *val.Validate) *QueryBuilder {
	return &QueryBuilder{
		Query:  newQuery(apikey, ""),
		val:    val,
		errors: nil,
	}
}

func (q *QueryBuilder) BuildLatestNewsQuery(apikey string, data *pageform.NEWSDATAIOLatestNews) (*Query, error) {
	q.SetEndpoint(data.Endpoint())
	if data.Keyword != "" {
		if err := q.val.Var(data.Keyword, "max=512"); err != nil {
			q.errors = append(q.errors, err)
		} else {
			q.Query = q.WithKeywords(data.Keyword)
		}
	}
	validateAndAppend(q, strings.Split(data.Domains, ","), qDomain, API_MAX_NUM_DOMAIN, VAL_TAG_DOMAIN)
	validateAndAppend(q, data.Language, qLanguage, API_MAX_NUM_LANGUAGE, VAL_TAG_LANGUAGE)
	validateAndAppend(q, data.Category, qCategory, API_MAX_NUM_CATEGORY, VAL_TAG_CATEGORY)
	validateAndAppend(q, data.Country, qCountry, API_MAX_NUM_COUNTRY, VAL_TAG_COUNTRY)

	if q.errors != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		for _, e := range q.errors {
			ecErr.WithDetails(e.Error())
		}
		return q.Query, ecErr
	}
	return q.Query, nil
}

func (q *QueryBuilder) BuildNewsArchive(apikey string, data *pageform.NEWSDATAIONewsArchive) (*Query, error) {
	q.SetEndpoint(data.Endpoint())
	if data.Keyword != "" {
		if err := q.val.Var(data.Keyword, "max=512"); err != nil {
			q.errors = append(q.errors, err)
		} else {
			q.Query = q.WithKeywords(data.Keyword)
		}
	}
	validateAndAppend(q, strings.Split(data.Domains, ","), qDomain, API_MAX_NUM_DOMAIN, VAL_TAG_DOMAIN)
	validateAndAppend(q, data.Language, qLanguage, API_MAX_NUM_LANGUAGE, VAL_TAG_LANGUAGE)
	validateAndAppend(q, data.Category, qCategory, API_MAX_NUM_CATEGORY, VAL_TAG_CATEGORY)
	validateAndAppend(q, data.Country, qCountry, API_MAX_NUM_COUNTRY, VAL_TAG_COUNTRY)

	if err := q.val.Struct(data.TimeRange); err != nil {
		q.errors = append(q.errors, err)
	}

	if q.errors != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		for _, e := range q.errors {
			ecErr.WithDetails(e.Error())
		}
		return q.Query, ecErr
	}
	return q.Query, nil
}

func (q *QueryBuilder) BuildNewsSources(apikey string, data *pageform.NEWSDATAIONewsSources) (*Query, error) {
	q.SetEndpoint(data.Endpoint())
	validateAndAppend(q, data.Language, qLanguage, API_MAX_NUM_LANGUAGE, VAL_TAG_LANGUAGE)
	validateAndAppend(q, data.Category, qCategory, API_MAX_NUM_CATEGORY, VAL_TAG_CATEGORY)
	validateAndAppend(q, data.Country, qCountry, API_MAX_NUM_COUNTRY, VAL_TAG_COUNTRY)

	if q.errors != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest)
		for _, e := range q.errors {
			ecErr.WithDetails(e.Error())
		}
		return q.Query, ecErr
	}
	return q.Query, nil
}

func validateAndAppend(q *QueryBuilder, field []string, which string, maxLen int, tag string) {
	if len(field) > 0 {
		if err := q.val.Var(
			field,
			fmt.Sprintf("max=%d", maxLen)); err != nil {
			q.errors = append(q.errors, err)
		} else {
			for _, f := range field {
				f = strings.TrimSpace(f)
				if err := q.val.Var(f, tag); err != nil {
					q.errors = append(q.errors, err)
				} else {
					q.Query.Params.Add(which, f)
				}
			}
		}
	}
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
		fmt.Printf("Add %s\n", strings.Join(v, ","))
		val.Add(k, strings.Join(v, ","))
	}
	return val
}

func (q *Query) WithKeywords(keyword string) *Query {
	q.Params.Add(qKeyword, keyword)
	return q
}

func (q *Query) WithKeywordsInTitle(keyword string) *Query {
	q.Params.Add(qKeywordInTitle, keyword)
	return q
}

func (q *Query) WithCountry(country string) *Query {
	q.Params.Add(qCountry, country)
	return q
}

func (q *Query) WithLanguage(lang string) *Query {
	q.Params.Add(qLanguage, lang)
	return q
}

func (q *Query) WithDomain(domain string) *Query {
	q.Params.Add(qDomain, domain)
	return q
}

func (q *Query) WithCategory(category string) *Query {
	q.Params.Add(qCategory, category)
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

func (q *Query) ToFrom(t time.Time) *Query {
	q.withTime(t, API_TIME_FORMAT, qToTime)
	return q
}

func (q *Query) SetEndpoint(ep string) (*Query, error) {
	switch ep {
	case pageform.NEWSDATAIOEPLatestNews:
		q.Endpoint = EPLatestNews
	case pageform.NEWSDATAIOEPNewsArchive:
		q.Endpoint = EPNewsArchive
	case pageform.NEWSDATAIOEPNewsSources:
		q.Endpoint = EPNewsSources
	case pageform.NEWSDATAEPCrypto:
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
