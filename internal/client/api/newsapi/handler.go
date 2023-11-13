package newsapi

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"
)

type EverythingHandler struct{}

func (hl EverythingHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.NEWSAPIEverything)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	r, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	var si SearchInField
	si.Parse(data.SearchIn)

	r.WithKeywords(data.Keyword).
		WithSources(data.Sources).
		WithDomains(data.Domains).
		WithExcludeDomains(data.ExcludeDomains).
		WithSearchIn(si).
		WithLanguage(data.Language).
		WithFrom(data.Form).
		WithTo(data.To)
	return r, nil
}

func (hl EverythingHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}

type TopHeadlinesHandler struct{}

func (hl TopHeadlinesHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.NEWSAPITopHeadlines)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	q, err := NewRequest(apikey).
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
	return ParseHTTPResponse(response)
}

type SourcesHandler struct{}

func (hl SourcesHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.NEWSAPISources)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	q, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Language).
		WithCountry(data.Country).
		WithCategory(data.Category)
	return q, nil
}

func (hl SourcesHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}
