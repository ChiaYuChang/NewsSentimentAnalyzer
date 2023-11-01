package api

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
)

var ErrTypeAssertionFailure = errors.New("type assertion failure")
var ErrNotNextPage = errors.New("there are no more pages to query")
var ErrNotImplemented = errors.New("not implemented")
var re = regexp.MustCompile("[\\p{Han}[:alnum:]]")

const QueryOriPageKey = true

func QueryOriginalPageContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, QueryOriPageKey, true)
}

type Key string

type Request interface {
	String() string
	ToHttpRequest() (*http.Request, error)
	json.Marshaler
	json.Unmarshaler
	Encode() string
	Decode(q string) error
}

type Values interface {
	Sep() string
	Add(key Key, val string)
	Set(key Key, val string)
	Del(key Key)
	Get(key Key) string
	Has(key Key) bool
	Clone() (Values, error)
}

type Response interface {
	String() string
	GetStatus() string
	HasNext() bool
	NextPageRequest(body io.Reader) (*http.Request, error)
	Len() int
	ToNews(ctx context.Context, wg *sync.WaitGroup, c chan<- *service.NewsCreateRequest)
}

func MD5Hash(title string, publishedAt time.Time, content ...string) (string, error) {
	hasher := md5.New()

	if _, err := hasher.Write(re.ReplaceAll([]byte(title), []byte{})); err != nil {
		return "", fmt.Errorf("error while writing to hasher: %w", err)
	}

	cs := strings.Join(content, "")
	if _, err := hasher.Write(re.ReplaceAll([]byte(cs), []byte{})); err != nil {
		return "", fmt.Errorf("error while writing to hasher: %w", err)
	}

	if _, err := hasher.Write(re.ReplaceAll([]byte(publishedAt.UTC().Format(time.DateTime)), []byte{})); err != nil {
		return "", fmt.Errorf("error while writing to hasher: %w", err)
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}
