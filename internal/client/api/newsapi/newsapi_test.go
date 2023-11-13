package newsapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/newsapi"
	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/newsapi"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
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
			Name:     "everything page 1",
			FileName: "example_response/everything_1.json",
			IsErr:    false,
		},
		{
			Name:     "everything page 2",
			FileName: "example_response/everything_2.json",
			IsErr:    false,
		},

		{
			Name:     "everything page 3",
			FileName: "example_response/everything_3.json",
			IsErr:    false,
		},
		{
			Name:     "everything page 4",
			FileName: "example_response/everything_4.json",
			IsErr:    false,
		},
		{
			Name:     "top headlines page 1",
			FileName: "example_response/top-headlines_1.json",
			IsErr:    false,
		},
		{
			Name:     "top headlines page 2",
			FileName: "example_response/top-headlines_1.json",
			IsErr:    false,
		},
		{
			Name:     "top headlines page 3",
			FileName: "example_response/top-headlines_1.json",
			IsErr:    false,
		},

		{
			Name:     "Unauthorized error",
			FileName: "example_response/error_ParameterInvalid.json",
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
					require.NotNil(t, apiErrResponse)
					require.Empty(t, apiErrResponse.Code)
				} else {
					apiResponse, err := cli.ParseResponse(respJson)
					require.NoError(t, err)
					require.NotNil(t, apiResponse)
					require.NotEmpty(t, apiResponse)
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
			Name:         "everything page 1",
			FileName:     "example_response/everything_1.json",
			RoutePattern: "/everything",
			StatusCode:   http.StatusOK,
		},
		{
			Name:         "top headlines page 1",
			FileName:     "example_response/top-headlines_1.json",
			RoutePattern: "/top-headlines",
			StatusCode:   http.StatusOK,
		},
		{
			Name:         "Unauthorized error",
			FileName:     "example_response/error_ParameterInvalid.json",
			RoutePattern: "/err",
			StatusCode:   http.StatusUnauthorized,
		},
		{
			Name:         "API key invalid error",
			FileName:     "example_response/error_ApiKeyInvalid.json",
			RoutePattern: "/err",
			StatusCode:   http.StatusUnauthorized,
		},
		{
			Name:         "Api key missing error",
			FileName:     "example_response/error_ApiKeyMissing.json",
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
					require.Equal(t, apiresp.GetStatus(), "ok")
				} else {
					require.Error(t, err)
				}
			},
		)
	}
}

func TestHandlers(t *testing.T) {
	type testCase struct {
		Name     string
		Handler  client.PageFormHandler
		PageForm pageform.PageForm
		Endpoint string
		Filename map[string]string
	}

	ft, err := time.Parse(time.DateOnly, "2023-06-30")
	require.NoError(t, err)

	tcs := []testCase{
		{
			Name:    "top-headlines",
			Handler: cli.TopHeadlinesHandler{},
			PageForm: srv.NEWSAPITopHeadlines{
				Keyword:  "世大運",
				Country:  srv.Taiwan,
				Category: srv.Sports,
			},
			Filename: map[string]string{
				"":  "example_response/top-headlines_1.json",
				"2": "example_response/top-headlines_2.json",
				"3": "example_response/top-headlines_3.json",
			},
			Endpoint: cli.EPTopHeadlines,
		},
		{
			Name:    "everything",
			Handler: cli.EverythingHandler{},
			PageForm: srv.NEWSAPIEverything{
				SearchIn: pageform.SearchIn{
					InTitle:       true,
					InDescription: true,
					InContent:     true},
				TimeRange: pageform.TimeRange{
					Form: ft,
				},
				Keyword:  "worldcoin",
				Language: srv.Chinese,
			},
			Filename: map[string]string{
				"":  "example_response/everything_1.json",
				"2": "example_response/everything_2.json",
				"3": "example_response/everything_3.json",
				"4": "example_response/everything_4.json",
			},
			Endpoint: cli.EPEverything,
		},
	}

	// setup test server
	mux := chi.NewRouter()
	for i := range tcs {
		tc := tcs[i]
		mux.Get("/"+path.Join(cli.API_VERSION, tc.Endpoint),
			func(w http.ResponseWriter, r *http.Request) {
				_ = r.ParseForm()
				if !(r.URL.Query().Get("apikey") != TEST_API_KEY ||
					r.Header.Get("X-Api-Key") != TEST_API_KEY ||
					r.Header.Get("Authorization") != TEST_API_KEY) {
					j, _ := os.ReadFile("example_response/error_ApiKeyInvalid.json")
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
	}
	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	srvrUrl, err := url.Parse(srvr.URL)
	require.NoError(t, err)

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.Name, func(t *testing.T) {
			// mock a redis db
			Cache := map[string]*api.PreviewCache{}
			NItem := 0

			// user first request
			req, err := tc.Handler.Handle(TEST_API_KEY, tc.PageForm)
			require.NoError(t, err)
			require.NotNil(t, req)

			ckey, cache := req.ToPreviewCache(TEST_USER_ID)
			Cache[ckey] = cache
			require.NotZero(t, ckey)
			require.NotNil(t, cache)

			for i := 1; i <= 10; i++ {
				httpReq, err := req.ToHttpRequest()
				require.NoError(t, err)
				require.NotNil(t, httpReq)
				httpReq.URL.Scheme = srvrUrl.Scheme
				httpReq.URL.Host = srvrUrl.Host

				require.Contains(t, httpReq.URL.String(), tc.Endpoint)

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
				resp, err := tc.Handler.Parse(httpResp)
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "ok", resp.GetStatus())

				next, prev := resp.ToNewsItemList()
				// update cache
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
				req, err = newsapi.RequestFromPreviewCache(cache)
				require.NoError(t, err)
			}
			require.Equal(t, NItem, len(Cache[ckey].NewsItem))
		})
	}
}
