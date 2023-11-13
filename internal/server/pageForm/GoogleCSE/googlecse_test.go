package googlecse_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	googlecse "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GoogleCSE"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/stretchr/testify/require"
)

func RequestWithForm(method, url string, formVal url.Values) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(formVal.Encode()))

	if method != http.MethodGet && formVal != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, err
}

const TEST_KEYWORD = "[[:TEST_KEYWORD:]]"
const TEST_SEARCH_ENGINE_ID = "[[:TEST_ENG_ID:]]"

func TestGoogleCSEFormValidation(t *testing.T) {
	val := validator.Validate

	decoder := form.NewDecoder()

	r := chi.NewRouter()
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		var cse googlecse.GoogleCSE
		w.Header().Add("Content-Type", "application/json")

		m := map[string]string{}
		if err := r.ParseForm(); err != nil {
			m["status"] = strconv.Itoa(http.StatusBadRequest)
			m["message"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			bs, _ := json.Marshal(m)
			w.Write(bs)
			return
		}

		if err := decoder.Decode(&cse, r.PostForm); err != nil {
			m["status"] = strconv.Itoa(http.StatusBadRequest)
			m["message"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			bs, _ := json.Marshal(m)
			w.Write(bs)
			return
		}

		if err := val.Struct(cse); err != nil {
			m["status"] = strconv.Itoa(http.StatusBadRequest)
			m["message"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			bs, _ := json.Marshal(m)
			w.Write(bs)
			return
		}

		w.WriteHeader(http.StatusOK)
		m["status"] = strconv.Itoa(http.StatusOK)
		m["message"] = "ok"
		m["keyword"] = cse.Keyword
		m["search-engine-id"] = cse.SearchEngineID
		m["date-restrict"] = cse.DateRestrict()
		m["date-restrict-value"] = strconv.Itoa(cse.DateRestrictValue)
		m["date-restrict-unit"] = cse.DateRestrictUnit
		bs, _ := json.Marshal(m)
		w.Write(bs)
		return
	})

	srvc := httptest.NewServer(r)
	type testCase struct {
		Name             string
		Code             int
		NewPageForm      func() url.Values
		TestResponseBody func(t *testing.T, body string)
	}

	tcs := []testCase{
		{
			Name: "OK",
			Code: http.StatusOK,
			NewPageForm: func() url.Values {
				formVal := url.Values{}
				formVal.Add("keyword", TEST_KEYWORD)
				formVal.Add("search-engine-id", TEST_SEARCH_ENGINE_ID)
				formVal.Add("date-restrict-value", "1")
				formVal.Add("date-restrict-unit", "w")
				return formVal
			},
			TestResponseBody: func(t *testing.T, body string) {
				require.Contains(t, body, fmt.Sprintf(`"keyword":"%s"`, TEST_KEYWORD))
				require.Contains(t, body, fmt.Sprintf(`"search-engine-id":"%s"`, TEST_SEARCH_ENGINE_ID))
				require.Contains(t, body, `"date-restrict-unit":"w"`)
				require.Contains(t, body, `"date-restrict-value":"1"`)
				require.Contains(t, body, `"date-restrict":"w1"`)
			},
		},
		{
			Name: "Missing Search Engine Id",
			Code: http.StatusBadRequest,
			NewPageForm: func() url.Values {
				formVal := url.Values{}
				formVal.Add("search-engine-id", TEST_SEARCH_ENGINE_ID)
				formVal.Add("date-restrict-value", "1")
				formVal.Add("date-restrict-unit", "w")
				return formVal
			},
			TestResponseBody: func(t *testing.T, body string) {
				require.Contains(t, body, `'Keyword' failed on the 'required' tag`)
			},
		},
		{
			Name: "Missing Keyword",
			Code: http.StatusBadRequest,
			NewPageForm: func() url.Values {
				formVal := url.Values{}
				formVal.Add("keyword", TEST_KEYWORD)
				formVal.Add("date-restrict-value", "1")
				formVal.Add("date-restrict-unit", "w")
				return formVal
			},
			TestResponseBody: func(t *testing.T, body string) {
				require.Contains(t, body, `'SearchEngineID' failed on the 'required' tag`)
			},
		},
		{
			Name: "Negative date value",
			Code: http.StatusBadRequest,
			NewPageForm: func() url.Values {
				formVal := url.Values{}
				formVal.Add("keyword", TEST_KEYWORD)
				formVal.Add("search-engine-id", TEST_SEARCH_ENGINE_ID)
				formVal.Add("date-restrict-value", "-1")
				formVal.Add("date-restrict-unit", "s")
				return formVal
			},
			TestResponseBody: func(t *testing.T, body string) {
				require.Contains(t, body, `'DateRestrictValue' failed on the 'gte' tag`)
			},
		},
		{
			Name: "Unknown date unit",
			Code: http.StatusBadRequest,
			NewPageForm: func() url.Values {
				formVal := url.Values{}
				formVal.Add("keyword", TEST_KEYWORD)
				formVal.Add("search-engine-id", TEST_SEARCH_ENGINE_ID)
				formVal.Add("date-restrict-value", "1")
				formVal.Add("date-restrict-unit", "s")
				return formVal
			},
			TestResponseBody: func(t *testing.T, body string) {
				require.Contains(t, body, `'DateRestrictUnit' failed on the 'oneof' tag`)
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.Name, func(t *testing.T) {
			req, err := RequestWithForm(
				http.MethodPost,
				fmt.Sprintf("%s/%s", srvc.URL, "test"),
				tc.NewPageForm(),
			)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.Code, resp.StatusCode)
			require.NoError(t, val.Var(string(body), "json"))
			if tc.TestResponseBody != nil {
				tc.TestResponseBody(t, string(body))
			}
		})
	}
}
