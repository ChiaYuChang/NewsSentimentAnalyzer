package newsdata

import (
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/code"
)

var API_URL = "https://newsdata.io/api"
var API_VERSION = "1"
var API_METHOD = http.MethodGet

var API_KEY string = ""

type Endpoint string
type Category string
type Domain string

const API_TIME_FORMAT = "2006-01-02"
const (
	API_MAX_NUM_DOMAIN   = 5
	API_MAX_NUM_COUNTRY  = 5
	API_MAX_NUM_CATEGORY = 5
)

const (
	EPLatestNews Endpoint = "news"
	EPCrypto     Endpoint = "crypto"
	EPArchive    Endpoint = "archive"
	EPSources    Endpoint = "sources"
)
const API_DEFAULT_ENDPOINT = EPLatestNews

const (
	CBusiness      Category = "business"
	CEntertainment Category = "entertainment"
	CEnvironment   Category = "environment"
	CFood          Category = "food"
	CHealth        Category = "health"
	CPolitics      Category = "politics"
	CScience       Category = "science"
	CSports        Category = "sports"
	CTechnology    Category = "technology"
	CTop           Category = "top"
	CTourism       Category = "tourism"
	CWorld         Category = "world"
)

const (
	DGVM          Domain = "gvm"      // 遠見雜誌
	DNewtalk      Domain = "newtalk"  // 新頭殼
	DTechNews     Domain = "technews" // 科技新報
	DLibertyTimes Domain = "ltn"      // 自由時報
	DETtoday      Domain = "ettoday"  // ETtoday新聞雲
	DGNN          Domain = "gnn"      // GNN 新聞網- 巴哈姆特
)

func SetDefaultAPIKey(key string) {
	API_KEY = key
}

type Client struct {
	ApiKey string
	*http.Client
}

func NewDefaultClient() Client {
	// using default client
	return NewClient(API_KEY, http.DefaultClient)
}

func NewClient(apiKey string, cli *http.Client) Client {
	return Client{apiKey, cli}
}

func (cli *Client) SetAPIKey(apikeys string) {
	cli.ApiKey = apikeys
}

func (cli Client) NewQuery(ep Endpoint, keyword, keywordInTitle string,
	country []code.CountryCode, category []Category, language []code.Language,
	page string, params ...Params) *Query {
	q := &Query{
		Endpoint:       ep,
		Keyword:        keyword,
		KeywordInTitle: keywordInTitle,
		Country:        country,
		Category:       category,
		Language:       language,
		Params:         map[string]Params{},
		Page:           page,
	}
	return q.AppendParams(params...)
}

func (cli Client) NewQueryWithDefaultVals() *Query {
	return cli.NewQuery(
		API_DEFAULT_ENDPOINT,
		"", "",
		nil, nil, nil, "",
	)
}

func (cli Client) NewContentFilter(hasFullContent, hasImage, hasVideo bool) ContentFilter {
	return ContentFilter{hasFullContent, hasImage, hasVideo}
}

func (cli Client) NewDomainFilter(domains ...Domain) []Domain {
	return domains
}

func (cli Client) NewArchiveParams(from, to time.Time) ArchiveParams {
	return ArchiveParams{from, to}
}

func (cli Client) NewArchiveParamsWithDefaultVals() ArchiveParams {
	return cli.NewArchiveParams(*new(time.Time), *new(time.Time))
}
