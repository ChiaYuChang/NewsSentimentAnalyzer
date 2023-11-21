package newsapi

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"
	"github.com/google/uuid"
)

type EverythingHandler struct{}

func (hl EverythingHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {
	data, ok := pf.(srv.NEWSAPIEverything)
	if !ok {
		return "", nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	r, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return "", nil, err
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

	ckey, cache = r.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (hl EverythingHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}

func (hl EverythingHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromCacheQuery(cq)
}

type TopHeadlinesHandler struct{}

func (hl TopHeadlinesHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {

	data, ok := pf.(srv.NEWSAPITopHeadlines)
	if !ok {
		return "", nil, api.ErrTypeAssertionFailure
	}

	q, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return "", nil, err
	}

	q.WithKeywords(data.Keyword).
		WithSources(data.Sources).
		WithCountry(data.Country).
		WithCategory(data.Category)
	ckey, cache = q.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (hl TopHeadlinesHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}

func (hl TopHeadlinesHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromCacheQuery(cq)
}

type SourcesHandler struct{}

func (hl SourcesHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {
	data, ok := pf.(srv.NEWSAPISources)
	if !ok {
		return "", nil, api.ErrTypeAssertionFailure
	}

	q, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return "", nil, err
	}

	q.WithKeywords(data.Language).
		WithCountry(data.Country).
		WithCategory(data.Category)

	ckey, cache = q.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (hl SourcesHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}

func (hl SourcesHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromCacheQuery(cq)
}
