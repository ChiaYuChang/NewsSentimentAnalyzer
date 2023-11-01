package cohere

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
)

const (
	API_SCHEME  = "https"
	API_HOST    = "api.cohere.ai"
	API_VERSION = "v1"
	API_METHOD  = http.MethodPost
)

var API_URL = fmt.Sprintf("%s://%s/%s", API_SCHEME, API_HOST, API_VERSION)

const (
	EPWordEmbedding string = "embed"
)

type TruncateType string

const (
	TruncateNone  TruncateType = "NONE"
	TruncateStart TruncateType = "START"
	TruncateEnd   TruncateType = "END"
)

const (
	Texts    api.Key = "texts"
	Truncate api.Key = "truncate"
	Model    api.Key = "model"
)

type WordEmbeddingRequest struct {
	*api.RequestProto
}

func NewWordEmbeddingRequest(apikey string) WordEmbeddingRequest {
	req := WordEmbeddingRequest{api.NewRequestProtoType("")}
	req.SetEndpoint(EPWordEmbedding)
	req.SetApiKey(apikey)
	return req
}

func (req WordEmbeddingRequest) String() string {
	b, _ := json.MarshalIndent(req, "", "    ")
	return string(b)
}

func (req *WordEmbeddingRequest) WithTexts(texts ...string) *WordEmbeddingRequest {
	for _, text := range texts {
		req.Add(Texts, text)
	}
	return req
}

func (req *WordEmbeddingRequest) WithModel(model string) *WordEmbeddingRequest {
	req.Add(Model, model)
	return req
}

func (req *WordEmbeddingRequest) WithTruncate(truncate TruncateType) *WordEmbeddingRequest {
	req.Add(Model, string(truncate))
	return req
}

func (req *WordEmbeddingRequest) ToHttpRequest() (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", API_URL, req.Endpoint())

	b, err := json.Marshal(req.Values)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", req.APIKey()))
	return r, nil
}

// type WordEmbeddingPayload struct {
// 	Texts    []string     `json:"texts"`
// 	Model    string       `json:"model"`
// 	Truncate TruncateType `json:"truncate"`
// }

// func (r WordEmbeddingPayload) String() string {
// 	b, _ := json.MarshalIndent(r, "", "    ")
// 	return string(b)
// }

// func (r WordEmbeddingPayload) Endpoint() string {
// 	return EPWordEmbedding
// }

// func Embedding(payload WordEmbeddingPayload) {
// 	url := "https://api.cohere.ai/v1/embed"

// 	b, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))

// 	req.Header.Add("accept", "application/json")
// 	req.Header.Add("content-type", "application/json")
// 	req.Header.Add("authorization", "Bearer ZQOEHgimFaZSDdWTwaob1ULC8SxHRa3tuLAlzzVn")

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	embd := EmbedResponse{}
// 	err = json.Unmarshal(body, &embd)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(embd)
// }
