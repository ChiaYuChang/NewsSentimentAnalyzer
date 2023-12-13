package newsapi

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/oklog/ulid/v2"
)

type Response struct {
	Status      string    `json:"status"`
	TotalResult int       `json:"total_results"`
	Articles    []Article `json:"articles"`
	CurrPage    int       `json:"current_page"`
}

func (resp Response) ContentProcessFunc(c string) (string, error) {
	return c, nil
}

// Retrieve whether the request was successful or not
func (resp Response) GetStatus() string {
	return resp.Status
}

// Check whether next page requery should be done.
func (resp Response) HasNext() bool {
	return resp.Len() > 0
}

// return the number of the articles in the response
func (resp Response) Len() int {
	return len(resp.Articles)
}

// fmt.Stringer interface
func (resp Response) String() string {
	b, _ := json.MarshalIndent(resp, "", "\t")
	return string(b)
}

func (resp Response) ToNewsItemList() (api.NextPageToken, []api.NewsPreview) {
	prev := make([]api.NewsPreview, resp.Len())
	for i, article := range resp.Articles {
		id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
		prev[i] = api.NewsPreview{
			Id:          id,
			Title:       article.Title,
			Link:        article.Link,
			Description: article.Description,
			Category:    "",
			Content:     article.Content,
			PubDate:     time.Time(article.PublishedAt),
		}
	}

	if resp.HasNext() {
		return api.IntNextPageToken(resp.CurrPage + 1), prev
	}
	return api.IntLastPageToken, prev
}

type Article struct {
	Source      ArticleSource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Link        string        `json:"url"`
	UrlToImage  string        `json:"urlToImage"`
	PublishedAt APIRespTime   `json:"publishedAt"`
	Content     string        `json:"content"`
}

// deal with api http response time format
type APIRespTime time.Time

func (respTime *APIRespTime) UnmarshalJSON(b []byte) error {
	s := bytes.Trim(b, "\"")
	t, err := time.Parse(API_RESP_TIME_FMT, string(s))
	if err != nil {
		return err
	}
	*respTime = APIRespTime(t)
	return nil
}

func (tm APIRespTime) ToTime() time.Time {
	return time.Time(tm)
}

type ArticleSource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (as ArticleSource) String() string {
	return fmt.Sprintf("Source: %s (%s)", as.Name, as.Id)
}

func extractCurrPageFromResp(resp *http.Response) (int, error) {
	pageStr := resp.Request.URL.Query().Get(string(Page))
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, fmt.Errorf("error while converting page string to int: %w", err)
	}
	return page, nil
}

func ParseHTTPResponse(resp *http.Response) (api.Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.Logger.Info().
			Err(err).
			Msg("error while reading response body")
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiErrResponse, err := ParseErrorResponse(body)
		if err != nil {
			global.Logger.Info().
				Int("status code", resp.StatusCode).
				Err(err).
				Str("body", string(body)).
				Str("url", resp.Request.URL.String()).
				Msg("error while parsing error response")
			return nil, err
		}
		global.Logger.Info().
			Int("status code", resp.StatusCode).
			Msg("error")
		return nil, apiErrResponse.ToEcError(resp.StatusCode)
	}

	apiResponse, err := ParseResponse(body)
	if err != nil {
		global.Logger.Info().
			Err(err).
			Msg("error while parsing response")
		return nil, err
	}

	page, err := extractCurrPageFromResp(resp)
	if err != nil {
		global.Logger.Info().
			Err(err).
			Msg("error while extracting current page")
		return nil, err
	}
	apiResponse.CurrPage = page
	return apiResponse, err
}

func ParseErrorResponse(b []byte) (*ErrorResponse, error) {
	var respErr *ErrorResponse
	if err := json.Unmarshal(b, &respErr); err != nil {
		return nil, fmt.Errorf("error while unmarshaling body: %v", err)
	}
	return respErr, nil
}

func ParseResponse(b []byte) (*Response, error) {
	var apiResponse Response
	if err := json.Unmarshal(b, &apiResponse); err != nil {
		return nil, fmt.Errorf("error while unmarshaling body of error response: %v", err)
	}
	return &apiResponse, nil
}
