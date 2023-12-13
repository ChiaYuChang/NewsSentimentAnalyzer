package openai

import (
	"crypto/rand"
	"encoding/json"
	"math"
	"math/big"
)

// func SentimentAnalysisRequest() Request[*ChatCompletionsRequestBody] {
// return
// }

type ChatCompletionsRequest Request[ChatCompletionsRequestBody]

// See https://platform.openai.com/docs/api-reference/chat
func NewChatCompletionsRequest(apikey string, messages ...Message) Request[ChatCompletionsRequestBody] {
	return Request[ChatCompletionsRequestBody]{
		Body: ChatCompletionsRequestBody{
			Messages: messages,
		},
		apikey: apikey,
	}
}

// see https://platform.openai.com/docs/api-reference/chat
type ChatCompletionsRequestBody struct {
	Messages         []Message       `json:"messages"                                                validate:"required,dive"`
	Model            string          `json:"model"                       mod:"default=gpt-3.5-turbo" validate:"required"`
	FrequencyPenalty float32         `json:"frequency_penalty,omitempty" mod:"default=0"             validate:"gte=-2.0,lte=2.0"`
	LogitBias        map[string]int  `json:"logit_bias,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	N                int             `json:"n,omitempty"                 mod:"default=1"             validate:"gt=0"`
	PresencePenalty  float32         `json:"presence_penalty,omitempty"  mod:"default=0"             validate:"gte=-2.0,lte=2.0"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Seed             *int64          `json:"seed,omitempty"`
	Stop             []string        `json:"stop,omitempty"                                          validate:"max=4"`
	Stream           bool            `json:"stream,omitempty"            mod:"default=false"`
	Temperature      *float32        `json:"temperature,omitempty"                                   validate:"omitempty,gte=0.0,lte=2.0,excluded_with_all=TopP"`
	TopP             *float32        `json:"top_p,omitempty"                                         validate:"omitempty,gt=0.0,lte=1.0,excluded_with_all=Temperature"`
	Tools            []Tool          `json:"tools,omitempty"`
	ToolChoice       ToolChoice      `json:"tool_choice,omitempty"`
	User             string          `json:"user,omitempty"`
}

func (body ChatCompletionsRequestBody) Endpoint() string {
	return EPChatCompletions
}

func (body *ChatCompletionsRequestBody) AppendSystemMessages(content, name string) *ChatCompletionsRequestBody {
	body.Messages = append(body.Messages, systemMessage{
		Content: content,
		Role:    "system",
		Name:    name,
	})
	return body
}

func (body *ChatCompletionsRequestBody) AppendUserMessages(content, name string) *ChatCompletionsRequestBody {
	body.Messages = append(body.Messages, userMessage{
		Content: content,
		Role:    "user",
		Name:    name,
	})
	return body
}

func (body *ChatCompletionsRequestBody) AppendAssistantMessages(content, name string, toolCalls ...*toolCall) *ChatCompletionsRequestBody {
	body.Messages = append(body.Messages, assistantMessage{
		basicMessage: basicMessage{
			Content: content,
			Role:    "assistant",
			Name:    name,
		},
		ToolCalls: toolCalls,
	})
	return body
}

func (body *ChatCompletionsRequestBody) AppendToolMessages(content, toolCallId string) *ChatCompletionsRequestBody {
	body.Messages = append(body.Messages, toolMessage{
		basicMessage: basicMessage{
			Content: content,
			Role:    "tool",
		},
		ToolCallId: toolCallId,
	})
	return body
}

func (body *ChatCompletionsRequestBody) SetModel(model string) *ChatCompletionsRequestBody {
	body.Model = model
	return body
}

func (body *ChatCompletionsRequestBody) SetTemperature(temperature float32) *ChatCompletionsRequestBody {
	body.Temperature = &temperature
	return body
}

func (body *ChatCompletionsRequestBody) SetTopP(topP float32) *ChatCompletionsRequestBody {
	body.Temperature = &topP
	return body
}

func (body *ChatCompletionsRequestBody) SetSeed(seed int64) *ChatCompletionsRequestBody {
	body.Seed = &seed
	return body
}

