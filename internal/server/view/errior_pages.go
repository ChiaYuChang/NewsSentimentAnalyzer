package view

import (
	"fmt"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/spf13/viper"
)

var ErrorPage500 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent(), Title: "500 error"},
	ErrorCode:          500,
	ErrorMessage:       "Sorry, unexpected error",
	ErrorDetail:        "The server encountered an internal error or misconfiguration and was unable to complete your request.",
	ShouldAutoRedirect: false,
}
var ErrorPage400 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent(), Title: "400 error"},
	ErrorCode:          http.StatusBadRequest,
	ErrorMessage:       "Bad request",
	ErrorDetail:        "There was a problem with your request.",
	ShouldAutoRedirect: false,
}

var ErrorPage401 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent(), Title: "401 error"},
	ErrorCode:          http.StatusUnauthorized,
	ErrorMessage:       "Unauthorized",
	ErrorDetail:        "You are not authorized to access this page.",
	ShouldAutoRedirect: true,
	RedirectPageUrl:    global.AppVar.App.RoutePattern.Page["sign-in"],
	RedirectPageName:   "log in page",
	CountDownFrom:      5,
}

var ErrorPage403 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent(), Title: "403 error"},
	ErrorCode:          http.StatusForbidden,
	ErrorMessage:       "Access denied",
	ErrorDetail:        "You do not have premission to access this page.",
	ShouldAutoRedirect: false,
}

var ErrorPage404 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent(), Title: "404 error"},
	ErrorCode:          http.StatusNotFound,
	ErrorMessage:       "Page not found",
	ErrorDetail:        "The page you are looking for may have been moved, deleted, or possibly never existed.",
	ShouldAutoRedirect: false,
}

var ErrorPage410 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent(), Title: "410 error"},
	ErrorCode:          http.StatusGone,
	ErrorMessage:       "Gone",
	ErrorDetail:        "The requested resource is no longer available at the server and no forwarding address is known.",
	ShouldAutoRedirect: true,
	RedirectPageUrl:    fmt.Sprintf("/%s%s", viper.GetString("APP_API_VERSION"), global.AppVar.App.RoutePattern.Page["endpoints"]),
	RedirectPageName:   "welcome page",
	CountDownFrom:      5,
}

var ErrorPage429 = object.ErrorPage{
	Page:               object.Page{HeadConent: SharedHeadContent(), Title: "409 error"},
	ErrorCode:          http.StatusTooManyRequests,
	ErrorMessage:       "Too Many Requests",
	ErrorDetail:        "You have sent too many requests to us recently. Please try again later.",
	ShouldAutoRedirect: false,
}
