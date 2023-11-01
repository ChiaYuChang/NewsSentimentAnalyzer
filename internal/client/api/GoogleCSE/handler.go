package googlecse

import (
	"context"
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

	q, err := NewSearchRequest(context.Background(),
		apikey, data.SearchEngineID)
	if err != nil {
		return nil, err
	}

	q = q.SetKeyword(data.Keyword).
		SetDateRestict(data.DateRestrict())

	return q, nil
}

func (hl CSEHandler) Parse(resp *http.Response) (api.Response, error) {
	return Response{}, nil
}
