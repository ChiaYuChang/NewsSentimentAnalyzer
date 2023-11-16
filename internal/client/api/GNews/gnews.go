package gnews

import (
	"fmt"
	"net/http"
	"time"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"
)

func init() {
	cli.RegisterPageForm(
		srv.GNewsHeadlines{},
		TopHeadlinesHandler{},
		EPform2client)

	cli.RegisterPageForm(
		srv.GNewsSearch{},
		SearchHandler{},
		EPform2client)
}

const (
	API_SCHEME          = "https"
	API_HOST            = "gnews.io"
	API_PATH            = "api"
	API_VERSION         = "v4"
	API_METHOD          = http.MethodGet
	API_MAX_NUM_ARTICLE = 100
	API_TIME_FORMAT     = "2006-01-02T15:04:05Z"
	API_RESP_TIME_FMT   = "2006-01-02T15:04:05Z"
)

var API_MIN_TIME, _ = time.Parse(time.DateOnly, "1900-01-01")
var API_URL = fmt.Sprintf("%s://%s/%s/%s", API_SCHEME, API_HOST, API_PATH, API_VERSION)

const (
	EPTopHeadlines = "top-headlines"
	EPSearch       = "search"
)

var EPform2client = map[string]string{
	srv.EPSearch:       EPSearch,
	srv.EPTopHeadlines: EPTopHeadlines,
}
