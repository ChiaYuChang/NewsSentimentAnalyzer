package api

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Errors

var ErrTypeAssertionFailure = errors.New("type assertion failure")
var ErrNotNextPage = errors.New("there are no more pages to query")
var ErrNotImplemented = errors.New("not implemented")
var ErrEndOfQuery = errors.New("end of query")
var ErrNextTokenAssertionFailure = errors.New("next token assertion failure")

var re = regexp.MustCompile("[\\p{Han}[:alnum:]]")

const (
	StrLastPageToken = StrNextPageToken("$")
	IntLastPageToken = IntNextPageToken(-1)
)

func IsLastPageToken(token NextPageToken) bool {
	return token.Equal(StrLastPageToken) || token.Equal(IntLastPageToken)
}

const QueryOriPageKey = true

func QueryOriginalPageContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, QueryOriPageKey, true)
}

type Key string

func (k Key) String() string {
	return string(k)
}

type Request interface {
	String() string
	ToHttpRequest() (*http.Request, error)
	json.Marshaler
	json.Unmarshaler
	Encode() string
	Decode(q string) error
	ToPreviewCache(uid uuid.UUID) (cKey string, c *PreviewCache)
}

type Response interface {
	String() string
	GetStatus() string
	HasNext() bool
	ToNewsItemList() (next NextPageToken, preview []NewsPreview)
	Len() int
	ContentProcessFunc(c string) (string, error)
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
