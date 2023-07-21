package newsdata_test

import (
	"encoding/json"
	"os"
	"testing"

	cli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/NEWSDATA"
	"github.com/stretchr/testify/require"
)

const TEST_API_KEY = "pub_00000x0000x00xx0x0000xxxx00x00xx0xx002"

func TestParseResponse(t *testing.T) {
	type testCase struct {
		name        string
		fileName    string
		hasNextPage bool
	}

	tcs := []testCase{
		{
			name:        "w/o next page",
			fileName:    "example_response/001.json",
			hasNextPage: false,
		},
		{
			name:        "w/ next page",
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

				var resp cli.Response
				err = json.Unmarshal(respJson, &resp)
				require.NoError(t, err)
				require.Equal(t, tc.hasNextPage, resp.HasNextPage())
				require.GreaterOrEqual(t, resp.TotalResult, resp.Len())
				if resp.HasNextPage() {
					t.Log(resp.NextPage)
				}
			},
		)
	}
}

// func TestBuildLatestNewsQuery(t *testing.T) {
// 	var err error
// 	v := val.New()
// 	err = v.RegisterValidation(
// 		newsdata.CatVal.Tag(),
// 		newsdata.CatVal.ValFun())
// 	require.NoError(t, err)

// 	err = v.RegisterValidation(
// 		newsdata.LangVal.Tag(),
// 		newsdata.LangVal.ValFun())
// 	require.NoError(t, err)

// 	err = v.RegisterValidation(
// 		newsdata.CtryVal.Tag(),
// 		newsdata.CtryVal.ValFun())
// 	require.NoError(t, err)

// 	v.RegisterAlias(newsdata.VAL_TAG_DOMAIN, "max=250")

// 	type testCase struct {
// 		Name      string
// 		Form      *pageform.NEWSDATAIOLatestNews
// 		CheckFunc func(form *pageform.NEWSDATAIOLatestNews)
// 	}

// 	tcs := []testCase{
// 		{
// 			Name: "OK",
// 			Form: &pageform.NEWSDATAIOLatestNews{
// 				Keyword:  "Golang AND Google",
// 				Language: []string{"zh", "en"},
// 				Country:  []string{"tw", "us", "gb"},
// 				Category: []string{"technology"},
// 			},
// 			CheckFunc: func(form *pageform.NEWSDATAIOLatestNews) {
// 				t.Log(form)

// 				b := newsdata.NewQueryBuilder(APIKey, v)
// 				q, err := b.BuildLatestNewsQuery(APIKey, form)
// 				require.NoError(t, err)
// 				require.NotNil(t, q)
// 				qURL := q.ToRequestURL(newsdata.API_URL)
// 				require.Contains(t, qURL, "Golang+AND+Google")
// 				require.Contains(t, qURL, "category=technology")
// 				require.Contains(t, qURL, "country=tw%2Cus%2Cgb")
// 				require.Contains(t, qURL, "language=zh%2Cen")
// 				require.NotContains(t, qURL, "domain=")
// 			},
// 		},
// 		{
// 			Name: "Error Country",
// 			Form: &pageform.NEWSDATAIOLatestNews{
// 				Keyword:  "Golang AND Google",
// 				Language: []string{"zh", "en"},
// 				Country:  []string{"tw", "us", "xx"},
// 				Category: []string{"technology"},
// 			},
// 			CheckFunc: func(form *pageform.NEWSDATAIOLatestNews) {
// 				b := newsdata.NewQueryBuilder(APIKey, v)
// 				q, err := b.BuildLatestNewsQuery(APIKey, form)
// 				require.Error(t, err)
// 				require.NotNil(t, q)
// 				require.ErrorContains(t, err, fmt.Sprintf("failed on the '%s' tag", newsdata.VAL_TAG_COUNTRY))
// 			},
// 		},
// 		{
// 			Name: "Error Category",
// 			Form: &pageform.NEWSDATAIOLatestNews{
// 				Keyword:  "Golang AND Google",
// 				Language: []string{"zh", "en"},
// 				Country:  []string{"tw", "us", "ph"},
// 				Category: []string{"technology", "nothiscategory"},
// 			},
// 			CheckFunc: func(form *pageform.NEWSDATAIOLatestNews) {
// 				b := newsdata.NewQueryBuilder(APIKey, v)
// 				q, err := b.BuildLatestNewsQuery(APIKey, form)
// 				require.Error(t, err)
// 				require.NotNil(t, q)
// 				require.ErrorContains(t, err, fmt.Sprintf("failed on the '%s' tag", newsdata.VAL_TAG_CATEGORY))
// 			},
// 		},
// 		{
// 			Name: "Error Language",
// 			Form: &pageform.NEWSDATAIOLatestNews{
// 				Keyword:  "Golang AND Google",
// 				Language: []string{"zh", "xx"},
// 				Country:  []string{"tw", "us", "ph"},
// 				Category: []string{"technology"},
// 			},
// 			CheckFunc: func(form *pageform.NEWSDATAIOLatestNews) {
// 				b := newsdata.NewQueryBuilder(APIKey, v)
// 				q, err := b.BuildLatestNewsQuery(APIKey, form)
// 				require.Error(t, err)
// 				require.NotNil(t, q)
// 				require.ErrorContains(t, err, fmt.Sprintf("failed on the '%s' tag", newsdata.VAL_TAG_LANGUAGE))
// 			},
// 		},
// 		{
// 			Name: "Keyword too long",
// 			Form: &pageform.NEWSDATAIOLatestNews{
// 				Keyword:  rg.Must[string](rg.AlphaNum.GenRdmString(513)),
// 				Language: []string{"zh", "en"},
// 				Country:  []string{"tw", "us", "ph"},
// 				Category: []string{"technology"},
// 			},
// 			CheckFunc: func(form *pageform.NEWSDATAIOLatestNews) {
// 				b := newsdata.NewQueryBuilder(APIKey, v)
// 				q, err := b.BuildLatestNewsQuery(APIKey, form)
// 				require.Error(t, err)
// 				require.NotNil(t, q)
// 				require.ErrorContains(t, err, "failed on the 'max' tag")
// 			},
// 		},
// 		{
// 			Name: "Too many Country",
// 			Form: &pageform.NEWSDATAIOLatestNews{
// 				Keyword:  "Golang AND Google",
// 				Language: []string{"zh", "en"},
// 				Country:  []string{"tw", "us", "ph", "mo", "jo", "jp"},
// 				Category: []string{"technology"},
// 			},
// 			CheckFunc: func(form *pageform.NEWSDATAIOLatestNews) {
// 				b := newsdata.NewQueryBuilder(APIKey, v)
// 				q, err := b.BuildLatestNewsQuery(APIKey, form)
// 				require.Error(t, err)
// 				require.NotNil(t, q)
// 				t.Log(err)
// 				require.ErrorContains(t, err, "failed on the 'max' tag")
// 			},
// 		},
// 	}

