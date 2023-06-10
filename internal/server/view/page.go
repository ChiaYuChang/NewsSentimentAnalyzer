package view

import "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"

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
		AddPair("href", "/static/css/style.css")

	head.Script.
		NewHTMLElement().
		AddPair("src", "/static/js/func.js")

	return head
}
