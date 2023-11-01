package gnews_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/GNews"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

const TEST_API_KEY = "0x0x000x0000x00x0xx0xxxxx0xx0000"

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
					apiResponse, err := cli.ParseResponse(respJson, 0)
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
				apiresp, err := cli.ParseHTTPResponse(resp, 0)
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

	mux := chi.NewRouter()
	mux.Get("/"+strings.Join([]string{
		cli.API_PATH,
		cli.API_VERSION,
		cli.EPTopHeadlines}, "/"),
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
	t.Log("/" + strings.Join([]string{
		cli.API_PATH,
		cli.API_VERSION,
		cli.EPTopHeadlines}, "/"))
	srvr := httptest.NewServer(mux)
	defer srvr.Close()

	cPars := make(chan *service.NewsCreateRequest)
	go func(c <-chan *service.NewsCreateRequest) {
		for p := range c {
			require.NotNil(t, p)
			require.NotEmpty(t, p.Title)
			require.NotEmpty(t, p.Md5Hash)
		}
	}(cPars)
	wg := &sync.WaitGroup{}

	q, err := h.Handle(TEST_API_KEY, pf)
	require.NoError(t, err)
	require.NotNil(t, q)

	qs := q.Encode()
	require.Contains(t, qs, "country="+strings.Join(pf.Country, "%2C"))
	require.Contains(t, qs, "lang="+strings.Join(pf.Language, "%2C"))
	require.NotContains(t, qs, "from=0001-01-01T23%3A59%3A59Z")
	require.NotContains(t, qs, "to=0001-01-01T23%3A59%3A59Z")

	r, err := q.ToHttpRequest()
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
		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		apiResp, err := h.Parse(resp)
		require.NoError(t, err)
		require.NotNil(t, apiResp)
		require.Equal(t, "success", apiResp.GetStatus())

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
