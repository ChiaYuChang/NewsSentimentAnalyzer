package newsdata

import (
	"fmt"
	"net/http"
	"net/url"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
)

func init() {
	cli.PageFormHandlerRepo.RegisterPageForm(
		srv.NEWSDATAIOLatestNews{}, HandleLatestNewsQuery)
	cli.PageFormHandlerRepo.RegisterPageForm(
		srv.NEWSDATAIONewsArchive{}, HandleNewsArchive)
	cli.PageFormHandlerRepo.RegisterPageForm(
		srv.NEWSDATAIONewsSources{}, HandleNewsSources)
}

const (
	API_ROOT             = "https://newsdata.io/api"
	API_VERSION          = "1"
	API_METHOD           = http.MethodGet
	API_TIME_FORMAT      = "2006-01-02"
	API_MAX_NUM_DOMAIN   = 5
	API_MAX_NUM_COUNTRY  = 5
	API_MAX_NUM_CATEGORY = 5
	API_MAX_NUM_LANGUAGE = 5
)

var API_URL, _ = url.Parse(fmt.Sprintf("%s/%s", API_ROOT, API_VERSION))

var c = srv.API_NAME

const (
	EPLatestNews  string = "news"
	EPCrypto      string = "crypto"
	EPNewsArchive string = "archive"
	EPNewsSources string = "sources"
)
