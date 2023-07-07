package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	URL_SCHEME  = "https"
	URL_HOST    = "api.openai.com"
	API_VERSION = "v1"
)

// API Endpoints
const (
	EPCompletions string = "completions"
	EPEmbeddings  string = "embeddings"
)

type Request interface {
	URL() *url.URL
	HttpMethod() string
	ToHttpRequest(apikey string) (*http.Request, error)
}

func toHttpRequest(req Request, apikey string) (*http.Request, error) {
	jsonObj, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		req.HttpMethod(),
		req.URL().String(),
		bytes.NewReader(jsonObj),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))
	return httpReq, nil
}
