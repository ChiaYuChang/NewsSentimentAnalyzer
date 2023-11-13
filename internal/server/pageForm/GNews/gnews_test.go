package gnews_test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	gnews "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GNews"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/stretchr/testify/require"
)

func RequestWithForm(method, url string, formVal url.Values) (*http.Request, error) {
	return http.NewRequest(method, url, strings.NewReader(formVal.Encode()))
}

func TestGnewsFormValidation(t *testing.T) {
	val := validator.Validate
	val.RegisterValidation(
		gnews.LanguageValidator.Tag(),
		gnews.LanguageValidator.ValFun())

	val.RegisterValidation(
		gnews.CountryValidator.Tag(),
		gnews.CountryValidator.ValFun())

	val.RegisterValidation(
		gnews.CategoryValidator.Tag(),
		gnews.CategoryValidator.ValFun())

	decoder := form.NewDecoder()
	decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse(time.DateOnly, vals[0])
	}, time.Time{})

	r := chi.NewRouter()
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		var headlines gnews.GNewsHeadlines
		w.Header().Add("Content-Type", "application/json")

		m := map[string]string{}
		if err := r.ParseForm(); err != nil {
			m["status"] = strconv.Itoa(http.StatusInternalServerError)
			m["message"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			bs, _ := json.Marshal(m)
			w.Write(bs)
			return
		}

		if err := decoder.Decode(&headlines, r.PostForm); err != nil {
			m["status"] = strconv.Itoa(http.StatusInternalServerError)
			m["message"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			bs, _ := json.Marshal(m)
			w.Write(bs)
			return
		}

		if err := val.Struct(headlines); err != nil {
			m["status"] = strconv.Itoa(http.StatusInternalServerError)
			m["message"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			bs, _ := json.Marshal(m)
			w.Write(bs)
			return
		}

		w.WriteHeader(http.StatusOK)
		m["status"] = strconv.Itoa(http.StatusOK)
		m["message"] = "ok"
		m["keyword"] = headlines.Keyword
		m["category"] = strings.Join(headlines.Category, ",")
		m["country"] = strings.Join(headlines.Country, ",")
		m["language"] = strings.Join(headlines.Language, ",")
		m["from"] = headlines.TimeRange.Form.Format(time.DateOnly)
		m["to"] = headlines.TimeRange.To.Format(time.DateOnly)
		bs, _ := json.Marshal(m)
		w.Write(bs)
		return
	})

	srvc := httptest.NewTLSServer(r)
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	t.Run(
		"OK",
		func(t *testing.T) {
			tTime := time.Now()
			fTime := time.Now().Add(-10 * time.Hour)

			formVal := url.Values{}
			formVal.Add("keyword", "")
			formVal.Add("language[0]", "en")
			formVal.Add("country[0]", "tw")
			formVal.Add("category[0]", "world")
			formVal.Add("from-time", fTime.Format(time.DateOnly))
			formVal.Add("to-time", tTime.Format(time.DateOnly))
			formVal.Add("timezone", "UTC")

			req, err := RequestWithForm(
				http.MethodPost,
				fmt.Sprintf("%s/%s", srvc.URL, "test"),
				formVal,
			)

			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := cli.Do(req)
			require.NoError(t, err)
			// require.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, val.Var(string(body), "json"))
			require.Contains(t, string(body), tTime.Format(time.DateOnly))
			require.Contains(t, string(body), fTime.Format(time.DateOnly))
		},
	)

	t.Run(
		"Error Language",
		func(t *testing.T) {
			tTime := time.Now()
			fTime := time.Now().Add(-10 * time.Hour)

			formVal := url.Values{}
			formVal.Add("keyword", "")
			formVal.Add("language[0]", "xx")
			formVal.Add("country[0]", "tw")
			formVal.Add("category[0]", "world")
			formVal.Add("from-time", fTime.Format(time.DateOnly))
			formVal.Add("to-time", tTime.Format(time.DateOnly))

			req, err := RequestWithForm(
				http.MethodPost,
				fmt.Sprintf("%s/%s", srvc.URL, "test"),
				formVal,
			)

			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := cli.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Contains(t, string(body), "'Language' failed on the 'gnews_lang' tag")
		},
	)

	t.Run(
		"Error Country",
		func(t *testing.T) {
			tTime := time.Now()
			fTime := time.Now().Add(-10 * time.Hour)

			formVal := url.Values{}
			formVal.Add("keyword", "")
			formVal.Add("language[0]", "en")
			formVal.Add("country[0]", "xx")
			formVal.Add("category[0]", "world")
			formVal.Add("from-time", fTime.Format(time.DateOnly))
			formVal.Add("to-time", tTime.Format(time.DateOnly))

			req, err := RequestWithForm(
				http.MethodPost,
				fmt.Sprintf("%s/%s", srvc.URL, "test"),
				formVal,
			)

			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := cli.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Contains(t, string(body), "'Country' failed on the 'gnews_ctry' tag")
		},
	)
}

func TestGnewsFormValidationStruct(t *testing.T) {
	val := validator.Validate
	err := validator.RegisterValidator(
		val,
		gnews.LanguageValidator,
		gnews.CountryValidator,
		gnews.CategoryValidator,
	)
	require.NoError(t, err)

	type valCatStruct struct {
		Category string `validate:"gnews_cat"`
	}

	type valCtryStruct struct {
		Country string `validate:"gnews_ctry"`
	}

	type valLangStruct struct {
		Language string `validate:"gnews_lang"`
	}

	require.NoError(t, val.Var("business", gnews.CategoryValidator.Tag()))
	require.NoError(t, val.Struct(valCatStruct{Category: "business"}))
	require.Error(t, val.Var("xx", gnews.CategoryValidator.Tag()))
	require.Error(t, val.Struct(valCatStruct{Category: "xx"}))

	require.NoError(t, val.Var("tw", gnews.CountryValidator.Tag()))
	require.NoError(t, val.Struct(valCtryStruct{Country: "tw"}))
	require.Error(t, val.Var("xx", gnews.CountryValidator.Tag()))
	require.Error(t, val.Struct(valCtryStruct{Country: "xx"}))

	require.NoError(t, val.Var("en", gnews.LanguageValidator.Tag()))
	require.NoError(t, val.Struct(valLangStruct{Language: "en"}))
	require.Error(t, val.Var("xx", gnews.LanguageValidator.Tag()))
	require.Error(t, val.Struct(valLangStruct{Language: "xx"}))
}
