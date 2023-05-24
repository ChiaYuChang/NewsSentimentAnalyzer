package newsapi_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	newsapi "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/newsAPI"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/code"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/stretchr/testify/require"
)

const (
	API_KEY_OK        string = "[[:OK]]"
	API_KEY_DISABLED         = "[[:DISABLE:]]"
	API_KEY_INVALID          = "[[:INVALID:]]"
	API_KEY_EXHAUSTED        = "[[:EXHAUSTED:]]"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var apikey string
		apikey = req.Header.Get("X-Api-Key")
		if apikey == "" {
			apikey = req.Header.Get("Authorization")
			apikey = strings.TrimPrefix(apikey, "Bearer")
		}

		apikey = strings.Trim(apikey, " ")
		if apikey != API_KEY_OK {
			errResp := newsapi.APIError{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
			}
			switch apikey {
			case "":
				errResp.Message = "apiKeyMissing"
			case API_KEY_DISABLED:
				errResp.Message = "apiKeyDisable"
			case API_KEY_EXHAUSTED:
				errResp.Message = "apiKeyExhausted"
			case API_KEY_INVALID:
				errResp.Message = "apiKeyInvalid"
			}
			jsonObj, _ := json.Marshal(errResp)
			w.WriteHeader(errResp.Code)
			w.Write(jsonObj)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func TestAPIAuthError(t *testing.T) {
	var err error
	var req *http.Request

	mux := http.NewServeMux()
	mux.HandleFunc(
		fmt.Sprintf("/%s/%s", newsapi.API_VERSION, newsapi.EPEverything),
		func(w http.ResponseWriter, r *http.Request) {
			return
		})

	srv := httptest.NewTLSServer(authMiddleware(mux))
	defer srv.Close()
	newsapi.API_URL = srv.URL

	cli := newsapi.NewClient("", &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}})

	query := cli.NewQueryWithDefaultVals()
	t.Run(
		"API KEY is missing",
		func(t *testing.T) {
			_, err = query.ToHTTPRequest(context.TODO(), cli.ApiKey)
			require.Equal(t, ec.ECUnauthorized, err.(*ec.Error).ErrorCode)
		},
	)

	type testCase struct {
		name   string
		apikey string
		eMsg   string
		eC     ec.ErrorCode
	}

	tcs := []testCase{
		{
			name:   "API KEY is disabled",
			apikey: API_KEY_DISABLED,
			eMsg:   "apiKeyDisable",
			eC:     ec.ECUnauthorized,
		},
		{
			name:   "API KEY is exhausted",
			apikey: API_KEY_EXHAUSTED,
			eMsg:   "apiKeyExhausted",
			eC:     ec.ECUnauthorized,
		},
		{
			name:   "API KEY is invalid",
			apikey: API_KEY_INVALID,
			eMsg:   "apiKeyInvalid",
			eC:     ec.ECUnauthorized,
		},
	}

	for _, tc := range tcs {
		t.Run(
			tc.name,
			func(t *testing.T) {
				cli.SetAPIKey(tc.apikey)
				req, err = query.ToHTTPRequest(context.TODO(), cli.ApiKey)
				require.NoError(t, err)
				require.NotNil(t, req)

				resp, err := cli.Do(req)
				require.NoError(t, err)
				jsonObj, err := newsapi.ParseHTTPResponse(resp)
				require.Error(t, err)

				ecErr, ok := err.(*ec.Error)
				require.True(t, ok)
				require.Equal(t, tc.eC, ecErr.ErrorCode)
				require.Contains(t, ecErr.Details, tc.eMsg)
				require.Nil(t, jsonObj)
			},
		)
	}
}

