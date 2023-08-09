package newsapi

import (
	"net/http"
	"net/url"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/newsapi"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
)

func getPage(q url.Values) int {
	page := 1
	if pStr := q.Get(qPage); pStr != "" {
		if p, err := convert.StrTo(pStr).Int(); err == nil {
			page = p
		}
	}
	return page
}

type EverythingHandler struct{}

func (hl EverythingHandler) Handle(apikey string, pf pageform.PageForm) (api.Query, error) {
	data, ok := pf.(srv.NEWSAPIEverything)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	q, err := newQuery(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Keyword).
		WithSources(data.Sources).
		WithDomains(data.Domains).
		WithExcludeDomains(data.ExcludeDomains).
		WithSearchIn(data.SearchIn.String()).
		WithLanguage(data.Language...).
		WithFrom(data.Form).
		WithTo(data.To)
	return q, nil
}

func (hl EverythingHandler) Parse(response *http.Response) (api.Response, error) {
	page := getPage(response.Request.URL.Query())
	return ParseHTTPResponse(response, page)
}

type TopHeadlinesHandler struct{}

func (hl TopHeadlinesHandler) Handle(apikey string, pf pageform.PageForm) (api.Query, error) {
	data, ok := pf.(srv.NEWSAPITopHeadlines)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	q, err := newQuery(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Keyword).
		WithSources(data.Sources).
		WithCountry(data.Country).
		WithCategory(data.Category)
	return q, nil
}

func (hl TopHeadlinesHandler) Parse(response *http.Response) (api.Response, error) {
	page := getPage(response.Request.URL.Query())
	return ParseHTTPResponse(response, page)
}
