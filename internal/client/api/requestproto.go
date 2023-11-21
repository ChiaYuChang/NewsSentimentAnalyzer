package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/google/uuid"
)

type RequestProto struct {
	api     string
	ep      string
	apikey  string
	cPrefix string
	cSuffix string
	*Params
}

func NewRequestProtoType(api, sep string) *RequestProto {
	return &RequestProto{api: api, Params: NewParams(sep)}
}

// api endpoint setter
func (reqproto *RequestProto) SetEndpoint(ep string) *RequestProto {
	reqproto.ep = ep
	return reqproto
}

// api endpoint getter
func (reqproto RequestProto) Endpoint() string {
	return reqproto.ep
}

// apikey setter
func (reqproto *RequestProto) SetApiKey(apikey string) *RequestProto {
	reqproto.apikey = apikey
	return reqproto
}

func (reqproto *RequestProto) SetCacheKeyPrefix(prefix string) *RequestProto {
	reqproto.cPrefix = prefix
	return reqproto
}

func (reqproto *RequestProto) SetCacheKeySuffix(suffix string) *RequestProto {
	reqproto.cSuffix = suffix
	return reqproto
}

// apikey getter
func (reqproto RequestProto) APIKey() string {
	return reqproto.apikey
}

// Request interface
func (reqproto RequestProto) String() string {
	return reqproto.Encode()
}

// return a request object with given url and http method, rawquery will not be appended
func (reqproto RequestProto) ToHTTPRequest(apiURL, apiMethod string, body io.Reader) (*http.Request, error) {
	ep := reqproto.Endpoint()
	if ep != "" {
		apiURL += "/" + ep
	}
	return http.NewRequest(apiMethod, apiURL, body)
}

func (reqproto RequestProto) AddAPIKeyToQuery(req *http.Request, key Key) *http.Request {
	req.URL.RawQuery = fmt.Sprintf("%s=%s&%s", key, reqproto.apikey, req.URL.RawQuery)
	return req
}

func (reqproto RequestProto) ToPreviewCache(uid uuid.UUID, next NextPageToken, other map[string]string) (cKey string, c *PreviewCache) {
	c = &PreviewCache{
		Query: CacheQuery{
			UserId: uid,
			API: API{
				Key:      reqproto.apikey,
				Name:     reqproto.api,
				Endpoint: reqproto.ep,
			},
			RawQuery: reqproto.Encode(),
			Body:     "",
			NextPage: next,
			Other:    other,
		},
		NewsItem:  []NewsPreview{},
		CreatedAt: time.Now().UTC(),
	}

	_ = c.AddRandomSalt(global.PREVIEW_CACHE_KEY_SALT_LEN)

	return c.Key(
		global.PREVIEW_CACHE_KEY_PREFIX,
		global.PREVIEW_CACHE_KEY_SUFFIX,
	), c
}
