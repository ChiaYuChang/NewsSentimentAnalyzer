package newsdata

import (
	"fmt"
	"net/http"
	"time"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/NEWSDATA"
)

func init() {
	cli.RegisterHandler(
		srv.NEWSDATAIOLatestNews{},
		LatestNewsHandler{},
		EPform2client)

	cli.RegisterHandler(
		srv.NEWSDATAIONewsArchive{},
		NewsArchiveHandler{},
		EPform2client)

	cli.RegisterHandler(
		srv.NEWSDATAIONewsSources{},
		NewsSourcesHandler{},
		EPform2client)
}

const (
	API_SCHEME           = "https"
	API_HOST             = "newsdata.io"
	API_PATH             = "api"
	API_VERSION          = "1"
	API_METHOD           = http.MethodGet
	API_TIME_FORMAT      = "2006-01-02"
	API_MAX_NUM_DOMAIN   = 5
	API_MAX_NUM_COUNTRY  = 5
	API_MAX_NUM_CATEGORY = 5
	API_MAX_NUM_LANGUAGE = 5
	API_RESP_TIME_FMT    = "2006-01-02 15:04:05"
)

var API_MIN_TIME, _ = time.Parse(time.DateOnly, "1900-01-01")
var API_URL = fmt.Sprintf("%s://%s/%s/%s", API_SCHEME, API_HOST, API_PATH, API_VERSION)

const (
	EPLatestNews  string = "news"
	EPCrypto      string = "crypto"
	EPNewsArchive string = "archive"
	EPNewsSources string = "sources"
)

var EPform2client = map[string]string{
	srv.EPLatestNews:  EPLatestNews,
	srv.EPCrypto:      EPCrypto,
	srv.EPNewsArchive: EPNewsArchive,
	srv.EPNewsSources: EPNewsSources,
}
