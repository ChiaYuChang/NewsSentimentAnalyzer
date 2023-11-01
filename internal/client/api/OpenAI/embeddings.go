package openai

import (
	"encoding/json"
)

type EmbeddingsRequest Request[*EmbeddingsRequestBody]

// See https://platform.openai.com/docs/api-reference/embeddings
func NewEmbeddingsRequest(apikey string) EmbeddingsRequest {
	return EmbeddingsRequest{
		Body:   &EmbeddingsRequestBody{},
		apikey: apikey,
	}
}

type EmbeddingsRequestBody struct {
	Model          string   `json:"model"            mod:"default=text-embedding-ada-002"  validate:"required"`
	Input          []string `json:"input"                                                  validate:"required,max=8192"`
	EncodingFormat string   `json:"encoding_format"  mod:"float"                           validate:"oneof=float base64"`
	User           string   `json:"user,omitempty"`
}

func (body EmbeddingsRequestBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(body)
}

func (body *EmbeddingsRequestBody) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, body)
}

func (body EmbeddingsRequestBody) String() string {
	b, _ := json.MarshalIndent(body, "", " ")
	return string(b)
}

func (body EmbeddingsRequestBody) Endpoint() string {
	return EPEmbeddings
}

type EmbeddingsObject struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"  validate:"min=1"`
	Object    string    `json:"object"     validate:"oneof=embedding"` // always embedding
}
