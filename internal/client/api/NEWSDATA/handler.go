package newsdata

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
)

type LatestNewsHandler struct{}

func (h LatestNewsHandler) Handle(apikey string, pf pageform.PageForm) (api.Query, error) {
	data, ok := pf.(newsdata.NEWSDATAIOLatestNews)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	q, err := newQuery(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithKeywords(data.Keyword).
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	return q, nil
}

func (h LatestNewsHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

type NewsArchiveHandler struct{}

func (h NewsArchiveHandler) Handle(apikey string, pf pageform.PageForm) (api.Query, error) {
	data, ok := pf.(newsdata.NEWSDATAIONewsArchive)
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
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To)

	return q, nil
}

func (h NewsArchiveHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

type NewsSourcesHandler struct{}

func (h NewsSourcesHandler) Handle(apikey string, pf pageform.PageForm) (api.Query, error) {
	data := pf.(newsdata.NEWSDATAIONewsSources)

	q, err := newQuery(apikey).SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	q.WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	return q, nil
}

func (h NewsSourcesHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}
