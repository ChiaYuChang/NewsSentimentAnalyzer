package newsdata

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/NEWSDATA"
)

type LatestNewsHandler struct{}

func (h LatestNewsHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(newsdata.NEWSDATAIOLatestNews)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	req, err := NewRequest(apikey).
		SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	req.WithKeywords(data.Keyword).
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	return req, nil
}

func (h LatestNewsHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

type NewsArchiveHandler struct{}

func (h NewsArchiveHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(newsdata.NEWSDATAIONewsArchive)
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
		WithDomain(data.Domains).
		WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...).
		WithFrom(data.Form).
		WithTo(data.To)

	return req, nil
}

func (h NewsArchiveHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}

type NewsSourcesHandler struct{}

func (h NewsSourcesHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data := pf.(newsdata.NEWSDATAIONewsSources)

	req, err := NewRequest(apikey).SetEndpoint(data.Endpoint())
	if err != nil {
		return nil, err
	}

	req.WithLanguage(data.Language...).
		WithCountry(data.Country...).
		WithCategory(data.Category...)

	return req, nil
}

func (h NewsSourcesHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}
