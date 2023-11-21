package googlecse

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GoogleCSE"
	"github.com/google/uuid"
)

type CSEHandler struct{}

func (hl CSEHandler) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {

	data, ok := pf.(srv.GoogleCSE)
	if !ok {
		return "", nil, api.ErrTypeAssertionFailure
	}

	req, err := NewRequest(apikey, data.SearchEngineID)
	if err != nil {
		return "", nil, err
	}
	req.SetEndpoint(pf.Endpoint())

	req = req.SetKeyword(data.Keyword).
		SetDateRestict(data.DateRestrict())

	ckey, cache = req.ToPreviewCache(uid)
	return ckey, cache, nil
}

func (hl CSEHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

func (hl CSEHandler) RequestFromCacheQuery(cq api.CacheQuery) (api.Request, error) {
	return RequestFromCacheQuery(cq)
}
