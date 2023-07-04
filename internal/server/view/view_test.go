package view_test

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"strings"
	"testing"
	ttemplate "text/template"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

const VIEWS_PATH = "../../../views"

func TestHeadConent(t *testing.T) {
	tmpl, err := htemplate.
		New("head").
		ParseFiles(VIEWS_PATH + "/template/head.gotmpl")
	require.NoError(t, err)

	head := object.HeadConent{
		Meta:   object.NewHTMLElementList("meta"),
		Link:   object.NewHTMLElementList("link"),
		Script: object.NewHTMLElementList("script"),
	}

	head.Meta.
		NewHTMLElement().
		AddPair("charset", "UTF-8")

	head.Meta.
		NewHTMLElement().
		AddPair("http-equiv", "X-UA-Compatible").
		AddPair("content", "IE=edge")

	head.Meta.
		NewHTMLElement().
		AddPair("name", "viewport").
		AddPair("content", "width=device-width, initial-scale=1.0")

	head.Link.
		NewHTMLElement().
		AddPair("rel", "preconnect").
		AddPair("href", "https://fonts.googleapis.com")

	head.Link.
		NewHTMLElement().
		AddPair("rel", "preconnect").
		AddPair("href", "https://fonts.gstatic.com").
		AddVal("crossorigin")

	head.Link.
		NewHTMLElement().
		AddPair("rel", "stylesheet").
		AddPair("href", "https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700&display=swap")

	head.Link.
		NewHTMLElement().
		AddPair("ref", "stylesheet").
		AddPair("href", "css/style.css")

	head.Script.
		NewHTMLElement().
		AddPair("src", "js/func.js")

	t.Run(
		"Use content",
		func(t *testing.T) {
			err = head.Execute(tmpl.Lookup("head"))
			require.NoError(t, err)
			require.True(t, head.HasExec())
			require.Contains(t, head.Content(), `charset="UTF-8"`)
			require.Contains(t, head.Content(), `ref="stylesheet" href="css/style.css">`)
			require.Contains(t, head.Content(), `src="js/func.js"`)
		},
	)

	t.Run(
		"Write to buffer",
		func(t *testing.T) {
			bf := bytes.NewBufferString("")
			err = tmpl.ExecuteTemplate(bf, "head", head)
			require.NoError(t, err)

			require.Contains(t, bf.String(), `charset="UTF-8"`)
			require.Contains(t, bf.String(), `ref="stylesheet" href="css/style.css">`)
			require.Contains(t, bf.String(), `src="js/func.js"`)
		},
	)
}

func TestPageWelcome(t *testing.T) {
	tmpl, err := htemplate.
		ParseFiles(
			VIEWS_PATH+"/template/head.gotmpl",
			VIEWS_PATH+"/template/welcome.gotmpl",
		)
	require.NoError(t, err)

	page := object.WelcomePage{
		Name: "User001",
		Page: object.Page{
			Title: "NewsSentimentanaylzer-Welcome",
			HeadConent: object.HeadConent{
				Meta:   object.NewHTMLElementList("meta"),
				Link:   object.NewHTMLElementList("link"),
				Script: object.NewHTMLElementList("script"),
			},
		},
	}
	page.HeadConent.Meta.
		NewHTMLElement().
		AddPair("charset", "UTF-8")

	page.HeadConent.Link.
		NewHTMLElement().
		AddPair("ref", "stylesheet").
		AddPair("href", "css/style.css")

	page.HeadConent.Script.
		NewHTMLElement().
		AddPair("src", "js/func.js")

	var sb strings.Builder
	var doc string
	page.Role = "admin"
	sb = strings.Builder{}
	err = tmpl.ExecuteTemplate(&sb, "welcome.gotmpl", page)
	require.NoError(t, err)
	doc = sb.String()
	require.Contains(t, doc, "Admin")

	page.Role = "user"
	sb = strings.Builder{}
	err = tmpl.ExecuteTemplate(&sb, "welcome.gotmpl", page)
	doc = sb.String()
	require.NoError(t, err)
	require.NotContains(t, doc, "Admin")
}

