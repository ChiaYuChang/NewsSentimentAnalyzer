package openai

import (
	"fmt"
	"net/http"
	"net/url"
)

// See https://platform.openai.com/docs/api-reference/completions
type CompletionsRequest struct {
	Model            string         `json:"model" validate:"required"`
	Prompt           []string       `json:"prompt" validate:"required"`
	Suffix           string         `json:"suffix,omitempty"`
	MaxTokens        int            `json:"max_token,omitempty"`
	Temperature      float32        `json:"temperature,omitempty" validate:"gte=0.0,lte=2.0"`
	TopP             float32        `json:"top_p,omitempty" validate:"gt=0.0,lte=1.0"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	LogProbs         int            `json:"logprobs,omitempty" validate:"gt=0.0,lte=5.0"`
	Echo             bool           `json:"echo,omitempty"`
	Stop             []string       `json:"stop,omitempty" validate:"max=4"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty" validate:"gte=-2.0,lte=2.0"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty" validate:"gte=-2.0,lte=2.0"`
	BestOf           int            `json:"best_of,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

type CompletionsResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []Choise `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choise struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	LogProbs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewCompletionsRequest() *CompletionsRequest {
	return &CompletionsRequest{}
}

func (req CompletionsRequest) URL() *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s://%s/%s/%s", URL_SCHEME, URL_HOST, API_VERSION, EPCompletions))
	return u
}

func (req CompletionsRequest) HttpMethod() string {
	return http.MethodPost
}

func (req CompletionsRequest) ToHttpRequest(apikey string) (*http.Request, error) {
	return toHttpRequest(req, apikey)
}

var defaultCompletionsRequest = CompletionsRequest{
	MaxTokens:        16,
	Temperature:      1,
	TopP:             1,
	N:                1,
	PresencePenalty:  0,
	FrequencyPenalty: 0,
	BestOf:           1,
}
