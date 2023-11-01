package openai

import (
	"encoding/json"
)

type ChatCompletionsRequest Request[*ChatCompletionsRequestBody]

// See https://platform.openai.com/docs/api-reference/chat
func NewChatCompletionsRequest(apikey string) Request[*ChatCompletionsRequestBody] {
	return Request[*ChatCompletionsRequestBody]{
		Body:   &ChatCompletionsRequestBody{},
		apikey: apikey,
	}
}

type ChatCompletionsRequestBody struct {
	Message          Message        `json:"message"                                                 validate:"required"`
	Model            string         `json:"model"                       mod:"default=gpt-3.5-turbo" validate:"required"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty" mod:"default=0"             validate:"gte=-2.0,lte=2.0"`
	FunctionCall     string         `json:"function_call,omitempty"     mod:"default=none"`
	Functions        []Function     `json:"functions,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	MaxTokens        int            `json:"max_token,omitempty"`
	N                int            `json:"n,omitempty"                 mod:"default=1"             validate:"gt=0"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"  mod:"default=0"             validate:"gte=-2.0,lte=2.0"`
	Stop             []string       `json:"stop,omitempty"                                          validate:"max=4"`
	Stream           bool           `json:"stream,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"       mod:"default=1"             validate:"gte=0.0,lte=2.0,excluded_with_all=TopP"`
	TopP             float32        `json:"top_p,omitempty"             mod:"default=1"             validate:"gt=0.0,lte=1.0,excluded_with_all=Temperature"`
	User             string         `json:"user,omitempty"`
}

type Message struct {
	Content      string `json:"content"    validate:"required"`
	FunctionCall struct {
		Arguments string `json:"arguments" validate:"required"`
		Name      string `json:"name"      validate:"required"`
	} `json:"function_call,omitempty"`
	Name string `json:"name,omitempty"`
	Role string `json:"role"               validate:"required"`
}

type Function struct {
	Name        string                      `json:"name"       validate:"required"`
	Description string                      `json:"description,omitempty"`
	Parameters  map[string]json.Unmarshaler `json:"parameters" validate:"required,json"`
}

func (body ChatCompletionsRequestBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(body)
}

func (body *ChatCompletionsRequestBody) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, body)
}

func (body ChatCompletionsRequestBody) Endpoint() string {
	return EPChatCompletions
}

func (body ChatCompletionsRequestBody) String() string {
	b, _ := json.MarshalIndent(body, "", " ")
	return string(b)
}

type ChatCompletionsChunkObject struct {
	Id      string `json:"id"`
	Choices []struct {
		Delta        Delta  `json:"delta"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
}

type Delta struct {
	Content      string `json:"content,omitempty"`
	FunctionCall struct {
		Arguments string `json:"arguments"`
		Name      string `json:"name"`
	} `json:"function_call,omitempty"`
	Role string `json:"role"`
}