func TestPageSignUp(t *testing.T) {
	tmpl, err := htemplate.
		ParseFiles(
			VIEWS_PATH+"/template/head.gotmpl",
			VIEWS_PATH+"/template/signup.gotmpl",
		)
	require.NoError(t, err)

	page := object.SignUpPage{
		Page: object.Page{
			Title: "NewsSentimentanaylzer-Sign up",
			HeadConent: object.HeadConent{
				Meta:   object.NewHTMLElementList("meta"),
				Link:   object.NewHTMLElementList("link"),
				Script: object.NewHTMLElementList("script"),
			},
		},
	}
	page.HeadConent.Meta.
		NewHTMLElement().
		AddPair("charset", "UTF-8")

	page.HeadConent.Link.
		NewHTMLElement().
		AddPair("ref", "stylesheet").
		AddPair("href", "css/style.css")

	page.HeadConent.Script.
		NewHTMLElement().
		AddPair("src", "js/func.js")

	var sb strings.Builder
	var doc string
	page.ShowUsernameHasUsedAlert = false
	sb = strings.Builder{}
	err = tmpl.ExecuteTemplate(&sb, "signup.gotmpl", page)
	require.NoError(t, err)
	doc = sb.String()
	require.NotContains(t, doc, "<span>Username already used</span>")

	page.ShowUsernameHasUsedAlert = true
	sb = strings.Builder{}
	err = tmpl.ExecuteTemplate(&sb, "signup.gotmpl", page)
	require.NoError(t, err)
	doc = sb.String()
	require.Contains(t, doc, "<span>Username already used</span>")
}

func TestPageLogin(t *testing.T) {
	tmpl, err := htemplate.
		ParseFiles(
			VIEWS_PATH+"/template/head.gotmpl",
			VIEWS_PATH+"/template/login.gotmpl",
		)
	require.NoError(t, err)

	page := object.LoginPage{
		Page: object.Page{
			Title: "NewsSentimentanaylzer-Log in",
			HeadConent: object.HeadConent{
				Meta:   object.NewHTMLElementList("meta"),
				Link:   object.NewHTMLElementList("link"),
				Script: object.NewHTMLElementList("script"),
			},
		},
	}
	page.HeadConent.Meta.
		NewHTMLElement().
		AddPair("charset", "UTF-8")

	page.HeadConent.Link.
		NewHTMLElement().
		AddPair("ref", "stylesheet").
		AddPair("href", "css/style.css")

	page.HeadConent.Script.
		NewHTMLElement().
		AddPair("src", "js/func.js")

	var sb strings.Builder
	var doc string

	type testCast struct {
		ShowUsernameNotFountAlert bool
		ShowPasswordMismatchAlert bool
	}

	tcs := []testCast{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("Case %d", i+1),
			func(t *testing.T) {
				page.ShowUsernameNotFountAlert = tc.ShowUsernameNotFountAlert
				page.ShowPasswordMismatchAlert = tc.ShowPasswordMismatchAlert
				sb = strings.Builder{}
				err = tmpl.ExecuteTemplate(&sb, "login.gotmpl", page)
				require.NoError(t, err)
				doc = sb.String()

				if tc.ShowUsernameNotFountAlert {
					require.Contains(t, doc, "<span>Couldn’t find your Account</span>")
				} else {
					require.NotContains(t, doc, "<span>Couldn’t find your Account</span>")
				}

				if tc.ShowPasswordMismatchAlert {
					require.Contains(t, doc, "<span>Wrong password. Please try again.</span>")
				} else {
					require.NotContains(t, doc, "<span>Wrong password. Please try again.</span>")
				}
			},
		)
	}
}

