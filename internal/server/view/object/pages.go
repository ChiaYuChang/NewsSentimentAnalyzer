package object

import (
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
)

type Page struct {
	HeadConent
	Title string
}

type WelcomePage struct {
	Page
	Name string
	Role string
}

type LoginPage struct {
	Page
	ShowUsernameNotFountAlert bool
	ShowPasswordMismatchAlert bool
}

type SignUpPage struct {
	Page
	ShowUsernameHasUsedAlert bool
}

type APIKeyPage struct {
	Page
	NewsAPIs     []*APIKey
	AnalyzerAPIs []*APIKey
}

type APIKey struct {
	Image *HTMLElement
	Input *HTMLElement
}

func APIKeyFromDBModel(page Page, rows []*model.ListAPIKeyRow) APIKeyPage {
	APIKeyPage := APIKeyPage{
		Page:         page,
		NewsAPIs:     []*APIKey{},
		AnalyzerAPIs: []*APIKey{},
	}

	for _, row := range rows {
		var which *[]*APIKey
		if row.Type == model.ApiTypeSource {
			which = &(APIKeyPage.NewsAPIs)
		} else if row.Type == model.ApiTypeLanguageModel {
			which = &(APIKeyPage.AnalyzerAPIs)
		}
		(*which) = append((*which), &APIKey{
			Image: NewHTMLElement("img").
				AddPair("alt", row.Name),
			Input: NewHTMLElement("input").
				AddPair("name", strings.ToLower(row.Name)+"-apikey").
				AddPair("id", strings.ToLower(row.Name)+"-apikey").
				AddPair("value", row.Key.String),
		})
	}
	return APIKeyPage
}

type ChangePasswordPage struct {
	Page
	ShowPasswordNotMatchAlert         bool
	ShowShouldNotUsedOldPasswordAlert bool
}

type APIEndpointPage struct {
	Page
	APIEndpoints []APIEndpoint
}

type APIEndpoint struct {
	Image       HTMLElement
	DocumentURL string
	Endpoints   *HTMLElementList
}
