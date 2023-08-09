package gnews

import (
	"net/http"
	"net/url"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/GNews"
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

type TopHeadlinesHandler struct{}

func (hl TopHeadlinesHandler) Handle(apikey string, pf pageform.PageForm) (api.Query, error) {
	data, ok := pf.(srv.GNewsHeadlines)
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
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To)
	return q, nil
}

func (hl TopHeadlinesHandler) Parse(response *http.Response) (api.Response, error) {
	page := getPage(response.Request.URL.Query())
	return ParseHTTPResponse(response, page)
}

type SearchHandler struct{}

func (s SearchHandler) Handle(apikey string, pf pageform.PageForm) (api.Query, error) {
	data, ok := pf.(srv.GNewsSearch)
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
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithFrom(data.Form).
		WithTo(data.To)

	return q, nil
}

func (s SearchHandler) Parse(response *http.Response) (api.Response, error) {
	page := getPage(response.Request.URL.Query())
	return ParseHTTPResponse(response, page)
}
