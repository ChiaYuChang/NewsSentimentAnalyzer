package newsdata_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/NEWSDATA"
	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/NEWSDATA"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

const TEST_API_KEY = "pub_00000x0000x00xx0x0000xxxx00x00xx0xx002"

func TestParseResponseBody(t *testing.T) {
	type testCase struct {
		name        string
		fileName    string
		hasNextPage bool
	}

	tcs := []testCase{
		{
			name:        "without next page",
			fileName:    "example_response/001.json",
			hasNextPage: false,
		},
		{
			name:        "with next page",
			fileName:    "example_response/002.json",
			hasNextPage: true,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.name,
			func(t *testing.T) {
				respJson, err := os.ReadFile(tc.fileName)
				require.NoError(t, err)

				apiResponse, err := cli.ParseResponse(respJson)
				require.NoError(t, err)

				if tc.hasNextPage {
					require.NotEmpty(t, apiResponse.NextPage)
				} else {
					require.Empty(t, apiResponse.NextPage)
				}
			},
		)
	}
}

func TestParseResponseBody1(t *testing.T) {
	type testCase struct {
		Name     string
		FileName []string
	}

	tcs := []testCase{
		{
			Name: "latest news",
			FileName: []string{
				"example_response/latest_news_1.json",
				"example_response/latest_news_2.json",
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.Name,
			func(t *testing.T) {
				for i, fn := range tc.FileName {
					respJson, err := os.ReadFile(fn)
					require.NoError(t, err)

					apiResponse, err := cli.ParseResponse(respJson)
					require.NoError(t, err)
					require.NotNil(t, apiResponse)

					t.Log(apiResponse.Len())
					if i < len(tc.FileName)-1 {
						require.True(t, apiResponse.HasNext())
					} else {
						require.False(t, apiResponse.HasNext())
					}
				}
			},
		)
	}
}

func TestParsErroreResponseBody(t *testing.T) {
	type testCase struct {
		File        string
		Status      string
		HttpCode    int
		Message     string
		MessageCode string
	}

	tcs := []testCase{
		{
			File:        "example_response/error_Unauthorized.json",
			Status:      "error",
			HttpCode:    http.StatusUnauthorized,
			Message:     "API key is invalid",
			MessageCode: "Unauthorized",
		},
		{
			File:        "example_response/error_FilterLimitExceed.json",
			Status:      "error",
			HttpCode:    http.StatusUnprocessableEntity,
			Message:     "Number of countries selected cannot be greater than 1",
			MessageCode: "FilterLimitExceed",
		},
		{
			File:        "example_response/error_UnsupportedFilter.json",
			Status:      "error",
			HttpCode:    http.StatusUnprocessableEntity,
			Message:     "Sorry! No such category exists in our database.",
			MessageCode: "UnsupportedFilter",
		},
	}

	for i := range tcs {
		tc := tcs[i]
		respJson, err := os.ReadFile(tc.File)
		require.NoError(t, err)

		apiErrResponse, err := cli.ParseErrorResponse(respJson)
		require.NoError(t, err)
		require.Error(t, apiErrResponse.ToError(tc.HttpCode))
		require.Equal(t, tc.Status, apiErrResponse.Status)
		require.Equal(t, tc.Message, apiErrResponse.Result["message"])
		require.Equal(t, tc.MessageCode, apiErrResponse.Result["code"])
		require.Contains(t, apiErrResponse.String(), tc.Message)
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
			Name:         "without next page",
			FileName:     "example_response/001.json",
			RoutePattern: "/001",
			StatusCode:   http.StatusOK,
		},
		{
			Name:         "with next page",
			FileName:     "example_response/002.json",
			RoutePattern: "/002",
			StatusCode:   http.StatusOK,
		},
		{
			Name:         "error",
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

func TestLatestNewsHandler(t *testing.T) {
	h := cli.LatestNewsHandler{}

	pf := srv.NEWSDATAIOLatestNews{
		Keyword:  "Typhoon AND Taiwan",
		Language: []string{srv.Chinese, srv.English},
		Country:  []string{srv.Taiwan, srv.UnitedStates},
	}

	type testCase struct {
		Name     string
		Filename map[string]string
	}

	tc := testCase{
		Name: "latest news",
		Filename: map[string]string{
			"": "example_response/latest_news_1.json",
			"16903560911db4c680d9b2461ebe967794a90f1b3f": "example_response/latest_news_2.json",
		},
	}

	mux := chi.NewRouter()
	mux.Get("/"+strings.Join([]string{
		newsdata.API_PATH,
		newsdata.API_VERSION,
		newsdata.EPLatestNews}, "/"),
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
			j, _ := os.ReadFile(tc.Filename[page])
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(j)
		})

	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	cPars := make(chan *model.CreateNewsParams)
	go func(c <-chan *model.CreateNewsParams) {
		for p := range c {
			require.NotEmpty(t, p.Title)
			require.NotEmpty(t, p.Md5Hash)
		}
	}(cPars)
	wg := &sync.WaitGroup{}

	q, err := h.Handle(TEST_API_KEY, pf)
	require.NoError(t, err)
	require.NotNil(t, q)

	qs := q.Params().ToQueryString()
	require.Contains(t, qs, "country=")
	require.Contains(t, qs, "language=")

	r, err := q.ToRequest()
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, cli.API_HOST, r.URL.Host)
	require.Equal(t, cli.API_SCHEME, r.URL.Scheme)
	require.Contains(t, r.URL.RawQuery, "apikey="+TEST_API_KEY)

	srvrUrl, err := url.Parse(srvr.URL)
	require.NoError(t, err)

	r.URL.Scheme = srvrUrl.Scheme
	r.URL.Host = srvrUrl.Host

	// t.Log(r.URL)
	resp, err := http.DefaultClient.Do(r)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp, err := h.Parse(resp)
	require.NoError(t, err)
	require.NotNil(t, apiResp)
	require.Equal(t, "success", apiResp.GetStatus())
	require.True(t, apiResp.HasNext())

	wg.Add(1)
	go apiResp.ToNews(context.Background(), wg, cPars)

	r, err = apiResp.NextPageRequest(nil)
	require.NoError(t, err)
	require.Contains(t, r.URL.String(), "16903560911db4c680d9b2461ebe967794a90f1b3f")

	resp, err = http.DefaultClient.Do(r)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp, err = h.Parse(resp)
	require.NoError(t, err)
	require.NotNil(t, apiResp)
	require.Equal(t, "success", apiResp.GetStatus())
	require.False(t, apiResp.HasNext())

	wg.Add(1)
	go apiResp.ToNews(context.Background(), wg, cPars)

	wg.Wait()
	close(cPars)
}

func TestNewsArchiveHandler(t *testing.T) {
	h := newsdata.NewsArchiveHandler{}

	pf := srv.NEWSDATAIONewsArchive{
		Keyword:  "Typhoon",
		Language: []string{srv.Chinese, srv.English},
		Country:  []string{srv.Taiwan, srv.UnitedStates},
	}

	q, err := h.Handle(TEST_API_KEY, pf)
	require.NoError(t, err)
	require.NotNil(t, q)

	qs := q.Params().ToQueryString()
	t.Log(qs)
	require.Contains(t, qs, "country=")
	require.Contains(t, qs, "language=")
	require.NotContains(t, qs, "from=0001-01-01")
	require.NotContains(t, qs, "to=0001-01-01")
}

func TestNewsSourcesHandler(t *testing.T) {
	h := newsdata.NewsSourcesHandler{}

	pf := srv.NEWSDATAIONewsSources{
		Language: []string{srv.Chinese, srv.English},
		Country:  []string{srv.Taiwan, srv.UnitedStates},
		Category: []string{srv.Business, srv.Environment},
	}

	q, err := h.Handle(TEST_API_KEY, pf)
	require.NoError(t, err)
	require.NotNil(t, q)

	qs := q.Params().ToQueryString()
	require.Contains(t, qs, "country=")
	require.Contains(t, qs, "language=")
	require.Contains(t, qs, "category=")

	r, err := q.ToRequest()
	require.NoError(t, err)
	require.NotNil(t, r)

	t.Log(r.URL.String())
}
