package view_test

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"testing"

	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/NEWSDATA"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/stretchr/testify/require"
)

const VIEWS_PATH = "../../../views"

func TestHeadConent(t *testing.T) {
	tmpl, err := template.
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

	err = tmpl.ExecuteTemplate(os.Stdout, "head", head)
	require.NoError(t, err)
}

func TestPageWelcome(t *testing.T) {
	tmpl, err := template.
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
	require.Contains(t, doc, "href='admin.html'")

	page.Role = "user"
	sb = strings.Builder{}
	err = tmpl.ExecuteTemplate(&sb, "welcome.gotmpl", page)
	doc = sb.String()
	t.Log(doc)
	require.NoError(t, err)
	require.NotContains(t, doc, "href='admin.html'")
}

func TestPageSignUp(t *testing.T) {
	tmpl, err := template.
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
	tmpl, err := template.
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

func TestJS(t *testing.T) {
	tmpl, err := template.
		New("head").
		ParseFiles(VIEWS_PATH + "/template/js/selector.gotmpl")
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	opts := []object.SelectOpts{
		{
			OptMap:         newsdata.CatList,
			MaxDiv:         5,
			DefaultValue:   "",
			DefaultText:    "All",
			InsertButtonId: "iCatBtn",
			DeleteButtonId: "dCatBtn",
			PositionId:     "category",
			AlertMessage:   "haha",
		},
		{
			OptMap:         newsdata.CtryList,
			MaxDiv:         5,
			DefaultValue:   "",
			DefaultText:    "All",
			InsertButtonId: "iCtryBtn",
			DeleteButtonId: "dCtryBtn",
			PositionId:     "country",
			AlertMessage:   "haha",
		},
		{
			OptMap:         newsdata.LangList,
			MaxDiv:         5,
			DefaultValue:   "",
			DefaultText:    "All",
			InsertButtonId: "iLangBtn",
			DeleteButtonId: "dLangBtn",
			PositionId:     "language",
			AlertMessage:   "haha",
		},
	}

	sb := &strings.Builder{}
	tmpl.ExecuteTemplate(sb, "selector.gotmpl", opts)
	t.Log(sb.String())
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
