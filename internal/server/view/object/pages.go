package object

import (
	"fmt"
	"html/template"
	"path"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
)

type WelcomePage struct {
	Page
	Name             string
	Role             string
	PageChangePWD    string
	PageEndpoint     string
	PageManageAPIKey string
	PageSeeResult    string
	PageAdmin        string
	PageLogout       string
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

type APIAdminPage struct {
	Page
}

type APIKeyPage struct {
	Page
	APIOption
	APIVersion   string
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

type APIEndpointSelectionPage struct {
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

func APIEndpointFromDBModel(page Page, apiVer string, rows []*model.ListEndpointByOwnerRow) APIEndpointSelectionPage {
	apiEndpointPage := APIEndpointSelectionPage{
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
			AddPair("src", path.Join(
				"/static",
				global.AppVar.Server.RoutePattern.StaticFiles.Image,
				row.Image.String,
			)).
			AddPair("alt", row.ApiName).
			AddPair("class", "api-logo api-logo-large")

		apiEndpointPage.
			Endpoints[row.ApiName].
			Endpoints.
			NewHTMLElement().
			AddPair("type", "button").
			AddPair("class", "btn").
			AddPair("onclick", fmt.Sprintf(
				"location.href='.%s/%s'",
				global.AppVar.Server.RoutePattern.Pages["endpoints"],
				strings.TrimSuffix(row.TemplateName, ".gotmpl"))).
			ToOpeningElement(template.HTML(row.EndpointName))
	}

	return apiEndpointPage
}

type APIEndpointPage struct {
	Page
	API      string
	Version  string
	Endpoint string
}
