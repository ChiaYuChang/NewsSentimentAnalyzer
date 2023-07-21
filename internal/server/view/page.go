package view

import (
	"net/http"

	gnews "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/GNews"
	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/newsapi"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
)

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

var SharedHeadContent = NewHeadContent()

var ErrorPage500 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent, Title: "500 error"},
	ErrorCode:          500,
	ErrorMessage:       "Sorry, unexpected error",
	ErrorDetail:        "The server encountered an internal error or misconfiguration and was unable to complete your request.",
	ShouldAutoRedirect: false,
}
var ErrorPage400 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent, Title: "400 error"},
	ErrorCode:          http.StatusBadRequest,
	ErrorMessage:       "Bad request",
	ErrorDetail:        "There was a problem with your request.",
	ShouldAutoRedirect: false,
}

var ErrorPage401 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent, Title: "401 error"},
	ErrorCode:          http.StatusUnauthorized,
	ErrorMessage:       "Unauthorized",
	ErrorDetail:        "You are not authorized to access this page.",
	ShouldAutoRedirect: true,
	RedirectPageUrl:    "/login",
	RedirectPageName:   "log in page",
	CountDownFrom:      5,
}

var ErrorPage403 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent, Title: "403 error"},
	ErrorCode:          http.StatusForbidden,
	ErrorMessage:       "Access denied",
	ErrorDetail:        "You do not have premission to access this page.",
	ShouldAutoRedirect: false,
}

var ErrorPage404 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent, Title: "404 error"},
	ErrorCode:          http.StatusNotFound,
	ErrorMessage:       "Page not found",
	ErrorDetail:        "The page you are looking for may have been moved, deleted, or possibly never existed.",
	ShouldAutoRedirect: false,
}

var ErrorPage429 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent, Title: "409 error"},
	ErrorCode:          http.StatusTooManyRequests,
	ErrorMessage:       "Too Many Requests",
	ErrorDetail:        "You have sent too many requests to us recently. Please try again later.",
	ShouldAutoRedirect: false,
}

var NEWSDATASelectOpts = []object.SelectOpts{
	{
		OptMap:         newsdata.Country,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-category-btn",
		DeleteButtonId: "delete-category-btn",
		PositionId:     "category",
		AlertMessage:   "You can only add up to 5 categories in a single query",
	},
	{
		OptMap:         newsdata.Category,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-country-btn",
		DeleteButtonId: "delete-country-btn",
		PositionId:     "country",
		AlertMessage:   "You can only add up to 5 countries in a single query",
	},
	{
		OptMap:         newsdata.Language,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-lang-btn",
		DeleteButtonId: "delete-lang-btn",
		PositionId:     "language",
		AlertMessage:   "You can only add up to 5 languages in a single query",
	},
}

var GnewsSelectOpts = []object.SelectOpts{
	{
		OptMap:         gnews.Category,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-category-btn",
		DeleteButtonId: "delete-category-btn",
		PositionId:     "category",
		AlertMessage:   "You can only add up to 5 categories in a single query",
	},
	{
		OptMap:         gnews.Country,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-country-btn",
		DeleteButtonId: "delete-country-btn",
		PositionId:     "country",
		AlertMessage:   "You can only add up to 5 countries in a single query",
	},
	{
		OptMap:         gnews.Language,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-lang-btn",
		DeleteButtonId: "delete-lang-btn",
		PositionId:     "language",
		AlertMessage:   "You can only add up to 5 languages in a single query",
	},
}

var NewsAPISelectOpts = []object.SelectOpts{
	{
		OptMap:         newsapi.Category,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-category-btn",
		DeleteButtonId: "delete-category-btn",
		PositionId:     "category",
		AlertMessage:   "You can only add up to 5 categories in a single query",
	},
	{
		OptMap:         newsapi.Country,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-country-btn",
		DeleteButtonId: "delete-country-btn",
		PositionId:     "country",
		AlertMessage:   "You can only add up to 5 countries in a single query",
	},
	{
		OptMap:         newsapi.Language,
		MaxDiv:         5,
		DefaultValue:   "",
		DefaultText:    "all",
		InsertButtonId: "insert-lang-btn",
		DeleteButtonId: "delete-lang-btn",
		PositionId:     "language",
		AlertMessage:   "You can only add up to 5 languages in a single query",
	},
}
