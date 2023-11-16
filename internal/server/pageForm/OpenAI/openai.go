package openai

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	pf "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
)

const API_NAME = "OpenAI"

const (
	EPCompletions     string = "Completions"
	EPChatCompletions string = "Chat Completions"
	EPEmbeddings      string = "Embeddings"
)

type OpenAICompletions struct {
	Model     string   `form:"model"                 mod:"default=gpt-3.5-turbo" validate:"required"`
	Prompt    []string `form:"prompt"                                            validate:"required"`
	MaxTokens int      `form:"max_tokens,omitempty"`
	N         int      `form:"n,omitempty"           mod:"default=1"             validate:"gt=0"`
}

func (f OpenAICompletions) Endpoint() string {
	return EPCompletions
}

func (f OpenAICompletions) API() string {
	return API_NAME
}

func (f OpenAICompletions) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pf.PageForm, error) {
	return pf.FormDecodeAndValidate[OpenAICompletions](decoder, val, postForm)
}

func (f OpenAICompletions) SelectionOpts() []object.SelectOpts {
	return nil
}

func (f OpenAICompletions) String() string {
	sb := strings.Builder{}
	sb.WriteString("OpenAI Completions\n")
	sb.WriteString(fmt.Sprintf("\t- Model     : %s\n", f.Model))
	sb.WriteString(fmt.Sprintf("\t- Max Token : %d\n", f.MaxTokens))
	sb.WriteString(fmt.Sprintf("\t- N         : %d\n", f.N))
	sb.WriteString("\t- Prompt     :\n")
	for _, p := range f.Prompt {
		sb.WriteString(fmt.Sprintf("\t\t - %s\n", p))
	}
	return sb.String()
}

func (f OpenAICompletions) Key() pageform.PageFormRepoKey {
	return pageform.NewPageFormRepoKey(f.API(), f.Endpoint())
}

type OpenAIChatCompletions struct {
	Message   Message `form:"message"                                          validate:"required"`
	Model     string  `form:"model"                mod:"default=gpt-3.5-turbo" validate:"required"`
	MaxTokens int     `form:"max_token,omitempty"`
	N         int     `form:"n,omitempty"          mod:"default=1"             validate:"gt=0"`
}

type Message struct {
	Content      string     `form:"content"    validate:"required"`
	Functions    []Function `form:"functions,omitempty"`
	FunctionCall struct {
		Arguments string `form:"arguments" validate:"required"`
		Name      string `form:"name"      validate:"required"`
	} `form:"function_call,omitempty"`
	Name string `form:"name,omitempty"`
	Role string `form:"role"               validate:"required"`
}

type Function struct {
	Name        string                      `form:"name"       validate:"required"`
	Description string                      `form:"description,omitempty"`
	Parameters  map[string]json.Unmarshaler `form:"parameters" validate:"required,json"`
}
