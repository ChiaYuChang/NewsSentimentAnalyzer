package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestBody interface {
	Endpoint() string
}

type Request[T RequestBody] struct {
	Body   T
	apikey string
}

func (r Request[T]) String() string {
	data, _ := json.MarshalIndent(r.Body, "", "    ")
	return string(data)
}

func (r Request[T]) EndPoint() string {
	return r.Body.Endpoint()
}

func (r *Request[T]) Modify(ctx context.Context) error {
	return GetModifier().Struct(ctx, &r.Body)
}

func (r Request[T]) Validate(ctx context.Context) error {
	return GetValidator().StructCtx(ctx, r.Body)
}

func (r Request[T]) ToHTTPRequest() (*http.Request, error) {
	body, err := json.Marshal(r.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/%s", API_URL, r.EndPoint()),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.apikey))
	return req, nil
}
