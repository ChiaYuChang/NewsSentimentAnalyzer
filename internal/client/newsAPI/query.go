package newsapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/code"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type Query struct {
	Keyword  string
	Endpoint Endpoint
	Params   map[string]Params
}

func (q *Query) WithKeywords(keyword string) *Query {
	q.Keyword = keyword
	return q
}

func (q *Query) SetEndPoint(ep Endpoint) *Query {
	q.Endpoint = ep
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
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

type Pager struct {
	Page     int
	PageSize int
}

func (p Pager) ToUrlVals(vals url.Values) (url.Values, error) {
	if p.Page > 1 {
		vals.Add("page", strconv.Itoa(p.Page))
	}

	if p.PageSize > 0 {
		if p.PageSize > API_MAX_PAGE_SIZE {
			err := ec.MustGetErr(ec.ECBadRequest).(*ec.Error)
			err.WithDetails(fmt.Sprintf(
				"page size should be less than %d, but get %d",
				API_MAX_PAGE_SIZE, p.PageSize,
			))
			return vals, err
		}
		vals.Add("pageSize", strconv.Itoa(p.PageSize))
	}
	return vals, nil
}

func (p Pager) ParamsName() string {
	return "Pager"
}

type NewsSources []string

func (src NewsSources) ToUrlVals(vals url.Values) (url.Values, error) {
	if len(src) > 0 {
		if len(src) > API_MAX_SOURCES_NUM {
			err := ec.MustGetErr(ec.ECBadRequest).(*ec.Error)
			err.WithDetails(
				fmt.Sprintf(
					"too many sources in a single request, max: %d, get: %d",
					API_MAX_SOURCES_NUM, len(src)),
			)
			return nil, err
		}
		vals.Add("sources", strings.Join(src, ","))
	}
	return vals, nil
}

func (src NewsSources) ParamsName() string {
	return "Sources"
}

type TopHeadlinesParams struct {
	Country  code.CountryCode
	Category Category
}

func (param TopHeadlinesParams) ToUrlVals(vals url.Values) (url.Values, error) {
	if vals.Get("sources") != "" && (!param.Country.IsEmpty() || !param.Category.IsEmpty()) {
		err := ec.MustGetErr(ec.ECBadRequest).(*ec.Error)
		err.WithDetails("mixing sources with the country or category params is not supported")
		return nil, err
	}

	if !param.Category.IsEmpty() {
		vals.Add("country", string(param.Country))
	}

	if !param.Category.IsEmpty() {
		vals.Add("category", string(param.Category))
	}
	return vals, nil
}

func (param TopHeadlinesParams) ParamsName() string {
	return "Top-Headlines"
}

type EverythingParams struct {
	SearchIn       []string
	Domains        []string
	ExcludeDomains []string
	From           time.Time
	To             time.Time
	Language       code.Language
	SortedBy       SortBy
}

func (param EverythingParams) ToUrlVals(vals url.Values) (url.Values, error) {
	if len(param.SearchIn) > 0 {
		vals.Add("searchIn", strings.Join(param.SearchIn, ","))
	}

	if len(param.Domains) > 0 {
		vals.Add("domains", strings.Join(param.Domains, ","))
	}

	if len(param.ExcludeDomains) > 0 {
		vals.Add("excludeDomains", strings.Join(param.ExcludeDomains, ","))
	}

	if !param.From.IsZero() {
		vals.Add("from", param.From.Format(API_TIME_FORMAT))
	}

	if !param.To.IsZero() {
		vals.Add("to", param.To.Format(API_TIME_FORMAT))
	}

	if !param.Language.IsEmpty() {
		vals.Add("language", string(param.Language))
	}

	if param.SortedBy != API_DEFAULT_SORTBY {
		vals.Add("sortBy", string(param.SortedBy))
	}

	return vals, nil
}

func (param EverythingParams) ParamsName() string {
	return "Everything"
}
