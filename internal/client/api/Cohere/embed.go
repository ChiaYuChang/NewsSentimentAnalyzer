package cohere

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	pgv "github.com/pgvector/pgvector-go"
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

func (body EmbedResponseBody) ToPgvVectors() <-chan pgv.Vector {
	vecChan := make(chan pgv.Vector)
	for _, embd := range body.Embeddings {
		vecChan <- pgv.NewVector(embd)
	}
	return vecChan
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
	return len(body.Embeddings)
}

func (body EmbedResponseBody) Unwind() []*EmbedObject {
	items := make([]*EmbedObject, 0, body.Len())
	for i := 0; i < body.Len(); i++ {
		items = append(items, &EmbedObject{
			Text:      body.Texts[i],
			Embedding: pgv.NewVector(body.Embeddings[i]),
		})
	}
	return items
}

type EmbedObject struct {
	Text      string     `json:"text"`
	Embedding pgv.Vector `json:"embedding"`
}

func (obj EmbedObject) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Text      string    `json:"text"`
		Embedding []float32 `json:"embedding"`
	}{
		Text:      obj.Text,
		Embedding: obj.Embedding.Slice(),
	}
	return json.Marshal(tmp)
}

func (item EmbedObject) String() string {
	sb := strings.Builder{}
	sb.WriteString("Cohere Embed Item:\n")

	txtLen := utf8.RuneCountInString(item.Text)
	txt := item.Text
	if txtLen > 40 {
		txt = fmt.Sprintf("%s...(%d words)", txt[:40], txtLen-40)
	}

	embdLen := len(item.Embedding.Slice())
	embd := item.Embedding.Slice()
	if len(embd) > 10 {
		embd = embd[:10]
	}
	sb.WriteString(fmt.Sprintf("  Text: %s\n", txt))
	sb.WriteString(fmt.Sprintf("  Embedding: %5.4f (len: %d)\n", embd[:], embdLen))
	return sb.String()
}
