package newsdata_test

import (
	"context"
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
	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/NEWSDATA"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/NEWSDATA"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

const TEST_API_KEY = "[[:TEST_API_KEY:]]"

var TEST_USER_ID, _ = uuid.Parse("741428c7-1ae0-4622-b615-9d44a141ff23")

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
		require.Error(t, apiErrResponse.ToEcError(tc.HttpCode))
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
	type testCase struct {
		Name         string
		Handler      client.Handler
		PageForm     pageform.PageForm
		RawQueryTest func(t *testing.T, rawquery string)
		Endpoint     string
		Filename     map[string]string
		NItem        int
	}

	tcs := []testCase{{
		Name:    "lastest-news",
		Handler: cli.LatestNewsHandler{},
		PageForm: srv.NEWSDATAIOLatestNews{
			Keyword:  "Typhoon AND Taiwan",
			Language: []string{srv.Chinese, srv.English},
			Country:  []string{srv.Taiwan, srv.UnitedStates},
		},
		RawQueryTest: nil,
		Filename: map[string]string{
			"": "example_response/latest_news_1.json",
			"16903560911db4c680d9b2461ebe967794a90f1b3f": "example_response/latest_news_2.json",
		},
		Endpoint: cli.EPLatestNews,
		NItem:    13,
	}}

	mux := chi.NewRouter()
	for i := range tcs {
		tc := tcs[i]
		mux.Get("/"+path.Join(cli.API_PATH, cli.API_VERSION, tc.Endpoint),
			func(w http.ResponseWriter, r *http.Request) {
				_ = r.ParseForm()

				w.Header().Set("Content-Type", "application/json")
				if !(r.URL.Query().Get("apikey") == TEST_API_KEY ||
					r.Header.Get("X-ACCESS-KEY") == TEST_API_KEY) {
					j, _ := os.ReadFile("example_response/error_Unauthorized.json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(j)
					return
				}

				page := r.URL.Query().Get("page")
				j, _ := os.ReadFile(tc.Filename[page])
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
			Cache := map[string]*api.PreviewCache{}

			// user's first request
			ckey, cache, err := tc.Handler.Handle(TEST_API_KEY, TEST_USER_ID, tc.PageForm)
			require.NoError(t, err)
			require.NotNil(t, cache)
			require.NotEmpty(t, ckey)

			Cache[ckey] = cache
			require.NotZero(t, ckey)
			require.NotNil(t, cache)

			// should stop when i == len(tc.Filename)
			for i := 1; i <= len(tc.Filename)+1; i++ {
				req, err := tc.Handler.RequestFromCacheQuery(cache.Query)
				require.NoError(t, err)

				httpReq, err := req.ToHttpRequest()
				require.NoError(t, err)
				require.NotNil(t, httpReq)
				httpReq.URL.Scheme = srvrUrl.Scheme
				httpReq.URL.Host = srvrUrl.Host

				require.Contains(t, httpReq.URL.String(), tc.Endpoint)

				if i == 1 {
					require.NotContains(t, httpReq.URL.RawQuery, "page=")
				} else {
					require.Contains(t, httpReq.URL.RawQuery, "page=")
				}

				if tc.RawQueryTest != nil {
					tc.RawQueryTest(t, httpReq.URL.RawQuery)
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
				if len(prev) > 0 {
					// append items to cache
					cache.NewsItem = append(cache.NewsItem, prev...)
					Cache[ckey] = cache
				}

				if next.Equal(api.StrLastPageToken) {
					break
				}

				// next client request comes in
				// build request from cache
				req, err = cli.RequestFromCacheQuery(cache.Query)
				require.NoError(t, err)
			}
			require.Equal(t, tc.NItem, len(Cache[ckey].NewsItem))
		})
	}
}

func TestNewsArchiveHandler(t *testing.T) {
	h := cli.NewsArchiveHandler{}

	pf := srv.NEWSDATAIONewsArchive{
		Keyword:  "Typhoon",
		Language: []string{srv.Chinese, srv.English},
		Country:  []string{srv.Taiwan, srv.UnitedStates},
	}

	ckey, cache, err := h.Handle(TEST_API_KEY, TEST_USER_ID, pf)
	require.NoError(t, err)
	require.NotNil(t, cache)
	require.NotEmpty(t, ckey)

	qs := cache.Query.RawQuery
	require.Contains(t, qs, "country=")
	require.Contains(t, qs, "language=")
	require.NotContains(t, qs, "from=0001-01-01")
	require.NotContains(t, qs, "to=0001-01-01")
}

func TestNewsSourcesHandler(t *testing.T) {
	h := cli.NewsSourcesHandler{}

	pf := srv.NEWSDATAIONewsSources{
		Language: []string{srv.Chinese, srv.English},
		Country:  []string{srv.Taiwan, srv.UnitedStates},
		Category: []string{srv.Business, srv.Environment},
	}

	ckey, cache, err := h.Handle(TEST_API_KEY, TEST_USER_ID, pf)
	require.NoError(t, err)
	require.NotNil(t, cache)
	require.NotEmpty(t, ckey)

	qs := cache.Query.RawQuery
	require.Contains(t, qs, "country=")
	require.Contains(t, qs, "language=")
	require.Contains(t, qs, "category=")

	req, err := h.RequestFromCacheQuery(cache.Query)
	require.NoError(t, err)
	r, err := req.ToHttpRequest()
	require.NoError(t, err)
	require.NotNil(t, r)

	t.Log(r.URL.String())
}

func TestWriteToRedis(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Log("skip test, .env file not found")
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Log("skip test, failed to connect to redis")
		return
	}

	rh := rejson.NewReJSONHandler()
	defer func() {
		err := rdb.Close()
		require.NoError(t, err)
	}()

	rh.SetGoRedisClientWithContext(context.Background(), rdb)
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

	p := path.Join(
		cli.API_PATH,
		cli.API_VERSION,
		cli.EPLatestNews,
	)

	mux := chi.NewRouter()
	mux.Get("/"+p,
		func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()

			if !(r.URL.Query().Get("apikey") == TEST_API_KEY || r.Header.Get("X-ACCESS-KEY") == TEST_API_KEY) {
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

	srvrUrl, err := url.Parse(srvr.URL)
	require.NoError(t, err)
	NItem := 0

	// user's first request
	ckey, cache, err := h.Handle(TEST_API_KEY, TEST_USER_ID, pf)
	require.NoError(t, err)
	require.NotNil(t, cache)
	require.NotEmpty(t, ckey)

	// err = cache.AddRandomSalt(64)
	// require.NoError(t, err)

	ctx := context.Background()
	rhResp, err := rh.JSONSet(ckey, ".", cache)
	require.NoError(t, err)
	require.Equal(t, "OK", rhResp.(string))
	t.Log(ckey)
	rdb.Expire(ctx, ckey, 20*time.Minute)

	// should stop when i == len(tc.Filename)
	for i := 1; i <= len(tc.Filename)+1; i++ {
		req, err := h.RequestFromCacheQuery(cache.Query)
		require.NoError(t, err)

		httpReq, err := req.ToHttpRequest()
		require.NoError(t, err)
		require.NotNil(t, httpReq)
		httpReq.URL.Scheme = srvrUrl.Scheme
		httpReq.URL.Host = srvrUrl.Host

		if i == 1 {
			require.NotContains(t, httpReq.URL.RawQuery, "page=")
		} else {
			require.Contains(t, httpReq.URL.RawQuery, "page=")
		}

		httpResp, err := http.DefaultClient.Do(httpReq)
		require.NoError(t, err)
		require.NotNil(t, httpResp)
		require.Equal(t, http.StatusOK, httpResp.StatusCode)

		resp, err := h.Parse(httpResp)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "success", resp.GetStatus())

		next, prev := resp.ToNewsItemList()
		// update cache
		rh.JSONSet(ckey, ".query.next_page", next)
		cache.SetNextPage(next)

		if len(prev) > 0 {
			NItem += len(prev)
			t.Logf("add %d items\n", len(prev))

			// append items to cache
			for j := range prev {
				rh.JSONArrAppend(ckey, ".news_item", prev[j])
			}
			l, err := rh.JSONArrLen(ckey, ".news_item")
			require.NoError(t, err)
			t.Logf("n items: %d\n", l.(int64))
		}

		if next.Equal(api.StrLastPageToken) {
			break
		}

		// next client request comes in
		// build request from cache
		b, err := rh.JSONGet(ckey, ".query")
		if err != nil {
			require.NoError(t, err)
		}

		var cq api.CacheQuery
		err = json.Unmarshal(b.([]byte), &cq)
		require.NoError(t, err)

		req, err = cli.RequestFromCacheQuery(cq)
		require.NoError(t, err)
	}

	b, err := rh.JSONGet(ckey, ".")
	if err != nil {
		require.NoError(t, err)
	}

	err = json.Unmarshal(b.([]byte), cache)
	require.NoError(t, err)
	require.Equal(t, NItem, cache.Len())
}
