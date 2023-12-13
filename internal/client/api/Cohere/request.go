package cohere

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

var Modifier *mold.Transformer = modifiers.New()

type Body interface {
	EmbedRequestBody | GenerateRequestBody | ChatRequestBody
	Endpoint() string
}

type Request[T Body] struct {
	Body   T
	apikey string
}

func (r Request[T]) Endpoint() string {
	return r.Body.Endpoint()
}

func (r Request[T]) Apikey() string {
	return r.apikey
}

func (r Request[T]) String() string {
	b, _ := json.MarshalIndent(r.Body, "", "  ")
	return string(b)
}

func (r *Request[T]) Modify(ctx context.Context) error {
	return Modifier.Struct(ctx, &r.Body)
}

func (r Request[T]) Validate(ctx context.Context) error {
	val, err := validator.GetDefaultValidate()
	if err != nil {
		return err
	}
	return val.StructCtx(ctx, r.Body)
}

func (r Request[T]) ToHTTPRequest() (*http.Request, error) {
	body, err := json.Marshal(r.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(API_METHOD,
		fmt.Sprintf("%s/%s", API_URL, r.Endpoint()),
		bytes.NewBuffer(body),
	)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.apikey))
	req.Header.Set("Accept", "application/json")
	return req, nil
}
