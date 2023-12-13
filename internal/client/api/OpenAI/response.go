package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseBody interface {
	Endpoint() string
}

type Response[T ResponseBody] struct {
	StatusCode int
	Body       T
}

func (r *Response[T]) Endpoint() string {
	return r.Body.Endpoint()
}

func (r *Response[T]) String() string {
	return fmt.Sprintf("code: %d\nbody: %v\n", r.StatusCode, r.Body)
}

func ParseHTTPResponse[T ResponseBody](resp *http.Response) (*Response[T], error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}
	defer resp.Body.Close()

	r := &Response[T]{}
	if r.StatusCode = resp.StatusCode; r.StatusCode != http.StatusOK {
		var errResp ErrorResponseBody
		if err := json.Unmarshal(body, &errResp); err != nil {
			return r, fmt.Errorf("error while parsing error response: %w", err)
		}
		ecErr := errResp.ToEcError(r.StatusCode).
			WithDetails("error while Unmarshal error response")
		return r, ecErr
	}

	if err := json.Unmarshal(body, &r.Body); err != nil {
		return nil, fmt.Errorf("error while Unmarshal response: %w", err)
	}

	return r, nil
}
