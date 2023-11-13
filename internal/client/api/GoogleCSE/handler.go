package googlecse

import (
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GoogleCSE"
)

type CSEHandler struct{}

func (hl CSEHandler) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	data, ok := pf.(srv.GoogleCSE)
	if !ok {
		return nil, api.ErrTypeAssertionFailure
	}

	req, err := NewRequest(apikey, data.SearchEngineID)
	if err != nil {
		return nil, err
	}

	req = req.SetKeyword(data.Keyword).
		SetDateRestict(data.DateRestrict())
	return req, nil
}

func (hl CSEHandler) Parse(resp *http.Response) (api.Response, error) {
	return ParseHTTPResponse(resp)
}
