package cohere_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	cohere "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/Cohere"
	"github.com/stretchr/testify/require"
)

const (
	TEST_API_KEY = "[[::TEST_API_KEY::]]"
)

func TestCoEmbed(t *testing.T) {
	apikey := TEST_API_KEY

	req := cohere.NewEmbedRequest(
		apikey,
		"義大利首都羅馬有了第一家專門服務狗狗的美食餐廳。菜單上的餐點按照狗狗的體重有4種不同的份量，避掉了可能誘發過敏的食材，並有專門為狗烹煮的魚料理，吸引不少愛狗人士上門消費。",
		"立法院今（1）日三讀通過《道路交通安全基本法》，以「2050年道路交通事故零死亡為目標」，未來由行政院統籌交通安全會報。而目前，以行人環境優化、加強執法等面向，來降低死傷。學者分析，台灣交通要轉換為'人本'主義，除了交通設施逐步改善，回歸教育，用路人行為改變，也是關鍵。",
		"以哈衝突延長休戰無望，以色列總理納坦雅胡1日宣稱武裝組織哈瑪斯違反協議，在休戰最後1小時用火箭攻擊以色列遭以軍攔截，雙方隨即恢復在加 薩的戰鬥狀態。而休戰最後1天內，哈瑪斯釋放8名以色列人質，以色列也放了30名巴勒斯坦囚犯。",
	)
	require.NotNil(t, req)
	require.NotNil(t, cohere.Modifier)
	require.Empty(t, req.Body.Truncate)
	require.Empty(t, req.Body.Model)
	require.Empty(t, req.Body.InputType)

	err := req.Modify(context.Background())
	require.NoError(t, err)
	require.Equal(t, cohere.TruncateEnd, req.Body.Truncate)
	require.Equal(t, cohere.EmbedModelMultilingualLightv3, req.Body.Model)
	require.Equal(t, cohere.InputTypeClustering, req.Body.InputType)

	err = req.Validate(context.Background())
	require.NoError(t, err)

	httpReq, err := req.ToHTTPRequest()
	require.NoError(t, err)
	require.NotNil(t, httpReq)
	require.Equal(t, "https://api.cohere.ai/v1/embed", httpReq.URL.String())

	f, err := os.ReadFile("./example_response/embeddings.json")
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc(
		fmt.Sprintf("/%s/%s", cohere.API_VERSION, req.Endpoint()),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if !strings.Contains(r.Header.Get("Authorization"), apikey) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "invalid api token"}`))
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

	resp, err := cohere.ParseHTTPResponse[cohere.EmbedResponseBody](httpResp)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "a5778bcc-842c-4ebf-9203-2de79e09decd", resp.Body.Id)

	for _, item := range resp.Body.Unwind() {
		t.Log(item)
	}
}

func TestPasreCoChatResponse(t *testing.T) {
	data, err := os.ReadFile("./example_response/chat.json")
	require.NoError(t, err)

	var body cohere.ChatResponseBody
	err = json.Unmarshal(data, &body)
	require.NoError(t, err)
	require.Equal(t, "d43e69ae-9710-455f-9b7a-ca4128126819", body.ResponseId.String())
	require.Equal(t, "3", body.Text)
	require.Equal(t, "e4d0bfa2-e9e2-4080-a37e-1f62ab92f2e4", body.GenerationId.String())

	require.NotNil(t, body.TokenCount)
	require.Equal(t, 349, body.TokenCount.PromptTokens)
	require.Equal(t, 1, body.TokenCount.ResponseTokens)
	require.Equal(t, 350, body.TokenCount.TotalTokens)
	require.Equal(t, 334, body.TokenCount.BilledTokens)

	require.NotNil(t, body.Meta)
	require.Equal(t, "1", body.Meta.APIVersion.Version)
	require.Equal(t, 333, body.Meta.BilledUnits.InputTokens)
	require.Equal(t, 1, body.Meta.BilledUnits.OutputTokens)
}

func TestCoChat(t *testing.T) {
	apikey := TEST_API_KEY

	req := cohere.NewChatRequest(apikey, "[^]義大利首都羅馬有了第一家專門服務狗狗的美食餐廳。菜單上的餐點按照狗狗的體重有4種不同的份量，避掉了可能誘發過敏的食材，並有專門為狗烹煮的魚料理，吸引不少愛狗人士上門消費。[$]")
	req.Body.
		AppendChatHistory("CHATBOT", cohere.SentimentAnalysisPrompt, "").
		SetModel(cohere.GenerateModelCommand).
		SetTemperature(0.001)
	require.Empty(t, req.Body.CitationQuality)
	require.Equal(t, cohere.GenerateModelCommand, req.Body.Model)

	err := req.Modify(context.Background())
	require.NoError(t, err)

	err = req.Validate(context.Background())
	require.NoError(t, err)
	require.Equal(t, "accurate", req.Body.CitationQuality)
	require.Equal(t, float32(0.001), req.Body.Temperature)

	httpReq, err := req.ToHTTPRequest()
	require.NoError(t, err)
	require.NotNil(t, httpReq)

	f, err := os.ReadFile("./example_response/chat.json")
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc(
		fmt.Sprintf("/%s/%s", cohere.API_VERSION, req.Endpoint()),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if !strings.Contains(r.Header.Get("Authorization"), apikey) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "invalid api token"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(f)
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

	resp, err := cohere.ParseHTTPResponse[cohere.ChatResponseBody](httpResp)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "d43e69ae-9710-455f-9b7a-ca4128126819", resp.Body.ResponseId.String())
	require.Equal(t, "3", resp.Body.Text)
	require.Equal(t, "e4d0bfa2-e9e2-4080-a37e-1f62ab92f2e4", resp.Body.GenerationId.String())

}
