package openai

// See https://platform.openai.com/docs/api-reference/completions
// func NewCompletionsRequest(apikey string) Request[CompletionsRequestBody] {
// 	return Request[CompletionsRequestBody]{
// 		Body:   CompletionsRequestBody{},
// 		apikey: apikey,
// 	}
// }

// type CompletionsRequestBody struct {
// 	Model            string         `json:"model"                       mod:"default=gpt-3.5-turbo" validate:"required"`
// 	Prompt           []string       `json:"prompt"                                                  validate:"required"`
// 	BestOf           int            `json:"best_of"                     mod:"default=1"             validate:"gte=1 required_with=MaxTokens Stop"`
// 	Echo             bool           `json:"echo,omitempty"`
// 	FrequencyPenalty float32        `json:"frequency_penalty,omitempty" mod:"default=0"             validate:"gte=-2.0,lte=2.0"`
// 	LogitBias        map[string]int `json:"logit_bias,omitempty"`
// 	LogProbs         float32        `json:"logit_probs,omitempty"                                   validate:"gte=0.0,lte=5.0"`
// 	MaxTokens        int            `json:"max_tokens,omitempty"`
// 	N                int            `json:"n,omitempty"                 mod:"default=1"             validate:"gt=0"`
// 	PresencePenalty  float32        `json:"presence_penalty,omitempty"  mod:"default=0"             validate:"gte=-2.0,lte=2.0"`
// 	Stop             []string       `json:"stop,omitempty"                                          validate:"max=4"`
// 	Stream           bool           `json:"stream,omitempty"`
// 	Suffix           string         `json:"suffix,omitempty"`
// 	Temperature      float32        `json:"temperature,omitempty"       mod:"default=1"             validate:"gte=0.0,lte=2.0,excluded_with_all=TopP"`
// 	TopP             float32        `json:"top_p,omitempty"             mod:"default=1"             validate:"gt=0.0,lte=1.0,excluded_with_all=Temperature"`
// 	User             string         `json:"user,omitempty"`
// }

// func (body CompletionsRequestBody) Endpoint() string {
// 	return EPCompletions
// }

// type CompletionsObject struct {
// 	Id      string `json:"id"`
// 	Choices []struct {
// 		FinishReason string                 `json:"finish_reason"`
// 		Index        int                    `json:"index"`
// 		LogProbs     map[string]interface{} `json:"log_probs"`
// 		Text         string                 `json:"text"`
// 	} `json:"choices"`
// 	Created int    `json:"created"`
// 	Model   string `json:"model"`
// 	Object  string `json:"object"`
// 	Usage   Usage  `json:"usage"`
// }
