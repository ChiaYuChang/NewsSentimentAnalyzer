package openai

import (
	"fmt"
	"net/http"
	"net/url"
)

// See https://platform.openai.com/docs/api-reference/embeddings
type EmbeddingsRequest struct {
	Model string   `json:"model" validate:"required"`
	Input []string `json:"input" validate:"required"`
	User  string   `json:"user,omitempty"`
}

type EmbeddingsResponse struct {
	Object string           `json:"object"`
	Data   []EmbeddingsData `json:"data"`
	Model  string           `json:"model"`
	Usage  Usage            `json:"usage"`
}

type EmbeddingsData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

func NewEmbeddingsRequest() *EmbeddingsRequest {
	return &EmbeddingsRequest{}
}

func (req EmbeddingsRequest) URL() *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s://%s/%s/%s", URL_SCHEME, URL_HOST, API_VERSION, EPEmbeddings))
	return u
}

func (req EmbeddingsRequest) HttpMethod() string {
	return http.MethodPost
}

func (req EmbeddingsRequest) ToHttpRequest(apikey string) (*http.Request, error) {
	return toHttpRequest(req, apikey)
}
