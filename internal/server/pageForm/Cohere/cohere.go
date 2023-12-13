package cohere

import (
	"fmt"
	"net/url"
	"strings"

	pf "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
)

const API_NAME = "Cohere"

const (
	EPChat  string = "Chat"
	EPEmbed string = "Embed"
)

type EmbeddingsOptions struct {
	Model     string `form:"embedding-model" validate:"required,oneof=embed-english-v3.0 embed-multilingual-v3.0 embed-english-light-v3.0 embed-multilingual-light-v3.0"`
	InputType string `form:"input-type"      validate:"required,oneof=search_document search_query classification clustering"`
}

func (f EmbeddingsOptions) Endpoint() string {
	return EPEmbed
}

func (f EmbeddingsOptions) API() string {
	return API_NAME
}

func (f EmbeddingsOptions) Key() pf.PageFormRepoKey {
	return pf.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f EmbeddingsOptions) SelectionOpts() []object.SelectOpts {
	return nil
}

func (f EmbeddingsOptions) String() string {
	sb := strings.Builder{}
	sb.WriteString("Cohere Embeddings Options\n")
	sb.WriteString(fmt.Sprintf("\t- Model     : %s\n", f.Model))
	sb.WriteString(fmt.Sprintf("\t- Input Type: %s\n", f.InputType))
	return sb.String()
}

func (f EmbeddingsOptions) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pf.PageForm, error) {
	return pf.FormDecodeAndValidate[EmbeddingsOptions](decoder, val, postForm)
}

type SentimentAnalysisOptions struct {
	Prompt   string `form:"prompt"`
	MaxToken int    `form:"max-token"  validate:"gt=0,lte=2048"`
	Truncate string `form:"truncate"   validate:"required,oneof=END START"`
}

func (f SentimentAnalysisOptions) Endpoint() string {
	return EPChat
}

func (f SentimentAnalysisOptions) API() string {
	return API_NAME
}

func (f SentimentAnalysisOptions) Key() pf.PageFormRepoKey {
	return pf.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f SentimentAnalysisOptions) SelectionOpts() []object.SelectOpts {
	return nil
}

func (f SentimentAnalysisOptions) String() string {
	sb := strings.Builder{}
	sb.WriteString("Cohere Sentiment Anaysis Options\n")
	sb.WriteString(fmt.Sprintf("\t- Prompt   : %s\n", f.Prompt))
	sb.WriteString(fmt.Sprintf("\t- Max Token: %d\n", f.MaxToken))
	sb.WriteString(fmt.Sprintf("\t- Truncate : %s\n", f.Truncate))
	return sb.String()
}

func (f SentimentAnalysisOptions) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pf.PageForm, error) {
	return pf.FormDecodeAndValidate[EmbeddingsOptions](decoder, val, postForm)
}