func TestAPIErrToErr(t *testing.T) {
	type testCase struct {
		statusCode int
		ecError    *ec.Error
		message    string
	}

	tcs := []testCase{
		{200, ec.MustGetErr(ec.Success).(*ec.Error), ""},
		{400, ec.MustGetErr(ec.ECBadRequest).(*ec.Error), ""},
		{499, ec.MustGetErr(ec.ECBadRequest).(*ec.Error), "Unknown error"},
		{401, ec.MustGetErr(ec.ECUnauthorized).(*ec.Error), "api key is missing"},
		{429, ec.MustGetErr(ec.ECTooManyRequests).(*ec.Error), "too many request"},
		{500, ec.MustGetErr(ec.ECServerError).(*ec.Error), ""},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%d", i+1, tc.statusCode),
			func(t *testing.T) {
				apiErr := newsapi.APIError{
					Code:    tc.statusCode,
					Status:  tc.ecError.Message,
					Message: tc.message,
				}

				err, ok := apiErr.ToError().(*ec.Error)
				require.True(t, ok)
				require.True(t, tc.ecError.IsEqual(err))
				if tc.message != "" {
					require.Contains(t, err.Details, tc.message)
				}
			},
		)
	}
}

func TestPager(t *testing.T) {
	var err error
	cli := newsapi.NewDefaultClient()

	p1 := cli.NewPagerWithDefaultVals()
	require.Equal(t, newsapi.API_MAX_PAGE_SIZE, p1.PageSize)
	require.Equal(t, newsapi.API_DEFAULT_PAGE, p1.Page)
	_, err = p1.ToUrlVals(url.Values{})
	require.NoError(t, err)

	p2 := cli.NewPager(newsapi.API_MAX_PAGE_SIZE+1, 2)
	_, err = p2.ToUrlVals(url.Values{})
	require.Error(t, err)
	require.True(t, ec.MustGetErr(ec.ECBadRequest).(*ec.Error).IsEqual(err))
}

func TestNewsSources(t *testing.T) {
	var err error
	cli := newsapi.NewDefaultClient()

	srcs := make([]string, newsapi.API_MAX_SOURCES_NUM)
	for i := range srcs {
		srcs[i] = fmt.Sprintf("source%03d", i+1)
	}
	src1 := cli.NewNewsSources(srcs...)
	val1, err := src1.ToUrlVals(url.Values{})
	require.NoError(t, err)
	require.Equal(t, strings.Join(srcs, ","), val1.Get("sources"))

	srcs = append(srcs, "extra sources")
	src2 := cli.NewNewsSources(srcs...)
	_, err = src2.ToUrlVals(url.Values{})
	require.Error(t, err)
	require.True(t, ec.MustGetErr(ec.ECBadRequest).(*ec.Error).IsEqual(err))
}

type failedParams struct {
	testErr error
}

func (p failedParams) ToUrlVals(vals url.Values) (url.Values, error) {
	return vals, p.testErr
}

func (p failedParams) ParamsName() string {
	return "Must-Failed"
}

func TestParamsToUrlFailed(t *testing.T) {
	testErr := errors.New("test error")
	cli := newsapi.NewDefaultClient()
	query := cli.NewQueryWithDefaultVals()

	query.AppendParams(failedParams{testErr})
	_, err := query.ToHTTPRequest(context.TODO(), API_KEY_OK)
	require.ErrorIs(t, testErr, err)
}

