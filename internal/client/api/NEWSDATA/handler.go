package newsdata

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/NEWSDATA"
	"github.com/google/uuid"
)

type LatestNewsHandler struct{}

func (hl LatestNewsHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {
	data, ok := pf.(newsdata.NEWSDATAIOLatestNews)
	if !ok {
		return "", nil, api.ErrTypeAssertionFailure
	}

	req, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return "", nil, err
	}

	req.WithKeywords(data.Keyword).
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	ckey, cache = req.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (hl LatestNewsHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

func (hl LatestNewsHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromCacheQuery(cq)
}

type NewsArchiveHandler struct{}

func (hl NewsArchiveHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {

	data, ok := pf.(newsdata.NEWSDATAIONewsArchive)
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
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To)

	ckey, cache = req.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (hl NewsArchiveHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

func (hl NewsArchiveHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromCacheQuery(cq)
}

type NewsSourcesHandler struct{}

func (hl NewsSourcesHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {
	data := pf.(newsdata.NEWSDATAIONewsSources)

	req, err := NewRequest(apikey).SetEndpoint(data.Endpoint())
	if err != nil {
		return "", nil, err
	}

	req.WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	ckey, cache = req.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (h1 NewsSourcesHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

func (hl NewsSourcesHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromCacheQuery(cq)
}
