package newsdata

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/code"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type Query struct {
	Endpoint       Endpoint
	Keyword        string
	KeywordInTitle string
	Country        []code.CountryCode
	Category       []Category
	Language       []code.Language
	Params         map[string]Params
	Page           string
}

func (q *Query) WithKeywords(keyword string) *Query {
	q.Keyword = keyword
	return q
}

func (q *Query) WithKeywordsInTitle(keyword string) *Query {
	q.KeywordInTitle = keyword
	return q
}

func (q *Query) WithCountry(country ...code.CountryCode) *Query {
	q.Country = append(q.Country, country...)
	return q
}

func (q *Query) WithLanguage(lang ...code.Language) *Query {
	q.Language = append(q.Language, lang...)
	return q
}

func (q *Query) AppendParams(params ...Params) *Query {
	for _, p := range params {
		q.Params[p.ParamsName()] = p
	}
	return q
}

func (q *Query) ToHTTPRequest(ctx context.Context, apiKey string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx, API_METHOD,
		fmt.Sprintf("%s/%s/%s", API_URL, API_VERSION, q.Endpoint), nil)

	if err != nil {
		return nil, fmt.Errorf("error while new request with context: %w", err)
	}

	_, err = q.appendAPIKeyToReq(req, apiKey)
	if err != nil {
		return nil, err
	}
	return q.appendParamsToReq(req)
}

func (q *Query) appendAPIKeyToReq(req *http.Request, apiKey string) (*http.Request, error) {
	if apiKey == "" {
		err := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
		err.WithDetails("API key is missing")
		return nil, err
	}
	req.Header.Add("X-ACCESS-KEY", apiKey)
	return req, nil
}

func (q *Query) appendParamsToReq(req *http.Request) (*http.Request, error) {
	var err error
	query := make(url.Values)

	if q.Keyword != "" {
		query.Add("q", q.Keyword)
	}

	for _, p := range q.Params {
		fmt.Printf("Append %s parameters...", p.ParamsName())
		query, err = p.ToUrlVals(query)
		if err != nil {
			fmt.Println("Failed")
			return nil, err
		}
		fmt.Println("OK")
	}

	req.URL.RawQuery = query.Encode()
	return req, nil
}

type Params interface {
	ToUrlVals(vals url.Values) (url.Values, error)
	ParamsName() string
}

type ContentFilter struct {
	HasFullContent bool
	HasImage       bool
	HasVideo       bool
}

func (params ContentFilter) ToUrlVals(vals url.Values) (url.Values, error) {
	if params.HasFullContent {
		vals.Add("full_content", "1")
	}

	if params.HasImage {
		vals.Add("image", "1")
	}

	if params.HasVideo {
		vals.Add("video", "1")
	}
	return vals, nil
}

func (params ContentFilter) ParamsName() string {
	return "Content-Filter"
}

type DomainFilter []Domain

func (params DomainFilter) ToUrlVals(vals url.Values) (url.Values, error) {
	var ds []string
	if l := len(params); l < API_MAX_NUM_DOMAIN {
		ecErr := ec.MustGetErr(ec.ECBadRequest).(*ec.Error)
		ecErr.WithDetails(fmt.Sprintf(
			"number of domain in a query should be less than or equal to 5, but get %d", l,
		))
		return nil, ecErr
	} else {
		ds = make([]string, len(params))
	}

	for i, d := range params {
		ds[i] = string(d)
	}

	if len(params) > 0 {
		vals.Add("domain", strings.Join(ds, ","))
	}
	return vals, nil
}

func (Params DomainFilter) ParamsName() string {
	return "Domain"
}

type ArchiveParams struct {
	From time.Time
	To   time.Time
}

func (params ArchiveParams) ToUrlVals(vals url.Values) (url.Values, error) {
	if !params.From.IsZero() {
		vals.Add("from_date", params.From.Format(API_TIME_FORMAT))
	}

	if !params.To.IsZero() {
		vals.Add("to_date", params.From.Format(API_TIME_FORMAT))
	}

	return vals, nil
}

func (params ArchiveParams) ParamsName() string {
	return "Archive-Time"
}
