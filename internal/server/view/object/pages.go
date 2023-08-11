package object

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
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
	PageSignOut      string
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

type APIResultPage struct {
	Page
	Jobs []Job
}

type Job struct {
	ID        int32     `json:"ID"`
	Status    JobStatus `json:"Status"`
	NewsSrc   string    `json:"NewsSrc"`
	Analyzer  string    `json:"Analyzer"`
	CreatedAt string    `json:"CreatedAt"`
	UpdatedAt string    `json:"UpdatedAt"`
}

func StatusToClass(jobStatus model.JobStatus) string {
	switch jobStatus {
	case model.JobStatusCreated:
		return "job-status-created"
	case model.JobStatusRunning:
		return "job-status-running"
	case model.JobStatusDone:
		return "job-status-done"
	case model.JobStatusFailure:
		return "job-status-faliure"
	case model.JobStatusCanceled:
		return "job-status-canceled"
	}
	return ""
}

func (p *APIResultPage) SetJobs(jl []*model.GetJobsByOwnerRow) *APIResultPage {
	sort.Slice(jl, func(i, j int) bool {
		return jl[i].UpdatedAt.Time.Before(jl[j].UpdatedAt.Time)
	})

	p.Jobs = make([]Job, len(jl))
	for i, j := range jl {
		p.Jobs[i] = Job{
			ID: j.ID,
			Status: JobStatus{
				Class: StatusToClass(j.Status),
				Text:  strings.ToTitle(string(j.Status)),
			},
			NewsSrc:   j.NewsSrc,
			Analyzer:  j.Analyzer,
			CreatedAt: j.CreatedAt.Time.UTC().Format(time.DateOnly),
			UpdatedAt: j.UpdatedAt.Time.UTC().Format(time.DateOnly),
		}
	}
	return p
}

type JobDetails struct {
	JID           int32             `json:"JobID"`
	Owner         string            `json:"Owner"`
	Status        JobStatus         `json:"Status"`
	NewsAPI       string            `json:"NewsAPI"`
	NewsAPIQuery  string            `json:"NewsAPIQuery"`
	Analyzer      string            `json:"Analyzer"`
	AnalyzerQuery map[string]string `json:"AnalyzerQuery"`
	CreatedAt     string            `json:"CreatedAt"`
	UpdatedAt     string            `json:"UpdatedAt"`
}

func NewJobDetails(owner string, j *model.GetJobsByJobIdRow) JobDetails {
	analyzerQuery := map[string]string{}
	json.Unmarshal(j.LlmQuery, &analyzerQuery)

	return JobDetails{
		JID:   j.ID,
		Owner: owner,
		Status: JobStatus{
			Class: StatusToClass(j.Status),
			Text:  strings.ToTitle(string(j.Status)),
		},
		NewsAPI:       j.NewsSrc,
		NewsAPIQuery:  j.SrcQuery,
		Analyzer:      j.Analyzer,
		AnalyzerQuery: analyzerQuery,
		CreatedAt:     j.CreatedAt.Time.UTC().Format(time.DateTime),
		UpdatedAt:     j.CreatedAt.Time.UTC().Format(time.DateTime),
	}
}

type JobStatus struct {
	Class string `json:"Class"`
	Text  string `json:"Text"`
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
			Icon: row.Image,
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
				global.AppVar.App.StaticFile.SubFolder["image"],
				row.Image,
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
				global.AppVar.App.RoutePattern.Page["endpoints"],
				fmt.Sprintf("%s-%s", row.ApiName, row.EndpointName))).
			// AddPair("onclick", fmt.Sprintf(
			// 	"location.href='.%s/%s'",
			// 	global.AppVar.Server.RoutePattern.Pages["endpoints"],
			// 	strings.TrimSuffix(row.TemplateName, ".gotmpl"))).
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
