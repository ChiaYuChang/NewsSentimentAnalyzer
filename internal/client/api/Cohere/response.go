package cohere

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseBody interface {
	EmbedResponseBody | ChatResponseBody
}

type Response[T ResponseBody] struct {
	StatusCode int
	Body       T
}

func ParseHTTPResponse[T ResponseBody](resp *http.Response) (*Response[T], error) {
	r := &Response[T]{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}
	defer resp.Body.Close()

	if r.StatusCode = resp.StatusCode; r.StatusCode != http.StatusOK {
		var errResp ErrorResponseBody
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("error while parsing error response: %w", err)
		}
		errResp.Code = resp.StatusCode
		err := errResp.ToEcError()
		err.WithDetails("error while parsing response")
		return nil, err
	}

	if err := json.Unmarshal(body, &r.Body); err != nil {
		return nil, fmt.Errorf("error while parsing response: %w", err)
	}

	return r, nil
}
