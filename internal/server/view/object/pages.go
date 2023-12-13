package object

import (
	"fmt"
	"html/template"
	"path"
	"strings"
	"sync"
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
	NJobs       map[string]int `json:"n_jobs"`
	TotalJobKey string         `json:"total_job_key"`
	PageSize    int            `json:"page_size"`
}

func (p APIResultPage) TotalJobs() int {
	return p.NJobs[p.TotalJobKey]
}

func (p APIResultPage) NCanceled() int {
	return p.NJobs[string(model.JobStatusCanceled)]
}

func (p APIResultPage) NCreated() int {
	return p.NJobs[string(model.JobStatusCreated)]
}

func (p APIResultPage) NDone() int {
	return p.NJobs[string(model.JobStatusDone)]
}

func (p APIResultPage) NFailed() int {
	return p.NJobs[string(model.JobStatusFailed)]
}

func (p APIResultPage) NRunning() int {
	return p.NJobs[string(model.JobStatusRunning)]
}

func (p APIResultPage) NPage(js string) int {
	n := p.NJobs[js]
	q := n / p.PageSize
	if n%p.PageSize > 0 {
		q++
	}
	return q
}

type Job struct {
	Id        int32  `json:"job-id"`
	Status    string `json:"job-status"`
	NewsSrc   string `json:"job-news_src"`
	Analyzer  string `json:"job-analyzer"`
	CreatedAt string `json:"job-created_at"`
	UpdatedAt string `json:"job-updated_at"`
}

func StatusToClass(jobStatus model.JobStatus) string {
	switch jobStatus {
	case model.JobStatusCreated:
		return "created"
	case model.JobStatusRunning:
		return "running"
	case model.JobStatusDone:
		return "done"
	case model.JobStatusFailed:
		return "failed"
	case model.JobStatusCanceled:
		return "canceled"
	}
	return ""
}

type JobDetails struct {
	JID           int64  `json:"job-id"`
	Owner         string `json:"job-owner"`
	Status        string `json:"job-status"`
	NewsAPI       string `json:"job-news_api"`
	NewsAPIQuery  string `json:"job-news_api_query"`
	Analyzer      string `json:"job-analyzer"`
	AnalyzerQuery string `json:"job-analyzer_query"`
	CreatedAt     string `json:"job-created_at"`
	UpdatedAt     string `json:"job-updated_at"`
}

func NewJobDetails(owner string, j *model.GetJobsByJobIdRow) JobDetails {
	return JobDetails{
		JID:           j.ID,
		Owner:         owner,
		Status:        StatusToClass(j.Status),
		NewsAPI:       j.NewsSrc,
		NewsAPIQuery:  j.SrcQuery,
		Analyzer:      j.Analyzer,
		AnalyzerQuery: string(j.LlmQuery),
		CreatedAt:     j.CreatedAt.Time.UTC().Format(time.DateTime),
		UpdatedAt:     j.UpdatedAt.Time.UTC().Format(time.DateTime),
	}
}

type APIAdminPage struct {
	Page
}

type PasswordInput struct {
	IdPrefix              string
	Name                  string
	PlaceHolder           string
	PasswordStrengthCheck bool
	PasswordCreteria      []PasswordCreterion
	AlertMessage          string
}

func (p PasswordInput) ShowAlert() bool {
	return p.AlertMessage != ""
}

type PasswordCreterion struct {
	Id       string
	Name     string
	Min, Max int
	Class    []string
	Regx     template.JS
}

type passwordCreteriaSingleton struct {
	PasswordCreterion []PasswordCreterion
	sync.Once
}

var defaultPasswordCreteria passwordCreteriaSingleton

func GetDefaultPasswordCreteria() []PasswordCreterion {
	defaultPasswordCreteria.Do(func() {
		defaultPasswordCreteria.PasswordCreterion = []PasswordCreterion{
			{
				Id:    "new_pwd_length",
				Name:  "length",
				Min:   global.AppVar.Password.MinLength,
				Max:   global.AppVar.Password.MaxLength,
				Class: []string{"invalid"}, Regx: "/./g",
			},
			{
				Id:    "new_pwd_n_lower",
				Name:  "lower case",
				Min:   global.AppVar.Password.MinNumLower,
				Max:   -1,
				Class: []string{"invalid"}, Regx: "/[a-z]/g",
			},
			{
				Id:    "new_pwd_n_upper",
				Name:  "upper case",
				Min:   global.AppVar.Password.MinNumUpper,
				Max:   -1,
				Class: []string{"invalid"}, Regx: "/[A-Z]/g",
			},
			{
				Id:    "new_pwd_n_number",
				Name:  "digit",
				Min:   global.AppVar.Password.MinNumDigit,
				Max:   -1,
				Class: []string{"invalid"}, Regx: `/\d/g`,
			},
			{
				Id:    "new_pwd_n_special",
				Name:  "special character",
				Min:   global.AppVar.Password.MinNumSpecial,
				Max:   -1,
				Class: []string{"invalid"}, Regx: `/[-#$.%&@!+=<>*\\/]/g`,
			},
		}
	})
	return defaultPasswordCreteria.PasswordCreterion
}

func (c PasswordCreterion) ClassList() string {
	return strings.Join(c.Class, " ")
}

func (c PasswordCreterion) Message() string {
	if c.Max < 0 {
		if c.Min == 1 {
			return fmt.Sprintf("Must contain at least %d %s", c.Min, c.Name)
		}
		return fmt.Sprintf("Must contain at least %d %ss", c.Min, c.Name)
	}
	return fmt.Sprintf("Must be between %d and %d %ss", c.Min, c.Max, c.Name)
}

type APIKeyPage struct {
	Page
	APIOption
	APIVersion   string
	NewsAPIs     []*APIKey
	AnalyzerAPIs []*APIKey
}

type APIOption struct {
	Source   map[int16]string `json:"source"`
	Analyzer map[int16]string `json:"analyzer"`
}

type APIKey struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Key  string `json:"key"`
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
	OldPassword PasswordInput
	NewPassword PasswordInput
	// ShowPasswordNotMatchAlert         bool
	// ShowShouldNotUsedOldPasswordAlert bool
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
	API        string
	SelectOpts []SelectOpts
	Version    string
	Endpoint   string
}

func (apiEP APIEndpointPage) HasSelectOpts() bool {
	return len(apiEP.SelectOpts) > 0
}

type ResultSecectorPage struct {
	Page
	Version string
}

type AnalyzerPage struct {
	Page
	Prompt  map[string]string
	Version string
}