func TestQueryWithParams(t *testing.T) {
	var err error
	var req *http.Request
	tstmp1, _ := time.Parse(newsapi.API_TIME_FORMAT, "2023-05-03T18:38:20Z")
	tstmp2, _ := time.Parse(newsapi.API_TIME_FORMAT, "2023-05-09T05:00:00Z")
	articles := []newsapi.Article{
		{
			ArticleSource: newsapi.ArticleSource{
				Id:   "null",
				Name: "Genbeta.com",
			},
			Author:      "Marcos Merino",
			Title:       "No, Susana Griso y Chicote no te animan a invertir en bitcoin, ni les detuvieron por difundir un \"vacío legal para hacernos ricos\"",
			Description: "Alberto Chicote, el famoso chef y presentador televisivo de 'pesadilla...",
			Url:         "https://www.genbeta.com/seguridad/no-susana-griso-chicote-no-te-animan-a-invertir-bitcoin-les-detuvieron-difundir-vacio-legal-para-hacernos-ricos",
			UrlToImage:  "https://i.blogs.es/291173/chicote_griso/840_560.jpeg",
			PublishedAt: tstmp1,
			Content:     "Alberto Chicote, el famoso chef y presentador televisivo de 'pesadilla en la cocina', ha sido arrestado por la policía tras revelar, sin querer, ante un micrófono abierto en su última entrevista en...",
		},
		{
			ArticleSource: newsapi.ArticleSource{
				Id:   "null",
				Name: "Techbang.com",
			},
			Author:      "Odaily",
			Title:       "「聰」的計價時代已到來！BRC-20讓等了14年的比特幣生態成為可能",
			Description: "如果BRC-20 能幫助比特幣實現生態拓展，那麼很大的共識會回到比特幣本身，這就會對其他項目造成影響。...",
			Url:         "https://www.techbang.com/posts/106060-sat-bitcoin-ecosystem",
			UrlToImage:  "https://cdn0.techbang.com/system/excerpt_images/106060/original/af0c884b841ff8ac72449d5bc44c72e0.jpg?1683606608",
			PublishedAt: tstmp2,
			Content:     "14 Casey Rodarmor  \r\n2022 12 Rodarmor Ordinals Inscriptions NFT2023 1 21 Ordinals 0.4.0 BTC NFT  \r\nBTC NFT\r\nBTC NFT \r\nNFTRodarmorNFTBTC \r\n2100satSatoshi11Rodarmor \r\n 0 2,100,000,000,000,000 NFT  \r\nRo...",
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		fmt.Sprintf("/%s/%s", newsapi.API_VERSION, newsapi.EPEverything),
		func(w http.ResponseWriter, r *http.Request) {
			ecOK := ec.MustGetErr(ec.Success).(*ec.Error)
			w.WriteHeader(ecOK.HttpStatusCode)

			jsonObj := newsapi.Response{
				Status:      ecOK.Message,
				TotalResult: len(articles),
				Articles:    articles,
			}
			b, _ := json.Marshal(jsonObj)
			w.Write(b)
			return
		})

	srv := httptest.NewTLSServer(authMiddleware(mux))
	defer srv.Close()
	newsapi.API_URL = srv.URL
	cli := newsapi.NewClient(API_KEY_OK, &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}})

	query := cli.NewQueryWithDefaultVals()
	query.WithKeywords("bitcoin AND NFT")
	query.AppendParams(cli.NewPager(10, 1))

	params := cli.NewEverythingParamsWithDefaultVals()
	params.SearchIn = []string{
		string(newsapi.InTitle),
		string(newsapi.InDescription)}

	tTo, tFrom := "2023-05-04T00:00:00Z", "2023-05-08T00:00:00Z"
	params.To, err = time.Parse(newsapi.API_TIME_FORMAT, tTo)
	require.NoError(t, err)
	params.From, _ = time.Parse(newsapi.API_TIME_FORMAT, tFrom)
	params.Language = code.LEnglish
	params.SortedBy = newsapi.ByPopularity
	query.AppendParams(params)

	req, err = query.ToHTTPRequest(context.TODO(), cli.ApiKey)
	require.NoError(t, err)
	require.NotNil(t, req)
	require.Equal(t, "", req.URL.Query().Get("page"))
	require.Equal(t, strconv.Itoa(10), req.URL.Query().Get("pageSize"))
	require.Equal(t, fmt.Sprintf("%s,%s", newsapi.InTitle, newsapi.InDescription), req.URL.Query().Get("searchIn"))
	require.Equal(t, string(code.LEnglish), req.URL.Query().Get("language"))
	require.Equal(t, tTo, req.URL.Query().Get("to"))
	require.Equal(t, tFrom, req.URL.Query().Get("from"))
	require.Equal(t, string(newsapi.ByPopularity), req.URL.Query().Get("sortBy"))

	t.Logf("req: %s", req.URL)

	resp, err := cli.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	t.Logf("resp: %s", resp.Request.URL)
	obj, err := newsapi.ParseHTTPResponse(resp)
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Equal(t, len(articles), obj.TotalResult)
	require.Equal(t, "success", obj.Status)
	for i := range articles {
		require.Equal(t, articles[i].Title, obj.Articles[i].Title)
	}
}
