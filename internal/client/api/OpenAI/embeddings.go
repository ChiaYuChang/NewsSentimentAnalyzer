package openai

import (
	"fmt"
	"strings"
)

// See https://platform.openai.com/docs/api-reference/embeddings
func NewEmbeddingsRequest(apikey string, input ...string) Request[EmbeddingsRequestBody] {
	return Request[EmbeddingsRequestBody]{
		Body:   EmbeddingsRequestBody{Input: input},
		apikey: apikey,
	}
}

type EmbeddingsRequestBody struct {
	Model          string   `json:"model"            mod:"default=text-embedding-ada-002"  validate:"required"`
	Input          []string `json:"input"                                                  validate:"required,max=8192"`
	EncodingFormat string   `json:"encoding_format"  mod:"default=float"                   validate:"oneof=float base64"`
	User           string   `json:"user,omitempty"`
}

func (body EmbeddingsRequestBody) Endpoint() string {
	return EPEmbeddings
}

func (body *EmbeddingsRequestBody) WithModel(model string) {
	body.Model = model
}

func (body *EmbeddingsRequestBody) WithUser(user string) {
	body.User = user
}

func (body *EmbeddingsRequestBody) WithEncodingFormat(format string) {
	body.EncodingFormat = format
}

type EmbeddingsResponseBody struct {
	Object string             `json:"object"`
	Data   []EmbeddingsObject `json:"data"`
	Model  string             `json:"model"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type EmbeddingsObject struct {
	Index     int       `json:"index"`
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
}

func (body EmbeddingsObject) String() string {
	sb := strings.Builder{}
	sb.WriteString("Embeddings object:\n")
	sb.WriteString(fmt.Sprintf("\t- Index  : %d\n", body.Index))
	sb.WriteString(fmt.Sprintf("\t- Object : %s\n", body.Object))
	sb.WriteString(fmt.Sprintf("\t- Size   : %d\n", len(body.Embedding)))
	return sb.String()
}

func (body EmbeddingsObject) fString(indent string) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%sEmbeddings object:\n", indent))
	sb.WriteString(fmt.Sprintf("%s\t- Index  : %d\n", indent, body.Index))
	sb.WriteString(fmt.Sprintf("%s\t- Object : %s\n", indent, body.Object))
	sb.WriteString(fmt.Sprintf("%s\t- Size   : %d\n", indent, len(body.Embedding)))
	return sb.String()
}

func (body EmbeddingsResponseBody) Endpoint() string {
	return EPEmbeddings
}

func (body EmbeddingsResponseBody) Len() int {
	return len(body.Data)
}

func (body EmbeddingsResponseBody) String() string {
	sb := strings.Builder{}
	sb.WriteString("Embeddings Response Body:\n")
	sb.WriteString(fmt.Sprintf("\t- Model  : %s\n", body.Model))
	sb.WriteString(fmt.Sprintf("\t- Data   : %d\n", len(body.Data)))
	for _, data := range body.Data {
		sb.WriteString(data.fString("\t\t"))
	}
	sb.WriteString("\t- Usage  :\n")
	sb.WriteString(fmt.Sprintf("\t\t- Prompt Tokens : %d\n", body.Usage.PromptTokens))
	sb.WriteString(fmt.Sprintf("\t\t- Total Tokens  : %d\n", body.Usage.TotalTokens))
	return sb.String()
}
