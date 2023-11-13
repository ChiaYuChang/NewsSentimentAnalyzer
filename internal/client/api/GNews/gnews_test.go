package gnews_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/GNews"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const TEST_API_KEY = "[[:TEST_API_KEY:]]"

var TEST_USER_ID, _ = uuid.Parse("741428c7-1ae0-4622-b615-9d44a141ff23")

func TestParseResponseBody(t *testing.T) {
	type testCase struct {
		Name     string
		FileName string
		IsErr    bool
	}

	tcs := []testCase{
		{
			Name:     "top-headlines",
			FileName: "example_response/001.json",
			IsErr:    false,
		},
		{
			Name:     "top-headlines page 1",
			FileName: "example_response/top-headlines_1.json",
			IsErr:    false,
		},
		{
			Name:     "top-headlines page 3",
			FileName: "example_response/top-headlines_3.json",
			IsErr:    false,
		},
		{
			Name:     "Unauthorized error",
			FileName: "example_response/error_Unauthorized.json",
			IsErr:    true,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.Name,
			func(t *testing.T) {
				respJson, err := os.ReadFile(tc.FileName)
				require.NoError(t, err)

				if tc.IsErr {
					apiErrResponse, err := cli.ParseErrorResponse(respJson)
					require.NoError(t, err)
					require.NotEmpty(t, apiErrResponse.Error)
				} else {
					apiResponse, err := cli.ParseResponse(respJson)
					require.NoError(t, err)
					require.NotNil(t, apiResponse)
				}
			},
		)
	}
}

func TestParseResponse(t *testing.T) {
	type testCase struct {
		Name         string
		FileName     string
		RoutePattern string
		StatusCode   int
	}

	tcs := []testCase{
		{
			Name:         "OK",
			FileName:     "example_response/001.json",
			RoutePattern: "/001",
			StatusCode:   http.StatusOK,
		},
		{
			Name:         "top-headlines first page",
			FileName:     "example_response/top-headlines_1.json",
			RoutePattern: "/th-h",
			StatusCode:   http.StatusOK,
		},
		{
			Name:         "top-headlines last page",
			FileName:     "example_response/top-headlines_3.json",
			RoutePattern: "/th-t",
			StatusCode:   http.StatusOK,
		},
		{
			Name:         "Unauthorized error",
			FileName:     "example_response/error_Unauthorized.json",
			RoutePattern: "/err",
			StatusCode:   http.StatusUnauthorized,
		},
	}

	mux := chi.NewRouter()
	for i := range tcs {
		tc := tcs[i]
		mux.Get(tc.RoutePattern, func(w http.ResponseWriter, r *http.Request) {
			j, _ := os.ReadFile(tc.FileName)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(tc.StatusCode)
			w.Write(j)
		})
	}

	srvr := httptest.NewServer(mux)
	defer srvr.Close()
	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.Name,
			func(t *testing.T) {
				resp, err := http.DefaultClient.Get(srvr.URL + tc.RoutePattern)
				require.NoError(t, err)
				apiresp, err := cli.ParseHTTPResponse(resp)
				if tc.StatusCode == http.StatusOK {
					require.NoError(t, err)
					require.Equal(t, apiresp.GetStatus(), "success")
				} else {
					require.Error(t, err)
				}
			},
		)
	}
}

func TestHeadlinesHandler(t *testing.T) {
	h := cli.TopHeadlinesHandler{}

	pf := srv.GNewsHeadlines{
		Keyword:  "Typhoon",
		Language: []string{srv.Chinese},
		Country:  []string{srv.Taiwan, srv.UnitedStates},
		Category: []string{srv.General},
	}

	type testCase struct {
		Name     string
		Filename map[string]string
	}

	tc := testCase{
		Name: "latest news",
		Filename: map[string]string{
			"":  "example_response/top-headlines_1.json",
			"2": "example_response/top-headlines_2.json",
			"3": "example_response/top-headlines_3.json",
		},
	}

	p := path.Join(
		cli.API_PATH,
		cli.API_VERSION,
		cli.EPTopHeadlines,
	)

	// new test server
	mux := chi.NewRouter()
	mux.Get("/"+p,
		func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			if r.URL.Query().Get("apikey") != TEST_API_KEY {
				j, _ := os.ReadFile("example_response/error_Unauthorized.json")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(j)
				return
			}

			page := r.URL.Query().Get("page")
			w.Header().Set("Content-Type", "application/json")
			j, _ := os.ReadFile(tc.Filename[page])
			w.WriteHeader(http.StatusOK)
			w.Write(j)
		})
	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	srvrUrl, err := url.Parse(srvr.URL)
	require.NoError(t, err)

	// mock a database
	Cache := map[string]*api.PreviewCache{}
	NItem := 0

	// user's first request
	req, err := h.Handle(TEST_API_KEY, pf)
	require.NoError(t, err)
	require.NotNil(t, req)

	ckey, cache := req.ToPreviewCache(TEST_USER_ID)
	Cache[ckey] = cache
	require.NotZero(t, ckey)
	require.NotNil(t, cache)

	for i := 1; i <= 10; i++ {
		// make API request
		httpReq, err := req.ToHttpRequest()
		require.NoError(t, err)
		require.NotNil(t, httpReq)
		httpReq.URL.Scheme = srvrUrl.Scheme
		httpReq.URL.Host = srvrUrl.Host

		require.Contains(t, httpReq.URL.String(), cli.EPTopHeadlines)
		if i == 1 {
			// page=1 is default
			require.NotContains(t, httpReq.URL.RawQuery, "page=")
		} else {
			require.Contains(t, httpReq.URL.RawQuery, fmt.Sprintf("page=%d", i))
		}

		httpResp, err := http.DefaultClient.Do(httpReq)
		require.NoError(t, err)
		require.NotNil(t, httpResp)
		require.Equal(t, http.StatusOK, httpResp.StatusCode)

		// parse response
		resp, err := h.Parse(httpResp)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "success", resp.GetStatus())

		next, prev := resp.ToNewsItemList()

		// update next token in cache
		cache, ok := Cache[ckey]
		require.True(t, ok)
		cache.SetNextPage(next)

		if next.Equal(api.IntLastPageToken) {
			// last page
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
		req, err = cli.RequestFromPreviewCache(cache)
		require.NoError(t, err)
	}
	require.Equal(t, NItem, len(Cache[ckey].NewsItem))
}

func TestSep(t *testing.T) {
	text := "...構。（資料照，路透）\n2023/05/25 11:11\n首次上稿...颱風眼也能用肉眼清晰可見。\n隸屬美國國家航空暨太空總署（NASA）的太空人海因斯（Bob Hines），推特帳號..."
	s := strings.ReplaceAll(text, "。\n\r", "\n")
	s = strings.ReplaceAll(s, "。\n", "。"+global.SEPToken)
	s = strings.ReplaceAll(s, "\n", "")
	s = global.CLSToken + s
	t.Log(s)
}