// 	for i := range tcs {
// 		tc := tcs[i]
// 		t.Run(
// 			tcs[i].Name,
// 			func(t *testing.T) {
// 				tc.CheckFunc(tc.Form)
// 			},
// 		)
// 	}
// }

// func TestBuildNewsArchive(t *testing.T) {
// 	var err error
// 	v := val.New()
// 	err = v.RegisterValidation(
// 		newsdata.CatVal.Tag(),
// 		newsdata.CatVal.ValFun())
// 	require.NoError(t, err)

// 	err = v.RegisterValidation(
// 		newsdata.LangVal.Tag(),
// 		newsdata.LangVal.ValFun())
// 	require.NoError(t, err)

// 	err = v.RegisterValidation(
// 		newsdata.CtryVal.Tag(),
// 		newsdata.CtryVal.ValFun())
// 	require.NoError(t, err)

// 	v.RegisterAlias(newsdata.VAL_TAG_DOMAIN, "max=250")

// 	type testCase struct {
// 		Name      string
// 		Form      *pageform.NEWSDATAIONewsArchive
// 		CheckFunc func(form *pageform.NEWSDATAIONewsArchive)
// 	}

// 	nw := time.Now()
// 	tcs := []testCase{
// 		{
// 			Name: "OK",
// 			Form: &pageform.NEWSDATAIONewsArchive{
// 				TimeRange: pageform.TimeRange{
// 					Form: nw.Add(-24 * time.Hour),
// 					To:   nw.Add(-1 * time.Hour),
// 				},
// 				Keyword:  "Golang AND Google",
// 				Domains:  "nytimes,bbc",
// 				Country:  []string{"tw", "us", "gb"},
// 				Category: []string{"technology"},
// 			},
// 			CheckFunc: func(form *pageform.NEWSDATAIONewsArchive) {
// 				t.Log(form)

// 				b := newsdata.NewQueryBuilder(APIKey, v)
// 				q, err := b.BuildNewsArchive(APIKey, form)
// 				require.NoError(t, err)
// 				require.NotNil(t, q)
// 				qURL := q.ToRequestURL(newsdata.API_URL)
// 				require.Contains(t, qURL, "Golang+AND+Google")
// 				require.Contains(t, qURL, "category=technology")
// 				require.Contains(t, qURL, "country=tw%2Cus%2Cgb")
// 				require.Contains(t, qURL, "domain=nytimes%2Cbbc")
// 				require.NotContains(t, qURL, "language=")
// 				t.Log(qURL)
// 			},
// 		},
// 	}

// 	for i := range tcs {
// 		tc := tcs[i]
// 		t.Run(
// 			tcs[i].Name,
// 			func(t *testing.T) {
// 				tc.CheckFunc(tc.Form)
// 			},
// 		)
// 	}
// }
