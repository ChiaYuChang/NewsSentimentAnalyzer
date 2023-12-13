package cohere

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func NewEmbedRequest(apikey string, text ...string) Request[EmbedRequestBody] {
	return Request[EmbedRequestBody]{
		apikey: apikey,
		Body: EmbedRequestBody{
			Text: text,
		},
	}
}

type EmbedRequestBody struct {
	Text      []string `json:"text"                                                    validate:"required"`
	Model     string   `json:"model"      mod:"default=embed-multilingual-light-v3.0"  validate:"required,cohere_embed_model"`
	InputType string   `json:"input_type" mod:"default=clustering"                     validate:"cohere_embed_input_type"`
	Truncate  string   `json:"truncate"   mod:"default=END"                            validate:"cohere_truncate"`
}

func (body EmbedRequestBody) Endpoint() string {
	return EPCoEmbed
}

func (body *EmbedRequestBody) WithModel(model string) {
	body.Model = model
}

func (body *EmbedRequestBody) WithInputType(inputType string) {
	body.InputType = inputType
}

func (body *EmbedRequestBody) WithTruncate(truncate string) {
	body.Truncate = truncate
}

type EmbedResponseBody struct {
	Id         string      `json:"id"`
	Texts      []string    `json:"texts"`
	Embeddings [][]float32 `json:"embeddings"`
	Meta       Meta        `json:"meta"`
}

type Meta struct {
	APIVersion struct {
		Version string `json:"version"`
	} `json:"api_version"`
	BilledUnits struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"billed_units"`
}

func (body EmbedResponseBody) Len() int {
	return len(body.Texts)
}

func (body EmbedResponseBody) Unwind() []*EmbedItem {
	items := make([]*EmbedItem, 0, body.Len())
	for i := 0; i < body.Len(); i++ {
		items = append(items, &EmbedItem{
			Text:      body.Texts[i],
			Embedding: body.Embeddings[i],
		})
	}
	return items
}

type EmbedItem struct {
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding"`
}

func (item EmbedItem) String() string {
	sb := strings.Builder{}
	sb.WriteString("Cohere Embed Item:\n")

	txtLen := utf8.RuneCountInString(item.Text)
	txt := item.Text
	if txtLen > 40 {
		txt = fmt.Sprintf("%s...(%d words)", txt[:40], txtLen-40)
	}

	embdLen := len(item.Embedding)
	embd := item.Embedding
	if len(embd) > 10 {
		embd = embd[:10]
	}
	sb.WriteString(fmt.Sprintf("  Text: %s\n", txt))
	sb.WriteString(fmt.Sprintf("  Embedding: %5.4f (len: %d)\n", embd[:], embdLen))
	return sb.String()
}