func TestMinifyJS(t *testing.T) {
	tmpl, err := ttemplate.
		New("selector.gotmpl").
		ParseFiles(VIEWS_PATH + "/template/js/selector.gotmpl")
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	for _, opts := range [][]object.SelectOpts{
		view.NEWSDATASelectOpts,
		view.GnewsSelectOpts,
		view.NewsAPISelectOpts,
	} {

		bf := bytes.NewBufferString("")
		err = tmpl.ExecuteTemplate(bf, "selector.gotmpl", opts)
		require.NoError(t, err)

		m := minify.New()
		m.AddFunc("application/javascript", js.Minify)

		sb := &strings.Builder{}
		err = m.Minify("application/javascript", sb, bf)
		require.NoError(t, err)
		require.NotEmpty(t, sb.String())
	}
}

// func TestAPIEndpointsPage(t *testing.T) {
// 	tmpl, err := template.
// 		ParseFiles(
// 			VIEWS_PATH+"/template/head.gotmpl",
// 			VIEWS_PATH+"/template/endpoint.gotmpl",
// 		)
// 	require.NoError(t, err)

// 	page := object.APIEndpointPage{
// 		Page: object.Page{
// 			Title: "NewsSentimentanaylzer-Log in",
// 			HeadConent: object.HeadConent{
// 				Meta:   object.NewHTMLElementList("meta"),
// 				Link:   object.NewHTMLElementList("link"),
// 				Script: object.NewHTMLElementList("script"),
// 			},
// 		},
// 		Endpoints: []object.APIEndpoint{},
// 	}
// 	page.HeadConent.Meta.
// 		NewHTMLElement().
// 		AddPair("charset", "UTF-8")

// 	page.HeadConent.Link.
// 		NewHTMLElement().
// 		AddPair("ref", "stylesheet").
// 		AddPair("href", "css/style.css")

// 	page.HeadConent.Script.
// 		NewHTMLElement().
// 		AddPair("src", "js/func.js")

// 	var ep *object.HTMLElementList
// 	ep = object.NewHTMLElementList("btn")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/NewsAPI/everything.html'").
// 		ToOpeningElement("Everything")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/NewsAPI/top-headlines.html'").
// 		ToOpeningElement("Top Headlines")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/NewsAPI/sources.html'").
// 		ToOpeningElement("Sources")
// 	page.Endpoints = append(
// 		page.Endpoints, object.APIEndpoint{
// 			Image:       object.NewHTMLElementList(""),
// 			DocumentURL: "https://newsapi.org/docs/endpoints",
// 			Endpoints:   ep,
// 		})
// 	page.Endpoints.

// 	ep = object.NewHTMLElementList("btn")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/NEWSDATA/lastest.html'").
// 		ToOpeningElement("Latest News")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/NEWSDATA/archive.html'").
// 		ToOpeningElement("News Archive")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/NEWSDATA/sources.html'").
// 		ToOpeningElement("News Sources")
// 	page.APIEndpoints = append(
// 		page.APIEndpoints, object.APIEndpoint{
// 			Image: *object.NewHTMLElement("").
// 				AddPair("src", "image/NEWSDATA.IO_logo.png").
// 				AddPair("alt", "NEWSDATA.IO"),
// 			DocumentURL: "https://newsdata.io/documentation/#first-api-request",
// 			Endpoints:   ep,
// 		})

// 	ep = object.NewHTMLElementList("btn")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/GNews/search.html'").
// 		ToOpeningElement("Search")
// 	ep.NewHTMLElement().
// 		AddPair("onclick", "location.href='endpoints/GNews/headlines.html'").
// 		ToOpeningElement("Top headlines")
// 	page.APIEndpoints = append(
// 		page.APIEndpoints, object.APIEndpoint{
// 			Image: *object.NewHTMLElement("").
// 				AddPair("src", "image/GNews_Logo.png").
// 				AddPair("alt", "GNews"),
// 			DocumentURL: "https://gnews.io/docs/v4",
// 			Endpoints:   ep,
// 		})

// 	sb := strings.Builder{}
// 	err = tmpl.ExecuteTemplate(&sb, "endpoint.gotmpl", page)
// 	require.NoError(t, err)
// 	doc := sb.String()
// 	t.Log(doc)
// }
