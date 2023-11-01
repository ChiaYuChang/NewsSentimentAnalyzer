package openai_test

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	openai "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/OpenAI"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCompletions(t *testing.T) {
	var exampleCompletionsResp = []byte(`{
  "id": "cmpl-uqkvlQyYK7bGYrRHQ0eXlWi7",
  "object": "text_completion",
  "created": 1589478378,
  "model": "text-davinci-003",
  "choices": [
    {
      "text": "\n\nThis is indeed a test",
      "index": 0,
      "logprobs": null,
      "finish_reason": "length"
    }
  ],
  "usage": {
    "prompt_tokens": 5,
    "completion_tokens": 7,
    "total_tokens": 12
  }
}`)

	r := chi.NewRouter()
	r.Post("/"+openai.EPCompletions, func(w http.ResponseWriter, r *http.Request) {
		btoken := r.Header.Get("Authorization")
		if btoken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("API token is missing"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(exampleCompletionsResp)
	})

	srv := httptest.NewTLSServer(r)
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	u, err := url.Parse(srv.URL + "/" + openai.EPCompletions)
	require.NoError(t, err)

	req := openai.NewCompletionsRequest("[[::API_KEY::]]")
	req.Body.Model = "text-davinci-003"
	req.Body.Prompt = append(req.Body.Prompt, "Say this is a test")
	req.Body.MaxTokens = 10

	httpReq, err := req.ToHTTPRequest()
	httpReq.URL = u
	require.NoError(t, err)

	httpResp, err := cli.Do(httpReq)
	require.NoError(t, err)
	require.NotNil(t, httpResp)
	require.Equal(t, http.StatusOK, httpResp.StatusCode)

	body, err := io.ReadAll(httpResp.Body)
	require.NoError(t, err)
	defer httpResp.Body.Close()

	var resp openai.CompletionsObject
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)
	require.Equal(t, resp.Id, "cmpl-uqkvlQyYK7bGYrRHQ0eXlWi7")
	require.Equal(t, resp.Model, "text-davinci-003")
	require.Equal(t, resp.Usage.CompletionTokens, 7)
	require.Equal(t, len(resp.Choices), 1)
	require.Equal(t, resp.Choices[0].FinishReason, "length")
	srv.Close()
}
