package newsdata

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/oklog/ulid/v2"
)

type Response struct {
	Status      string    `json:"status"`
	TotalResult int       `json:"totalResults"`
	Articles    []Article `json:"results"`
	NextPage    string    `json:"nextPage"`
}

func (resp Response) ContentProcessFunc(c string) (string, error) {
	return c, nil
}

// return response status (must be either "success" or "error")
func (resp Response) GetStatus() string {
	return resp.Status
}

// return url for querying next page
func (resp *Response) HasNext() bool {
	return resp.NextPage != ""
}

// return the number of the articles in the response
func (resp Response) Len() int {
	return len(resp.Articles)
}

func (resp Response) String() string {
	b, _ := json.MarshalIndent(resp, "", "\t")
	return string(b)
}

func (resp Response) ToNewsItemList() (next api.NextPageToken, preview []api.NewsPreview) {
	preview = make([]api.NewsPreview, len(resp.Articles))
	for i, article := range resp.Articles {
		id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
		preview[i] = api.NewsPreview{
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
		return api.StrNextPageToken(resp.NextPage), preview
	}
	return api.StrLastPageToken, preview
}

type Article struct {
	Title       string      `json:"title"`
	Link        string      `json:"link"`
	SourceId    string      `json:"source_id"`
	Keywords    []string    `json:"keywords"`
	Author      []string    `json:"creater"`
	UrlToImage  string      `json:"image_url"`
	UrlToVideo  string      `json:"video_url"`
	Description string      `json:"description"`
	PublishedAt APIRespTime `json:"pubDate"`
	Content     string      `json:"content"`
	Country     []string    `json:"country"`
	Category    []string    `json:"category"`
	Language    string      `json:"language"`
}

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

func ParseHTTPResponse(resp *http.Response) (*Response, error) {
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
		return nil, apiErrResponse.ToEcError(resp.StatusCode)
	}

	apiResponse, err := ParseResponse(body)
	if err != nil {
		return nil, err
	}

	return apiResponse, err
}

func ParseErrorResponse(b []byte) (*ErrorResponse, error) {
	var respErr *ErrorResponse
	err := json.Unmarshal(b, &respErr)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling body: %v", err)
	}
	return respErr, nil
}

func ParseResponse(b []byte) (*Response, error) {
	var apiResponse Response
	err := json.Unmarshal(b, &apiResponse)
	return &apiResponse, err
}
