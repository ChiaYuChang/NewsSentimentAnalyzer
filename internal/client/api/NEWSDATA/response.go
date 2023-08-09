package newsdata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
)

type Response struct {
	Status      string    `json:"status"`
	TotalResult int       `json:"totalResults"`
	Articles    []Article `json:"results"`
	NextPage    string    `json:"nextPage"`
	nextPageURL string    `json:"-"`
}

// return response status success/error
func (resp Response) GetStatus() string {
	return resp.Status
}

// return url for querying next page
func (resp *Response) HasNext() bool {
	return resp.NextPage != ""
}

func (resp *Response) NextPageRequest(body io.Reader) (*http.Request, error) {
	if resp.HasNext() {
		return http.NewRequest(API_METHOD, resp.nextPageURL, body)
	}
	return nil, api.ErrNotNextPage
}

// return the number of the articles in the response
func (resp Response) Len() int {
	return len(resp.Articles)
}

func (resp Response) String() string {
	b, _ := json.MarshalIndent(resp, "", "\t")
	return string(b)
}

// convert response to model.CreateNewsParams and return by a channel
func (resp Response) ToNews(ctx context.Context, wg *sync.WaitGroup, c chan<- *model.CreateNewsParams) {
	defer wg.Done()
	for i := 0; i < resp.Len(); i++ {
		select {
		case <-ctx.Done():
			break
		default:
			c <- &model.CreateNewsParams{
				Md5Hash: api.MD5Hash(
					resp.Articles[i].Title,
					resp.Articles[i].PublishedAt.ToTime(),
				),
				Title:       resp.Articles[i].Title,
				Url:         resp.Articles[i].Url,
				Description: resp.Articles[i].Description,
				Content:     resp.Articles[i].Content,
				Source:      convert.StrTo(resp.Articles[i].SourceId).PgText(),
				PublishAt:   convert.TimeTo(resp.Articles[i].PublishedAt.ToTime()).ToPgTimeStampZ(),
			}
		}
	}
	return
}

type Article struct {
	Title       string      `json:"title"`
	Url         string      `json:"link"`
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
		return nil, apiErrResponse.ToError(resp.StatusCode)
	}

	apiResponse, err := ParseResponse(body)
	if err != nil {
		return nil, err
	}

	if apiResponse.HasNext() {
		u, _ := url.Parse(resp.Request.URL.String())
		v, _ := url.ParseQuery(resp.Request.URL.Query().Encode())
		v.Add(qPage, apiResponse.NextPage)
		u.RawQuery = v.Encode()
		apiResponse.nextPageURL = u.String()
	}
	return apiResponse, err
}

func ParseErrorResponse(b []byte) (*APIError, error) {
	var respErr *APIError
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