func (body *ChatCompletionsRequestBody) RandomSetSeed() *ChatCompletionsRequestBody {
	n, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	body.SetSeed(n.Int64())
	return body
}

// currently only support gpt-4-1106-preview or gpt-3.5-turbo-1106
func (body *ChatCompletionsRequestBody) EnsureJsonResponse() *ChatCompletionsRequestBody {
	body.ResponseFormat = &ResponseFormat{Type: "json_object"}
	return body
}

// Messages
type Message interface {
	json.Marshaler
}

type basicMessage struct {
	Content string `json:"content"                  validate:"required"`
	Role    string `json:"role"                     validate:"required,oneof=system user assistant tool"`
	Name    string `json:"name,omitempty"`
}

func (msg basicMessage) MarshalJSON() ([]byte, error) {
	type InnerBasicMessage basicMessage
	tmp := InnerBasicMessage(msg)
	return json.Marshal(tmp)
}

type systemMessage basicMessage

func (msg systemMessage) MarshalJSON() ([]byte, error) {
	return basicMessage(msg).MarshalJSON()
}

type userMessage basicMessage

func (msg userMessage) MarshalJSON() ([]byte, error) {
	return basicMessage(msg).MarshalJSON()
}

type assistantMessage struct {
	basicMessage
	ToolCalls []*toolCall `json:"tool_calls,omitempty" validate:"dive"`
}

func (msg assistantMessage) MarshalJSON() ([]byte, error) {
	type InnerAssistantMessage assistantMessage
	tmp := InnerAssistantMessage(msg)
	return json.Marshal(tmp)
}

type toolCall struct {
	Id       string           `json:"id"       validate:"required"`
	Type     string           `json:"type"     validate:"required,oneof=function"`
	Function toolCallFunction `json:"function" validate:"required"`
}

type toolCallFunction struct {
	Name      string `json:"name"      validate:"required"`
	Arguments string `json:"arguments" validate:"required"`
}

type toolMessage struct {
	basicMessage
	ToolCallId string `json:"tool_call_id,omitempty" validate:"required"`
}

func (msg toolMessage) MarshalJSON() ([]byte, error) {
	type InnerToolMessage toolMessage
	tmp := InnerToolMessage(msg)
	return json.Marshal(tmp)
}

type ResponseFormat struct {
	Type string `json:"type" validate:"oneof=text json_object"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Description string                 `json:"description,omitempty"`
	Name        string                 `json:"name"`
	Parameters  ToolFunctionParameters `json:"parameters"`
}

type ToolFunctionParameters struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
}

type ToolChoice interface {
	json.Marshaler
}

type ToolChoiceObject struct {
	Type     string `json:"type"`
	Function struct {
		Name string `json:"name"`
	} `json:"function"`
}

func (tcObj ToolChoiceObject) MarshalJSON() ([]byte, error) {
	type InnerToolChoiceObject ToolChoiceObject
	tmp := InnerToolChoiceObject(tcObj)
	return json.Marshal(tmp)
}

type ToolChoiceString string

func (tcStr ToolChoiceString) MarshalJSON() ([]byte, error) {
	return []byte(string(tcStr)), nil
}

const (
	ToolChoiceNone = "none"
	ToolChoiceAuto = "auto"
)

type ChatCompletionsBasicObject struct {
	Id                string `json:"id"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Objcet            string `json:"object"`
}

func (obj ChatCompletionsBasicObject) Endpoint() string {
	return EPChatCompletions
}

type ChatCompletionsObject struct {
	ChatCompletionsBasicObject
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Message      struct {
			Content   string     `json:"content"`
			ToolCalls []toolCall `json:"tool_calls"`
			Role      string     `json:"role"`
		} `json:"message"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

type Usage struct {
	TotalTokens      int `json:"total_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type ChatCompletionsChunkObject struct {
	ChatCompletionsBasicObject
	Choices []struct {
		Delta struct {
			Content   string     `json:"content"`
			ToolCalls []toolCall `json:"tool_calls"`
			Role      string     `json:"role"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}
