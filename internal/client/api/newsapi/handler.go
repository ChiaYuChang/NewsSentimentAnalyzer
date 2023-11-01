package newsapi

import (
	"net/http"
	"net/url"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
)

func getPage(q url.Values) int {
	page := 1
	if pStr := q.Get(string(Page)); pStr != "" {
		if p, err := convert.StrTo(pStr).Int(); err == nil {
			page = p
		}
	}
	return page
}

type EverythingHandler struct{}

func (hl EverythingHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.NEWSAPIEverything)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	q, err := newRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	var si SearchInField
	si.Parse(data.SearchIn)

	q.WithKeywords(data.Keyword).
		WithSources(data.Sources).
		WithDomains(data.Domains).
		WithExcludeDomains(data.ExcludeDomains).
		WithSearchIn(si).
		WithLanguage(data.Language).
		WithFrom(data.Form).
		WithTo(data.To)
	return q, nil
}

func (hl EverythingHandler) Parse(response *http.Response) (api.Response, error) {
	page := getPage(response.Request.URL.Query())
	return ParseHTTPResponse(response, page)
}

type TopHeadlinesHandler struct{}

func (hl TopHeadlinesHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.NEWSAPITopHeadlines)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	q, err := newRequest(apikey).
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

type SourcesHandler struct{}

func (hl SourcesHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.NEWSAPISources)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	q, err := newRequest(apikey).
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
	page := getPage(response.Request.URL.Query())
	return ParseHTTPResponse(response, page)
}
