package pageform_test

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

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
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
		pageform.GnewsCatVal.Tag(),
		pageform.GnewsCatVal.ValFun())

	val.RegisterValidation(
		pageform.GNewsCtryVal.Tag(),
		pageform.GNewsCtryVal.ValFun())

	val.RegisterValidation(
		pageform.GNewsLangVal.Tag(),
		pageform.GNewsLangVal.ValFun())

	decoder := form.NewDecoder()
	decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse(time.DateOnly, vals[0])
	}, time.Time{})

	r := chi.NewRouter()
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		var headlines pageform.GNewsHeadlines
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
		m["category"] = headlines.Category
		m["country"] = headlines.Country
		m["language"] = headlines.Language
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
			formVal.Add("language", "en")
			formVal.Add("country", "tw")
			formVal.Add("category", "world")
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
			require.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.NoError(t, val.Var(string(body), "json"))
			require.Contains(t, string(body), tTime.Format(time.DateOnly))
			require.Contains(t, string(body), fTime.Format(time.DateOnly))
		},
	)

	t.Run(
		"To Time Error",
		func(t *testing.T) {
			tTime := time.Now().Add(2 * 24 * time.Hour)
			fTime := time.Now()

			formVal := url.Values{}
			formVal.Add("keyword", "")
			formVal.Add("language", "en")
			formVal.Add("country", "tw")
			formVal.Add("category", "world")
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
			require.Contains(t, string(body), "'To' failed on the 'lte' tag")
		},
	)

	t.Run(
		"Error Language",
		func(t *testing.T) {
			tTime := time.Now()
			fTime := time.Now().Add(-10 * time.Hour)

			formVal := url.Values{}
			formVal.Add("keyword", "")
			formVal.Add("language", "xx")
			formVal.Add("country", "tw")
			formVal.Add("category", "world")
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
			formVal.Add("language", "en")
			formVal.Add("country", "xx")
			formVal.Add("category", "world")
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
		pageform.GnewsCatVal,
		pageform.GNewsCtryVal,
		pageform.GNewsLangVal,
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

	require.NoError(t, val.Var("business", pageform.GnewsCatVal.Tag()))
	require.NoError(t, val.Struct(valCatStruct{Category: "business"}))
	require.Error(t, val.Var("xx", pageform.GnewsCatVal.Tag()))
	require.Error(t, val.Struct(valCatStruct{Category: "xx"}))

	require.NoError(t, val.Var("tw", pageform.GNewsCtryVal.Tag()))
	require.NoError(t, val.Struct(valCtryStruct{Country: "tw"}))
	require.Error(t, val.Var("xx", pageform.GNewsCtryVal.Tag()))
	require.Error(t, val.Struct(valCtryStruct{Country: "xx"}))

	require.NoError(t, val.Var("en", pageform.GNewsLangVal.Tag()))
	require.NoError(t, val.Struct(valLangStruct{Language: "en"}))
	require.Error(t, val.Var("xx", pageform.GNewsLangVal.Tag()))
	require.Error(t, val.Struct(valLangStruct{Language: "xx"}))
}
