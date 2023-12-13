package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	openai "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/OpenAI"
	"github.com/stretchr/testify/require"
)

const (
	TEST_USER          = "[[::TEST_USER::]]"
	TEST_API_KEY       = "[[::TEST_API_KEY::]]"
	TEST_ERROR_API_KEY = "[[::TEST_ERROR_API_KEY::]]"
)

func TestCompletions(t *testing.T) {
	// f, err := os.ReadFile("../../../../secrets/.APIKEY.yaml")
	// require.NoError(t, err)
	// require.NotNil(t, f)

	// m := map[string]string{}
	// err = yaml.Unmarshal(f, &m)
	// require.NoError(t, err)

	// apikey := m["OpenAI"]

	req := openai.NewEmbeddingsRequest(
		TEST_API_KEY,
		"義大利首都羅馬有了第一家專門服務狗狗的美食餐廳。菜單上的餐點按照狗狗的體重有4種不同的份量，避掉了可能誘發過敏的食材，並有專門為狗烹煮的魚料理，吸引不少愛狗人士上門消費。",
		"立法院今（1）日三讀通過《道路交通安全基本法》，以「2050年道路交通事故零死亡為目標」，未來由行政院統籌交通安全會報。而目前，以行人環境優化、加強執法等面向，來降低死傷。學者分析，台灣交通要轉換為'人本'主義，除了交通設施逐步改善，回歸教育，用路人行為改變，也是關鍵。",
		"以哈衝突延長休戰無望，以色列總理納坦雅胡1日宣稱武裝組織哈瑪斯違反協議，在休戰最後1小時用火箭攻擊以色列遭以軍攔截，雙方隨即恢復在加 薩的戰鬥狀態。而休戰最後1天內，哈瑪斯釋放8名以色列人質，以色列也放了30名巴勒斯坦囚犯。",
	)
	require.NotNil(t, req)
	require.NotNil(t, openai.GetModifier())
	req.Body.WithUser(TEST_USER)

	err := req.Modify(context.Background())
	require.NoError(t, err)
	require.Equal(t, TEST_USER, req.Body.User)

	httpReq, err := req.ToHTTPRequest()
	require.NoError(t, err)
	require.NotNil(t, httpReq)

	f, err := os.ReadFile("./example_response/embeddings.json")
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc(
		fmt.Sprintf("/%s/%s", openai.API_VERSION, req.Body.Endpoint()),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if !strings.Contains(r.Header.Get("Authorization"), TEST_API_KEY) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{
    "error": {
        "message": "Incorrect API key provided: sk-o39RB*********************************EFxY. You can find your API key at https://platform.openai.com/account/api-keys.",
        "type": "invalid_request_error",
        "param": null,
        "code": "invalid_api_key"
    }
}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(f))
		},
	)

	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	srvrUrl, err := url.Parse(srvr.URL)
	require.NoError(t, err)
	httpReq.URL.Scheme = srvrUrl.Scheme
	httpReq.URL.Host = srvrUrl.Host

	cli := http.Client{Timeout: 3 * time.Second}

	httpResp, err := cli.Do(httpReq)
	require.NoError(t, err)
	require.NotNil(t, httpResp)
	require.Equal(t, http.StatusOK, httpResp.StatusCode)

	resp, err := openai.ParseHTTPResponse[openai.EmbeddingsResponseBody](httpResp)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TEST_ERROR_API_KEY))
	httpResp, err = cli.Do(httpReq)
	require.NoError(t, err)
	require.NotNil(t, httpResp)
	require.Equal(t, http.StatusUnauthorized, httpResp.StatusCode)

	resp, err = openai.ParseHTTPResponse[openai.EmbeddingsResponseBody](httpResp)
	require.Error(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestChatCompletionsObjectUnmarshal(t *testing.T) {
	files := []string{
		"chat_completion_gpt-3.5-turbo.json",
		"chat_completion_gpt-4.json",
	}

	for _, fn := range files {
		b, err := os.ReadFile("example_response/" + fn)
		require.NoError(t, err)

		var respObj openai.ChatCompletionsObject
		err = json.Unmarshal(b, &respObj)
		require.NoError(t, err)
	}
}

const prompt = "As an AI specializing in language and emotion analysis, your task is to assess the sentiments conveyed in a set of statements. Each statement will be enclosed with the symbols [^] and [$]. Consider the overall tone, emotional nuances, and context within the statements. Classify each statement into one of five categories: 1 for very negative, 2 for negative, 3 for neutral, 4 for positive, and 5 for very positive. Please present your responses in a JSON list. Respond sequentially without repeating the provided sentences. For instance, if the sentence is '[^]I love this movie[$] [^]I hate you[$]'', your corresponding response should be [4, 1]."

func TestChatCompletions(t *testing.T) {
	files := []string{
		"bad_request.json",
		"chat_completion_gpt-3.5-turbo.json",
		"chat_completion_gpt-4.json",
		"unauthorizeation.json",
	}

	openAIResponses := make(map[string][]byte, len(files))
	for _, fn := range files {
		b, err := os.ReadFile("example_response/" + fn)
		require.NoError(t, err)
		k := strings.TrimSuffix(fn, ".json")
		openAIResponses[k] = b
	}

	type testCase struct {
		Name           string
		HttpStatusCode int
		ResponseKey    string
		NewRequestFunc func(t *testing.T) openai.Request[openai.ChatCompletionsRequestBody]
	}

	tcs := []testCase{
		{
			Name:           "ok - chat_completion_gpt-4",
			HttpStatusCode: http.StatusOK,
			ResponseKey:    "chat",
			NewRequestFunc: func(t *testing.T) openai.Request[openai.ChatCompletionsRequestBody] {
				req := openai.NewChatCompletionsRequest(TEST_API_KEY)
				req.Body.SetTemperature(0).
					SetModel("gpt-4").
					AppendSystemMessages(prompt, "").
					AppendUserMessages("[^]義大利首都羅馬有了第一家專門服務狗狗的美食餐廳。菜單上的餐點按照狗狗的體重有4種不同的份量，避掉了可能誘發過敏的食材，並有專門為狗烹煮的魚料理，吸引不少愛狗人士上門消費。[$][^]立法院今（1）日三讀通過《道路交通安全基本法》，以「2050年道路交通事故零死亡為目標」，未來由行政院統籌交通安全會報。而目前，以行人環境優化、加強執法等面向，來降低死傷。學者分析，台灣交通要轉換為'人本'主義，除了交通設施逐步改善，回歸教育，用路人行為改變，也是關鍵。[$]", "").
					RandomSetSeed()
				req.Body.MaxTokens = 100

				err := req.Validate(context.Background())
				require.Error(t, err)

				err = req.Modify(context.Background())
				require.NoError(t, err)
				require.Equal(t, "gpt-4", req.Body.Model)
				require.Equal(t, 1, req.Body.N)

				err = req.Validate(context.Background())
				require.NoError(t, err)
				return req
			},
		},
		{
			Name:           "ok - chat_completion_gpt-3.5-turbo",
			HttpStatusCode: http.StatusOK,
			ResponseKey:    "chat_completion_gpt-3.5-turbo",
			NewRequestFunc: func(t *testing.T) openai.Request[openai.ChatCompletionsRequestBody] {
				req := openai.NewChatCompletionsRequest(TEST_API_KEY)
				req.Body.SetTemperature(0).
					AppendSystemMessages(prompt, "").
					AppendUserMessages("[^]義大利首都羅馬有了第一家專門服務狗狗的美食餐廳。菜單上的餐點按照狗狗的體重有4種不同的份量，避掉了可能誘發過敏的食材，並有專門為狗烹煮的魚料理，吸引不少愛狗人士上門消費。[$][^]立法院今（1）日三讀通過《道路交通安全基本法》，以「2050年道路交通事故零死亡為目標」，未來由行政院統籌交通安全會報。而目前，以行人環境優化、加強執法等面向，來降低死傷。學者分析，台灣交通要轉換為'人本'主義，除了交通設施逐步改善，回歸教育，用路人行為改變，也是關鍵。[$]", "").
					RandomSetSeed()
				req.Body.MaxTokens = 100

				err := req.Validate(context.Background())
				require.Error(t, err)

				err = req.Modify(context.Background())
				require.NoError(t, err)
				require.Equal(t, "gpt-3.5-turbo", req.Body.Model)
				require.Equal(t, 1, req.Body.N)

				err = req.Validate(context.Background())
				require.NoError(t, err)
				return req
			},
		},
		{
			Name:           "bad request",
			HttpStatusCode: http.StatusBadRequest,
			ResponseKey:    "bad_request",
			NewRequestFunc: func(t *testing.T) openai.Request[openai.ChatCompletionsRequestBody] {
				req := openai.NewChatCompletionsRequest(TEST_API_KEY)
				req.Body.SetTemperature(0).
					AppendSystemMessages(prompt, "").
					AppendUserMessages("[^]義大利首都羅馬有了第一家專門服務狗狗的美食餐廳。菜單上的餐點按照狗狗的體重有4種不同的份量，避掉了可能誘發過敏的食材，並有專門為狗烹煮的魚料理，吸引不少愛狗人士上門消費。[$][^]立法院今（1）日三讀通過《道路交通安全基本法》，以「2050年道路交通事故零死亡為目標」，未來由行政院統籌交通安全會報。而目前，以行人環境優化、加強執法等面向，來降低死傷。學者分析，台灣交通要轉換為'人本'主義，除了交通設施逐步改善，回歸教育，用路人行為改變，也是關鍵。[$]", "").
					EnsureJsonResponse()
				_ = req.Validate(context.Background())
				_ = req.Modify(context.Background())
				_ = req.Validate(context.Background())
				return req
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		fmt.Sprintf("/%s/%s", openai.API_VERSION, openai.EPChatCompletions),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if !strings.Contains(r.Header.Get("Authorization"), TEST_API_KEY) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(openAIResponses["unauthorizeation"])
				return
			}

			data, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			defer r.Body.Close()

			m := map[string]any{}
			err = json.Unmarshal(data, &m)
			require.NoError(t, err)

			if m["response_format"] != nil {
				format, ok := m["response_format"].(map[string]any)
				require.True(t, ok)
				if format["type"] == "json_object" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(openAIResponses["bad_request"])
					return
				}
				t.Fatalf("unknown response format %s", format["type"])
			}

			w.WriteHeader(http.StatusOK)
			if model := m["model"].(string); model == "gpt-4" {
				w.Write(openAIResponses["chat_completion_gpt-4"])
			} else if model == "gpt-3.5-turbo" {
				w.Write(openAIResponses["chat_completion_gpt-3.5-turbo"])
			}
		},
	)
	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	srvrUrl, err := url.Parse(srvr.URL)
	require.NoError(t, err)

	cli := http.Client{Timeout: 3 * time.Second}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.Name, func(t *testing.T) {
			req := tc.NewRequestFunc(t)
			httpReq, err := req.ToHTTPRequest()
			require.NoError(t, err)
			require.NotNil(t, httpReq)

			httpReq.URL.Scheme = srvrUrl.Scheme
			httpReq.URL.Host = srvrUrl.Host

			httpResp, err := cli.Do(httpReq)
			require.NoError(t, err)
			require.NotNil(t, httpResp)
			require.Equal(t, tc.HttpStatusCode, httpResp.StatusCode)

			resp, err := openai.ParseHTTPResponse[openai.ChatCompletionsObject](httpResp)
			if tc.HttpStatusCode != http.StatusOK {
				require.Error(t, err)
				require.Equal(t, tc.HttpStatusCode, resp.StatusCode)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.HttpStatusCode, resp.StatusCode)
			require.Equal(t, "chat.completion", resp.Body.Objcet)

			saObj := openai.SentimentAnalysisObject(resp.Body)
			content, err := saObj.Content()
			require.NoError(t, err)
			for i, c := range content {
				t.Logf("content[%d]: %d", i, c)
			}
		})
	}
}

func TestContentMarshal(t *testing.T) {
	jsn := `{
		"message": {
			"role": "assistant",
			"content": "[4, 3]"
		}
	}`

	type Objcet struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}
	}

	var obj Objcet
	err := json.Unmarshal([]byte(jsn), &obj)
	require.NoError(t, err)

	t.Log(obj.Message.Role)
	t.Log(obj.Message.Content)

	var content []int
	err = json.Unmarshal([]byte(obj.Message.Content), &content)
	require.NoError(t, err)

	t.Log(content)

}
