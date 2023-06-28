package newsapi

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/code"
)

type Endpoint string

func (ep Endpoint) IsDefault() bool {
	return ep == API_DEFAULT_ENDPOINT
}

type SearchIn string
type Category string

func (c Category) IsEmpty() bool {
	return c == ""
}

type SortBy string

var API_HOST = "https://newsapi.org"
var API_VERSION = "v2"
var API_METHOD = http.MethodGet
var API_URL, _ = url.Parse(fmt.Sprintf("%s/%s/", API_HOST, API_VERSION))

const API_TIME_FORMAT = "2006-01-02T15:04:05Z"
const API_DEFAULT_PAGE_SIZE = 100
const API_MAX_PAGE_SIZE = 100
const API_DEFAULT_PAGE = 1
const API_MAX_SOURCES_NUM = 30
const API_TIME_DELAY_FOR_FREE_PLAN = 24 * 60 * time.Minute

const (
	EPEverything Endpoint = "everything"
	EPHeadlines  Endpoint = "top-headlines"
	EPSources    Endpoint = "top-headlines/sources"
)
const API_DEFAULT_ENDPOINT = EPEverything

const (
	CGeneral      Category = "general"
	CBusiness     Category = "business"
	CEntertaiment Category = "entertaiment"
	CHealth       Category = "health"
	CScience      Category = "science"
	CSports       Category = "sport"
	CTechnology   Category = "technology"
)

const (
	InTitle       SearchIn = "title"
	InDescription SearchIn = "description"
	InContent     SearchIn = "content"
)

const (
	ByPublishedAt SortBy = "publishedAt" // default
	ByRelevancy   SortBy = "relevancy"
	ByPopularity  SortBy = "popularity"
)
const API_DEFAULT_SORTBY = ByPublishedAt

type Client struct {
	ApiKey string
	*http.Client
}

// func NewDefaultClient() Client {
// 	// using default client
// 	return NewClient(API_KEY, http.DefaultClient)
// }

func NewClient(apiKey string, cli *http.Client) Client {
	return Client{apiKey, cli}
}

func (cli *Client) SetAPIKey(apiKey string) {
	cli.ApiKey = apiKey
}

func (cli Client) NewQuery(k string, ep Endpoint, params ...Params) *Query {
	q := &Query{Keyword: k, Endpoint: ep, Params: map[string]Params{}}
	return q.AppendParams(params...)
}

func (cli Client) NewQueryWithDefaultVals() *Query {
	return cli.NewQuery("", API_DEFAULT_ENDPOINT)
}

func (cli Client) NewPager(pageSize, page int) Pager {
	return Pager{PageSize: pageSize, Page: page}
}

func (cli Client) NewPagerWithDefaultVals() Pager {
	return cli.NewPager(API_MAX_PAGE_SIZE, API_DEFAULT_PAGE)
}

func (cli Client) NewNewsSources(srcs ...string) NewsSources {
	return NewsSources(srcs)
}

func (cli Client) NewTopHeadlinesParams(country code.CountryCode, category Category) TopHeadlinesParams {
	return TopHeadlinesParams{Country: country, Category: category}
}

func (cli Client) NewTopHeadlinesParamsWithDefaultVals() TopHeadlinesParams {
	return cli.NewTopHeadlinesParams("", "")
}

func (cli Client) NewEverythingParams(searchIn []SearchIn, includeDomains, excludeDomains []string,
	from, to time.Time, language code.Language, sortBy SortBy) EverythingParams {
	params := EverythingParams{
		SearchIn:       make([]string, len(searchIn)),
		Domains:        includeDomains,
		ExcludeDomains: excludeDomains,
		From:           from,
		To:             to,
		Language:       language,
		SortedBy:       sortBy,
	}

	for i, si := range searchIn {
		params.SearchIn[i] = string(si)
	}
	return params
}

func (cli Client) NewEverythingParamsWithDefaultVals() EverythingParams {
	return cli.NewEverythingParams(
		nil, nil, nil,
		*new(time.Time),
		*new(time.Time),
		"",
		API_DEFAULT_SORTBY,
	)
}
