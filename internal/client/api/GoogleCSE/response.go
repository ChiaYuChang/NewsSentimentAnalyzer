package googlecse

import (
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
)

type Response struct{}

func (r Response) String() string {
	return ""
}

func (r Response) GetStatus() string {
	return "OK"
}

func (r Response) HasNext() bool {
	return false
}

func (r Response) NextPageRequest(body io.Reader) (*http.Request, error) {
	return nil, nil
}

func (r Response) Len() int {
	return 0
}

func (r Response) ToNews(ctx context.Context, wg *sync.WaitGroup, c chan<- *service.NewsCreateRequest) {
	return
}
