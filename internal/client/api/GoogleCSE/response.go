package googlecse

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/utils"
	"github.com/oklog/ulid/v2"
)

type Response struct {
	Kind string `json:"kind"`
	URL  struct {
		Type     string `json:"type"`
		Template string `json:"template"`
	}
	Queries           Queries           `json:"queries"`
	Context           map[string]string `json:"context"`
	SearchInformation SearchInformation `json:"searchInformation"`
	Items             []Item            `json:"items"`
}

func (resp Response) String() string {
	b, _ := json.MarshalIndent(resp, "", "\t")
	return string(b)
}

func (resp Response) ContentProcessFunc(c string) (string, error) {
	return c, nil
}

func (resp Response) GetStatus() string {
	return "success"
}

func (resp Response) HasNext() bool {
	return len(resp.Queries.NextPage) > 0
}

func (resp Response) Len() int {
	return len(resp.Items)
}

func (resp Response) ToNewsItemList() (api.NextPageToken, []api.NewsPreview) {
	preview := make([]api.NewsPreview, resp.Len())
	for i, item := range resp.Items {
		id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
		preview[i] = api.NewsPreview{
			Id:          id,
			Title:       item.Title,
			Link:        item.Link.String(),
			Description: item.PageMap.Description(),
			Category:    item.PageMap.Category(),
			Content:     "",
			PubDate:     item.PageMap.PubDate(),
		}
	}

	if resp.HasNext() {
		return api.IntNextPageToken(resp.Queries.NextPage[0].StartIndex), preview
	}
	return api.IntLastPageToken, preview
}

type SearchInformation struct {
	SearchTime            float32 `json:"searchTime"`
	FormattedSearchTime   float32 `json:"formattedSearchTime,string"`
	TotalResults          int     `json:"totalResults,string"`
	FormattedTotalResults string  `json:"formattedTotalResults"`
}

type Queries struct {
	Request      []Query `json:"request"`
	NextPage     []Query `json:"nextPage"`
	PreviousPage []Query `json:"previousPage"`
}

type Query struct {
	Title          string         `json:"title"`
	TotalResults   int            `json:"totalResults,string"`
	SearchTerms    string         `json:"searchTerms"`
	Count          int            `json:"count"`
	StartIndex     int            `json:"startIndex"`
	InputEncoding  string         `json:"inputEncoding"`
	OutputEncoding string         `json:"outputEncoding"`
	Other          map[string]any `json:"-"`
}

func (q Query) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	for k, v := range q.Other {
		m[k] = v
	}
	m["title"] = q.Title
	m["totalResults"] = q.TotalResults
	m["searchTerms"] = q.SearchTerms
	m["count"] = q.Count
	m["startIndex"] = q.StartIndex
	m["inputEncoding"] = q.InputEncoding
	m["outputEncoding"] = q.OutputEncoding
	return json.Marshal(m)
}

func (q *Query) UnmarshalJSON(data []byte) error {
	type InnerQuery Query

	tmp := struct{ *InnerQuery }{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	var other map[string]any
	if err := json.Unmarshal(data, &other); err != nil {
		return err
	}

	tags, _ := utils.GetStructTags("json", Query{})
	for _, tag := range tags {
		delete(other, tag)
	}

	tmp.InnerQuery.Other = other
	(*q) = Query(*tmp.InnerQuery)
	return nil
}

type Item struct {
	Kind             string  `json:"kind"`
	Title            string  `json:"title"`
	HtmlTitle        string  `json:"htmlTitle"`
	Link             URL     `json:"link"`
	DisplayLink      string  `json:"displayLink"`
	Snippet          string  `json:"snippet"`
	HtmlSnippet      string  `json:"htmlSnippet"`
	CacheId          string  `json:"cacheId"`
	FormattedUrl     string  `json:"formattedUrl"`
	HtmlFormattedUrl string  `json:"htmlFormattedUrl"`
	PageMap          PageMap `json:"pageMap"`
}

type URL struct{ *url.URL }

func (u URL) MarshalJSON() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u *URL) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("error while unquoting %s: %v", u.String(), err)
	}

	(*u).URL, err = url.Parse(s)
	return err
}

type PageMap struct {
	Website      []Website           `json:"website"`
	MetaTags     []map[string]string `json:"metatags"`
	CSEImage     []map[string]string `json:"cse_image"`
	CSEThumbnail []map[string]string `json:"cse_thumbnail"`
	ListItem     []map[string]string `json:"listitme"`
}

type Website struct {
	Image            string `json:"image"`
	Datemodified     string `json:"datemodified"`
	Keywords         string `json:"keywords"`
	Articlesection   string `json:"articlesection"`
	Author           string `json:"author"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Publisher        string `json:"publisher"`
	Headline         string `json:"headline"`
	Url              URL    `json:"url"`
	Datepublished    string `json:"datepublished"`
	MainentityOfPage string `json:"mainentityofpage"`
}

var pmTitle = []string{
	"title",
	"twitter:title",
	"og:title",
	"apple-mobile-web-app-title",
}

var pmLink = []string{
	"twitter:url",
	"og:url",
}

var pmDescription = []string{
	"twitter:description",
	"og:description",
}

var pmCategory = []string{
	"article:section",
}

var pmPubDate = []string{
	"pubdate",
	"article:published_time",
	"date",
}

func (pm PageMap) findInPageMap(defaultVal string, keys []string) string {
	f := defaultVal
	for i := 0; f == "" && i < len(keys); i++ {
		for _, m := range pm.MetaTags {
			if val, ok := m[keys[i]]; ok {
				f = val
				break
			}
		}
	}
	return f
}

func (pm PageMap) Title() string {
	title := ""
	if len(pm.Website) > 0 {
		title = pm.Website[0].Headline
	}
	return pm.findInPageMap(title, pmTitle)
}

func (pm PageMap) Link() string {
	link := ""
	if len(pm.Website) > 0 {
		link = pm.Website[0].Url.String()
	}
	return pm.findInPageMap(link, pmLink)
}

func (pm PageMap) Description() string {
	description := ""
	if len(pm.Website) > 0 {
		description = pm.Website[0].Description
	}
	return pm.findInPageMap(description, pmDescription)
}

func (pm PageMap) Category() string {
	return pm.findInPageMap("", pmCategory)
}

func (pm PageMap) PubDate() time.Time {
	s := ""
	if len(pm.Website) > 0 {
		s = pm.Website[0].Datepublished
	}

	s = pm.findInPageMap(s, pmPubDate)
	if s == "" {
		// pubdate not found
		return time.Time{}
	}

	tm, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// malformed pubdate
		return time.Time{}
	}
	return tm.UTC()
}

func ParseHTTPResponse(resp *http.Response) (api.Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiErrResponse, err := ParseErrorResponse(body)
		if err != nil {
			return nil, err
		}
		return nil, apiErrResponse.ToError()
	}

	apiResponse, err := ParseResponse(body)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func ParseResponse(b []byte) (*Response, error) {
	var resp Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("error while unmarshaling body of error response: %v", err)
	}
	return &resp, nil
}

func ParseErrorResponse(b []byte) (*ErrorResponse, error) {
	var respErr *ErrorResponse
	if err := json.Unmarshal(b, &respErr); err != nil {
		return nil, fmt.Errorf("error while unmarshaling body: %v", err)
	}
	return respErr, nil
}
