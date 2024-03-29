package newsapi

import (
	"fmt"
	"net/http"
	"time"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"
)

func init() {
	cli.RegisterHandler(
		srv.NEWSAPITopHeadlines{},
		TopHeadlinesHandler{},
		EPform2client)

	cli.RegisterHandler(
		srv.NEWSAPIEverything{},
		EverythingHandler{},
		EPform2client)

	cli.RegisterHandler(
		srv.NEWSAPISources{},
		SourcesHandler{},
		EPform2client)
}

const (
	API_SCHEME                   = "https"
	API_HOST                     = "newsapi.org"
	API_PATH                     = ""
	API_VERSION                  = "v2"
	API_METHOD                   = http.MethodGet
	API_TIME_FORMAT              = "2006-01-02T15:04:05Z"
	API_RESP_TIME_FMT            = "2006-01-02T15:04:05Z"
	API_TIME_DELAY_FOR_FREE_PLAN = 24 * 60 * time.Minute
)

const (
	API_MAX_SOURCES_NUM = 30
	API_MAX_PAGE_SIZE   = 100
)

const (
	API_DEFAULT_PAGE_SIZE = 100
	API_DEFAULT_PAGE      = 1
	API_DEFAULT_ENDPOINT  = EPEverything
)

var API_URL = fmt.Sprintf("%s://%s/%s", API_SCHEME, API_HOST, API_VERSION)
var API_MIN_TIME, _ = time.Parse(time.DateOnly, "1900-01-01")

const (
	EPEverything   string = "everything"
	EPTopHeadlines string = "top-headlines"
	EPSources      string = "top-headlines/sources"
)

var EPform2client = map[string]string{
	srv.EPEverything:   EPEverything,
	srv.EPTopHeadlines: EPTopHeadlines,
	srv.EPSources:      EPSources,
}
