package view_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	htemplate "html/template"
	"strings"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/stretchr/testify/require"
)

const VIEWS_PATH = "../../../views"

func TestHeadConent(t *testing.T) {
	tmpl, err := htemplate.
		New("head").
		ParseFiles(VIEWS_PATH + "/template/head.gotmpl")
	require.NoError(t, err)

	head0 := object.HeadConent{
		Meta:   object.NewHTMLElementList("meta"),
		Link:   object.NewHTMLElementList("link"),
		Script: object.NewHTMLElementList("script"),
	}

	head0.Meta.
		NewHTMLElement().
		AddPair("charset", "UTF-8")

	head0.Meta.
		NewHTMLElement().
		AddPair("http-equiv", "X-UA-Compatible").
		AddPair("content", "IE=edge")

	head0.Meta.
		NewHTMLElement().
		AddPair("name", "viewport").
		AddPair("content", "width=device-width, initial-scale=1.0")

	head0.Link.
		NewHTMLElement().
		AddPair("rel", "preconnect").
		AddPair("href", "https://fonts.googleapis.com")

	head0.Link.
		NewHTMLElement().
		AddPair("rel", "preconnect").
		AddPair("href", "https://fonts.gstatic.com").
		AddVal("crossorigin")

	head0.Link.
		NewHTMLElement().
		AddPair("rel", "stylesheet").
		AddPair("href", "https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700&display=swap")

	head0.Link.
		NewHTMLElement().
		AddPair("ref", "stylesheet").
		AddPair("href", "css/style.css")

	head0.Script.
		NewHTMLElement().
		AddPair("src", "js/func.js")

	t.Run(
		"Use content",
		func(t *testing.T) {
			err = head0.Execute(tmpl.Lookup("head"))
			require.NoError(t, err)
			require.True(t, head0.HasExec())
			require.Contains(t, head0.Content(), `charset="UTF-8"`)
			require.Contains(t, head0.Content(), `ref="stylesheet" href="css/style.css">`)
			require.Contains(t, head0.Content(), `src="js/func.js"`)
		},
	)

	t.Run(
		"Write to buffer",
		func(t *testing.T) {
			bf := bytes.NewBufferString("")
			err = tmpl.ExecuteTemplate(bf, "head", head0)
			require.NoError(t, err)

			require.Contains(t, bf.String(), `charset="UTF-8"`)
			require.Contains(t, bf.String(), `ref="stylesheet" href="css/style.css">`)
			require.Contains(t, bf.String(), `src="js/func.js"`)
		},
	)

	t.Run(
		"From json",
		func(t *testing.T) {
			data, err := json.MarshalIndent(head0, "", "    ")
			require.NoError(t, err)

			var head1 object.HeadConent
			err = head1.FromJson(data, tmpl.Lookup("head"))
			require.NoError(t, err)
			require.Equal(t, head0.Content(), head1.Content())
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

// func TestMinifyJS(t *testing.T) {
// 	tmpl, err := ttemplate.
// 		New("selector.gotmpl").
// 		ParseFiles(VIEWS_PATH + "/template/js/selector.gotmpl")
// 	require.NoError(t, err)
// 	require.NotNil(t, tmpl)

// 	for _, opts := range [][]object.SelectOpts{
// 		view.NEWSDATASelectOpts,
// 		view.GnewsSelectOpts,
// 		view.NewsAPISelectOpts,
// 	} {

// 		bf := bytes.NewBufferString("")
// 		err = tmpl.ExecuteTemplate(bf, "selector.gotmpl", opts)
// 		require.NoError(t, err)

// 		m := minify.New()
// 		m.AddFunc("application/javascript", js.Minify)

// 		sb := &strings.Builder{}
// 		err = m.Minify("application/javascript", sb, bf)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, sb.String())

// 		t.Log(sb)
// 	}
// }

var tmplText string = `
document.addEventListener("DOMContentLoaded", function () {
        {{range .}}
        let {{.PositionId}}Opts = [
            {{$optMap := .OptMap}}
            { value: "{{.DefaultValue}}", txt: "{{.DefaultText}}"},
            {{range $key := .SortedOptKey}}
            { value: "{{$key}}", txt: "{{index $optMap $key}}" },{{end}}
        ];
        addListenerToBtn("{{.PositionId}}", "{{.InsertButtonId}}", "{{.DeleteButtonId}}", {{.MaxDiv}}, {{.PositionId}}Opts, "{{.AlertMessage}}");
        {{end}}
    })
`

// func TestNewEndPointOptSelector(t *testing.T) {
// 	tmpl, err := htemplate.New("selector").Parse(tmplText)
// 	require.NoError(t, err)
// 	require.NotNil(t, tmpl)

// 	type testCase struct {
// 		name string
// 		opts []object.SelectOpts
// 	}

// 	tcs := []testCase{
// 		{
// 			name: "NEWSDATA",
// 			opts: view.NEWSDATASelectOpts,
// 		},
// 		{
// 			name: "Gnews",
// 			opts: view.GnewsSelectOpts,
// 		},
// 		{
// 			name: "NewsAPI",
// 			opts: view.NewsAPISelectOpts,
// 		},
// 	}

// 	for i := range tcs {
// 		tc := tcs[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			sb := &strings.Builder{}
// 			tmpl.Execute(sb, tc.opts)
// 			t.Log(sb.String())
// 		})
// 	}
// }
