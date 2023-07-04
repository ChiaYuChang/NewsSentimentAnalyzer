package gnews

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	API_ROOT            = "https://gnews.io/api"
	API_VERSION         = "v4"
	API_METHOD          = http.MethodGet
	API_TIME_FORMAT     = "2006-01-02T15:04:05Z"
	API_MAX_NUM_ARTICLE = 100
)

var API_URL, _ = url.Parse(fmt.Sprintf("%s/%s", API_ROOT, API_VERSION))

const (
	EPTopHeadlines = "top-headlines"
	EPSearch       = "search"
)
