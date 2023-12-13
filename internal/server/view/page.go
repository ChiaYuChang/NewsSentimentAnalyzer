package view

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
)

type headContentSingleton struct {
	headContent object.HeadConent
	once        sync.Once
}

var sharedHeadContent headContentSingleton
var jobPageHeadContent headContentSingleton

func SharedHeadContent() object.HeadConent {
	sharedHeadContent.once.Do(func() {
		sharedHeadContent.headContent = NewHeadContent()
	})
	return sharedHeadContent.headContent
}

func JobPageHeadContent() object.HeadConent {
	jobPageHeadContent.once.Do(func() {
		jobPageHeadContent.headContent = SharedHeadContent().Copy()

		jobPageHeadContent.headContent.
			Script.NewHTMLElement().
			AddPair("src", "/static/js/job_funcs.js")

		jobPageHeadContent.headContent.
			Script.NewHTMLElement().
			AddPair("src", "//cdnjs.cloudflare.com/ajax/libs/list.js/2.3.1/list.min.js")

		jobPageHeadContent.headContent.
			Link.NewHTMLElement().
			AddPair("rel", "stylesheet").
			AddPair("href", "/static/css/animation.css")

		jobPageHeadContent.headContent.
			Link.NewHTMLElement().
			AddPair("rel", "stylesheet").
			AddPair("href", "/static/css/jobs.css")
	})
	return jobPageHeadContent.headContent
}

func NewHeadContent() object.HeadConent {
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
		AddPair("rel", "stylesheet").
		AddPair("href", "https://fonts.googleapis.com/css2?family=Roboto+Mono:wght@400;600&display=swap")

	head.Link.
		NewHTMLElement().
		AddPair("rel", "stylesheet").
		AddPair("href", "/static/css/style.css")

	// pure css
	head.Link.
		NewHTMLElement().
		AddPair("rel", "stylesheet").
		AddPair("href", "https://cdn.jsdelivr.net/npm/purecss@3.0.0/build/pure-min.css").
		AddPair("integrity", "sha384-X38yfunGUhNzHpBaEBsWLO+A0HDYOQi8ufWDkZ0k9e0eXz/tH3II7uKZ9msv++Ls").
		AddPair("crossorigin", "anonymous")

	head.Link.
		NewHTMLElement().
		AddPair("rel", "stylesheet").
		AddPair("href", "/static/css/fontawesome.css")

	head.Link.
		NewHTMLElement().
		AddPair("rel", "stylesheet").
		AddPair("href", "/static/css/brands.css")

	head.Link.
		NewHTMLElement().
		AddPair("rel", "stylesheet").
		AddPair("href", "/static/css/solid.css")

	head.Script.
		NewHTMLElement().
		AddPair("src", "/static/js/func.js")

	return head
}

func NewEndPointOptSelector(tmpl *template.Template, opts []object.SelectOpts) ([]byte, error) {
	script := bytes.NewBufferString("")

	if err := tmpl.Execute(script, opts); err != nil {
		return nil, err
	}

	minifiedScript := bytes.NewBufferString("")
	if err := global.Minifier().Minify("application/javascript", minifiedScript, script); err != nil {
		return nil, err
	}
	return minifiedScript.Bytes(), nil
}

var GnewsSelectOpts = []object.SelectOpts{}
