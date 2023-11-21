package gnews

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"
	"github.com/google/uuid"
)

type TopHeadlinesHandler struct{}

func (h TopHeadlinesHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {
	data, ok := pf.(srv.GNewsHeadlines)
	if !ok {
		return "", nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	req, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return "", nil, err
	}

	req.WithKeywords(data.Keyword).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To).
		WithPage(1)

	ckey, cache = req.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (h TopHeadlinesHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}

func (h TopHeadlinesHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromPreviewCache(cq)
}

type SearchHandler struct{}

func (h SearchHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {

	data, ok := pf.(srv.GNewsSearch)
	if !ok {
		return "", nil, api.ErrTypeAssertionFailure
	}
	data.TimeRange.ToUTC()

	req, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return "", nil, err
	}

	req.WithKeywords(data.Keyword).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithFrom(data.Form).
		WithTo(data.To).
		WithPage(1)

	ckey, cache = req.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (h SearchHandler) Parse(response *http.Response) (api.Response, error) {
	return ParseHTTPResponse(response)
}

func (h SearchHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromPreviewCache(cq)
}
