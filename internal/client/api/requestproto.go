package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RequestProto struct {
	ep     string
	apikey string
	*Params
}

func NewRequestProtoType(sep string) *RequestProto {
	return &RequestProto{Params: NewParams(sep)}
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

// apikey getter
func (reqproto RequestProto) APIKey() string {
	return reqproto.apikey
}

// Request interface
func (reqproto RequestProto) String() string {
	return reqproto.Encode()
}

func (reqproto RequestProto) ToHttpRequest(apiURL, apiMethod string,
	body io.Reader, req Request) (*http.Request, error) {
	u, err := reqproto.ToURL(apiURL, apiMethod)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(apiMethod, u.String(), body)
}

func (reqproto RequestProto) ToURL(apiURL, apiMethod string) (*url.URL, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	if apiMethod == http.MethodGet {
		u = u.JoinPath(reqproto.ep)
		u.RawQuery = reqproto.Encode()
	}
	return u, nil
}

func (reqproto RequestProto) ToHTTPRequest(apiURL, apiMethod string, body io.Reader) (*http.Request, error) {
	u, err := reqproto.ToURL(apiURL, apiMethod)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(apiMethod, u.String(), body)
}

func (reqproto RequestProto) AddAPIKeyToQuery(req *http.Request, key Key) *http.Request {
	req.URL.RawQuery = fmt.Sprintf("%s=%s&%s", key, reqproto.apikey, req.URL.RawQuery)
	return req
}
