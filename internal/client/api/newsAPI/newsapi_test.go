package newsapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/newsAPI"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/newsapi"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

const TEST_API_KEY = "00000xxxx00000xxxxx0xx0000x00xxx"

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
					apiResponse, err := cli.ParseResponse(respJson, 0)
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
				apiresp, err := cli.ParseHTTPResponse(resp, 0)
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

func TestHeadlinesHandler(t *testing.T) {
	h := cli.TopHeadlinesHandler{}

	pf := srv.NEWSAPITopHeadlines{
		Keyword:  "世大運",
		Country:  srv.Taiwan,
		Category: srv.Sports,
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

	mux := chi.NewRouter()
	mux.Get("/"+strings.Join([]string{
		cli.API_VERSION,
		cli.EPTopHeadlines}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			if r.URL.Query().Get("apikey") != TEST_API_KEY {
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

	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	cPars := make(chan *service.NewsCreateRequest)
	go func(c <-chan *service.NewsCreateRequest) {
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
	require.Contains(t, qs, "country="+pf.Country)
	require.Contains(t, qs, "category="+pf.Category)

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

	rs := []*http.Request{r}
	for i := 0; i < len(rs); i++ {
		resp, err := http.DefaultClient.Do(rs[i])
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		apiResp, err := h.Parse(resp)
		require.NoError(t, err)
		require.NotNil(t, apiResp)
		require.Equal(t, "ok", apiResp.GetStatus())

		if i < len(tc.Filename)-1 {
			require.True(t, apiResp.HasNext())
			r, err = apiResp.NextPageRequest(nil)
			require.NoError(t, err)
			rs = append(rs, r)
		} else {
			require.False(t, apiResp.HasNext())
		}

		wg.Add(1)
		go apiResp.ToNews(context.Background(), wg, cPars)
		require.Greater(t, len(tc.Filename), i)
	}

	wg.Wait()
	close(cPars)
}

func TestEverythingHandler(t *testing.T) {
	h := cli.EverythingHandler{}

	ft, err := time.Parse(time.DateOnly, "2023-06-30")
	require.NoError(t, err)

	pf := srv.NEWSAPIEverything{
		SearchIn: pageform.SearchIn{
			InTitle:       true,
			InDescription: true,
			InContent:     true},
		TimeRange: pageform.TimeRange{
			Form: ft,
		},
		Keyword:  "worldcoin",
		Language: []string{srv.Chinese},
	}

	type testCase struct {
		Name     string
		Filename map[string]string
	}

	tc := testCase{
		Name: "latest news",
		Filename: map[string]string{
			"":  "example_response/everything_1.json",
			"2": "example_response/everything_2.json",
			"3": "example_response/everything_3.json",
			"4": "example_response/everything_4.json",
		},
	}

	mux := chi.NewRouter()
	mux.Get("/"+strings.Join([]string{
		cli.API_VERSION,
		cli.EPEverything}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			if r.URL.Query().Get("apikey") != TEST_API_KEY {
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

	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	cPars := make(chan *service.NewsCreateRequest)
	go func(c <-chan *service.NewsCreateRequest) {
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
	require.NotContains(t, qs, "searchIn=")
	require.NotContains(t, qs, "to=")
	require.Contains(t, qs, "from=2023-06-30")
	require.Contains(t, qs, fmt.Sprintf("q=%s", pf.Keyword))

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

	rs := []*http.Request{r}
	for i := 0; i < len(rs); i++ {
		resp, err := http.DefaultClient.Do(rs[i])
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		apiResp, err := h.Parse(resp)
		require.NoError(t, err)
		require.NotNil(t, apiResp)
		require.Equal(t, "ok", apiResp.GetStatus())

		if i < len(tc.Filename)-1 {
			require.True(t, apiResp.HasNext())
			r, err = apiResp.NextPageRequest(nil)
			require.NoError(t, err)
			rs = append(rs, r)
		} else {
			require.False(t, apiResp.HasNext())
		}

		wg.Add(1)
		go apiResp.ToNews(context.Background(), wg, cPars)
		require.Greater(t, len(tc.Filename), i)
	}

	wg.Wait()
	close(cPars)
}
