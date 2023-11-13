package googlecse_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/GoogleCSE"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GoogleCSE"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const TEST_SEARCH_ENGINE_ID = "[[:TEST_SEARCH_ENGINE_ID:]]"
const TEST_API_KEY = "[[:TEST_API_KEY:]]"

var TEST_USER_ID, _ = uuid.Parse("741428c7-1ae0-4622-b615-9d44a141ff23")

func TestUnmarshalJsonResponse(t *testing.T) {
	b, err := os.ReadFile("./example_response/002.json")
	require.NoError(t, err)

	var resp cli.Response
	err = json.Unmarshal(b, &resp)
	require.NoError(t, err)

	require.Equal(t, "customsearch#search", resp.Kind)
	require.Equal(t, "application/json", resp.URL.Type)
	require.Equal(t, "customsearch#result", resp.Items[0].Kind)
	require.Equal(t, "日本海底火山8月噴發浮石漂到綠島海域船隻注意| 生活| 中央社CNA", resp.Items[0].Title)
	require.Equal(t, "https://www.cna.com.tw/news/ahel/202111300051.aspx", resp.Items[0].Link.String())

	require.Equal(t, "日本海底火山8月噴發 浮石漂到綠島海域船隻注意 | 生活 | 中央社 CNA", resp.Items[0].PageMap.Title())
	require.Equal(t, "日本小笠原群島海底火山「福德岡之場」今年8月噴發後，產生大量浮石最近隨著洋流漂到綠島海域；水產試驗所提醒海上航行船隻互相通報浮石位置，避免浮石堵塞船舶冷卻水管造成損壞。", resp.Items[0].PageMap.Description())
	require.Equal(t, "https://www.cna.com.tw/news/ahel/202111300051.aspx", resp.Items[0].PageMap.Link())
	require.Equal(t, "", resp.Items[0].PageMap.Category())

	pt, _ := time.Parse(time.RFC3339, "2021-11-30T10:16:00+08:00")
	require.Equal(t, pt.UTC(), resp.Items[0].PageMap.PubDate())
}

func TestUnmarshalError(t *testing.T) {
	type testCase struct {
		Name     string
		Code     int
		FileName string
		Status   string
		ErrMsg   string
		Reason   string
	}

	tcs := []testCase{
		{
			Name:     "400 Error",
			Code:     http.StatusBadRequest,
			FileName: "./example_response/error_400.json",
			Status:   "INVALID_ARGUMENT",
			ErrMsg:   "Request contains an invalid argument.",
			Reason:   "badRequest",
		},
		{
			Name:     "403 Error",
			Code:     http.StatusForbidden,
			FileName: "./example_response/error_403.json",
			Status:   "PERMISSION_DENIED",
			ErrMsg:   "The request is missing a valid API key.",
			Reason:   "forbidden",
		},
	}

	for i := range tcs {
		tc := tcs[i]
		b, err := os.ReadFile(tc.FileName)
		require.NoError(t, err)

		var errResp cli.ErrorResponse
		err = json.Unmarshal(b, &errResp)
		require.NoError(t, err)

		require.Equal(t, tc.Status, errResp.Error.Status)
		require.Equal(t, 1, len(errResp.Error.Errors))
		require.Equal(t, tc.Reason, errResp.Error.Errors[0].Reason)
		require.Equal(t, tc.ErrMsg, errResp.Error.Errors[0].Message)
		require.Equal(t, tc.Code, errResp.Error.Code)
	}
}

func TestCESHandler(t *testing.T) {
	tc := struct {
		Name     string
		Handler  client.PageFormHandler
		PageForm pageform.PageForm
		Filename map[string]string
		Endpoint string
	}{
		Name:    "Google CSE",
		Handler: cli.CSEHandler{},
		PageForm: srv.GoogleCSE{
			Keyword:        "日本",
			SearchEngineID: TEST_SEARCH_ENGINE_ID,
		},
		Filename: map[string]string{
			"":   "example_response/001.json",
			"11": "example_response/002.json",
			"21": "example_response/003.json",
		},
		Endpoint: cli.EPCustomSearch,
	}

	mux := chi.NewRouter()

	p := path.Join(cli.API_PATH, cli.API_VERSION)
	if tc.Endpoint != cli.EPCustomSearch {
		p += "/" + tc.Endpoint
	}

	mux.Get("/"+p, func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()

		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("key") != TEST_API_KEY {
			j, _ := os.ReadFile("example_response/error_403.json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(j)
			return
		}

		start := r.URL.Query().Get("start")
		j, _ := os.ReadFile(tc.Filename[start])
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	})

	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	srvrUrl, err := url.Parse(srvr.URL)
	require.NoError(t, err)

	t.Run(tc.Name, func(t *testing.T) {
		Cache := map[string]*api.PreviewCache{}
		NItem := 0

		// user's first request
		req, err := tc.Handler.Handle(TEST_API_KEY, tc.PageForm)
		require.NoError(t, err)
		require.NotNil(t, req)

		ckey, cache := req.ToPreviewCache(TEST_USER_ID)
		Cache[ckey] = cache
		require.NotZero(t, ckey)
		require.NotNil(t, cache)

		// should stop when i == len(tc.Filename)
		for i := 1; i <= len(tc.Filename)+1; i++ {
			httpReq, err := req.ToHttpRequest()
			require.NoError(t, err)
			require.NotNil(t, httpReq)
			httpReq.URL.Scheme = srvrUrl.Scheme
			httpReq.URL.Host = srvrUrl.Host

			require.Contains(t, httpReq.URL.String(), tc.Endpoint)

			if i == 1 {
				require.NotContains(t, httpReq.URL.RawQuery, "start=")
			} else {
				require.Contains(t, httpReq.URL.RawQuery, "start=")
			}

			httpResp, err := http.DefaultClient.Do(httpReq)
			require.NoError(t, err)
			require.NotNil(t, httpResp)
			require.Equal(t, http.StatusOK, httpResp.StatusCode)

			resp, err := tc.Handler.Parse(httpResp)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, "success", resp.GetStatus())

			next, prev := resp.ToNewsItemList()
			// update cache
			cache, ok := Cache[ckey]
			require.True(t, ok)
			cache.SetNextPage(next)

			if next.Equal(api.IntLastPageToken) {
				require.Nil(t, prev)
				require.Equal(t, 0, len(prev))
				require.Equal(t, len(tc.Filename), i)
				break
			} else {
				require.Less(t, 0, len(prev))
				NItem += len(prev)

				// append items to cache
				cache.NewsItem = append(cache.NewsItem, prev...)
				Cache[ckey] = cache
			}

			// next client request comes in
			// build request from cache
			req, err = cli.RequestFromPreviewCache(cache.Query)
			require.NoError(t, err)
		}
		require.Equal(t, NItem, len(Cache[ckey].NewsItem))
	})
}
