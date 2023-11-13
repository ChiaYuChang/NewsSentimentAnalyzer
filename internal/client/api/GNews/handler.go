package gnews

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"
)

type TopHeadlinesHandler struct{}

func (h TopHeadlinesHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.GNewsHeadlines)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	req, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	req.WithKeywords(data.Keyword).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To).
		WithPage(1)

	return req, nil
}

func (h TopHeadlinesHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}

type SearchHandler struct{}

func (h SearchHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.GNewsSearch)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	req, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	req.WithKeywords(data.Keyword).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithFrom(data.Form).
		WithTo(data.To).
		WithPage(1)

	return req, nil
}

func (h SearchHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}
