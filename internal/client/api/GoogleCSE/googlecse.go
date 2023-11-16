package googlecse

import (
	"fmt"
	"net/http"
	"time"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GoogleCSE"
)

func init() {
	cli.RegisterPageForm(
		srv.GoogleCSE{},
		CSEHandler{},
		EPform2client)
}

const (
	API_SCHEME          = "https"
	API_HOST            = "www.googleapis.com"
	API_PATH            = "customsearch"
	API_VERSION         = "v1"
	API_METHOD          = http.MethodGet
	API_MAX_NUM_ARTICLE = 100
	API_TIME_FORMAT     = "2006-01-02T15:04:05Z"
	API_RESP_TIME_FMT   = "2006-01-02T15:04:05Z"
)

var API_MIN_TIME, _ = time.Parse(time.DateOnly, "1900-01-01")
var API_URL = fmt.Sprintf("%s://%s/%s/%s", API_SCHEME, API_HOST, API_PATH, API_VERSION)

const (
	DEFAULT_PAGE_SIZE = 10
)

const (
	EPCustomSearch   string = ""
	EPSiteRestricted string = "siterestrict"
)

var EPform2client = map[string]string{
	srv.EPCustomSearch:   EPCustomSearch,
	srv.EPSiteRestricted: EPSiteRestricted,
}
