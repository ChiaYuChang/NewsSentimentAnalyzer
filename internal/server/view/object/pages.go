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
	Title string
	Name  string
	Role  string
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
		if !row.Type.Valid {
			continue
		} else if row.Type.ApiType == model.ApiTypeSource {
			which = &(APIKeyPage.NewsAPIs)
		} else if row.Type.ApiType == model.ApiTypeLanguageModel {
			which = &(APIKeyPage.AnalyzerAPIs)
		}
		(*which) = append((*which), &APIKey{
			Image: NewHTMLElement("img").
				AddPair("alt", row.Name.String),
			Input: NewHTMLElement("input").
				AddPair("name", strings.ToLower(row.Name.String)+"-apikey").
				AddPair("id", strings.ToLower(row.Name.String)+"-apikey").
				AddPair("value", row.Key),
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
