package object

import (
	"fmt"
	"html/template"

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
	Username                  string
}

type SignUpPage struct {
	Page
	ShowUsernameHasUsedAlert bool
}

type APIKeyPage struct {
	Page
	APIOption
	NewsAPIs     []*APIKey
	AnalyzerAPIs []*APIKey
}

type APIOption struct {
	Source   map[int16]string
	Analyzer map[int16]string
}

type APIKey struct {
	ID   int16
	Name string
	Icon string
	Key  string
}

func (apikey APIKey) InputID() string {
	return fmt.Sprintf("api-id-%03d", apikey.ID)
}

func APIKeyFromDBModel(page Page, rows []*model.ListAPIKeyRow) APIKeyPage {
	apiKeyPage := APIKeyPage{
		Page:         page,
		NewsAPIs:     []*APIKey{},
		AnalyzerAPIs: []*APIKey{},
	}

	for _, row := range rows {
		var which *[]*APIKey
		if row.Type == model.ApiTypeSource {
			which = &(apiKeyPage.NewsAPIs)
		} else if row.Type == model.ApiTypeLanguageModel {
			which = &(apiKeyPage.AnalyzerAPIs)
		}

		(*which) = append((*which), &APIKey{
			ID:   row.ApiID,
			Name: row.Name,
			Icon: row.Image.String,
			Key:  row.Key,
		})
	}
	return apiKeyPage
}

type ChangePasswordPage struct {
	Page
	ShowPasswordNotMatchAlert         bool
	ShowShouldNotUsedOldPasswordAlert bool
}

type APIEndpointPage struct {
	Page
	Endpoints           map[string]*APIEndpoint
	NoAvailableEndpoint bool
}

type APIEndpoint struct {
	Name        string
	Image       *HTMLElement
	DocumentURL string
	Endpoints   *HTMLElementList
}

func APIEndpointFromDBModel(page Page, rows []*model.ListEndpointByOwnerRow) APIEndpointPage {
	apiEndpointPage := APIEndpointPage{
		Page:      page,
		Endpoints: make(map[string]*APIEndpoint, len(rows)),
	}

	if len(rows) < 1 {
		apiEndpointPage.NoAvailableEndpoint = true
		return apiEndpointPage
	}

	for _, row := range rows {
		if _, ok := apiEndpointPage.Endpoints[row.ApiName]; !ok {
			apiEndpointPage.Endpoints[row.ApiName] = &APIEndpoint{
				Name:        row.EndpointName,
				Image:       NewHTMLElement("img"),
				DocumentURL: row.DocumentUrl,
				Endpoints:   NewHTMLElementList("button"),
			}
		}
		apiEndpointPage.
			Endpoints[row.ApiName].
			Image.
			AddPair("src", "/static/image/"+row.Image.String).
			AddPair("alt", row.ApiName).
			AddPair("class", "api-logo api-logo-large")

		apiEndpointPage.
			Endpoints[row.ApiName].
			Endpoints.
			NewHTMLElement().
			AddPair("type", "button").
			AddPair("class", "btn").
			AddPair("oneclick", "#").
			ToOpeningElement(template.HTML(row.EndpointName))
	}

	return apiEndpointPage
}
