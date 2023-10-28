package gnews

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/parser"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
)

type Response struct {
	TotalArticles int       `json:"totalArticles"`
	Articles      []Article `json:"articles"`
	nextPageURL   string    `json:"-"`
}

func (resp Response) GetStatus() string {
	return "success"
}

func (resp Response) HasNext() bool {
	return resp.Len() > 0
}

func (resp Response) NextPageRequest(body io.Reader) (*http.Request, error) {
	if resp.HasNext() {
		return http.NewRequest(API_METHOD, resp.nextPageURL, body)
	}
	return nil, api.ErrNotNextPage
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

// convert response to model.CreateNewsParams and return by a channel
func (resp Response) ToNews(ctx context.Context, wg *sync.WaitGroup, c chan<- *service.NewsCreateRequest) {
	defer wg.Done()
	for i := 0; i < resp.Len(); i++ {
		select {
		case <-ctx.Done():
			break
		default:
			link := resp.Articles[i].Link
			u, _ := url.Parse(link)

			md5hash, _ := api.MD5Hash(
				resp.Articles[i].Title,
				resp.Articles[i].PublishedAt.ToTime(),
				resp.Articles[i].Content,
			)

			var req *service.NewsCreateRequest
			if val, ok := ctx.Value(api.QueryOriPageKey).(bool); ok && val {
				q := parser.ParseURL(u)
				req = q.ToNewsCreateParam(md5hash)
			} else {
				c <- &service.NewsCreateRequest{
					Md5Hash:     md5hash,
					Guid:        parser.ToGUID(u),
					Author:      []string{},
					Title:       resp.Articles[i].Title,
					Link:        link,
					Description: resp.Articles[i].Description,
					Language:    "",
					Content:     []string{resp.Articles[i].Content},
					Source:      u.Host,
					RelatedGuid: []string{},
					PublishedAt: resp.Articles[i].PublishedAt.ToTime(),
				}
			}
			c <- req
		}
	}
	return
}

type Article struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Content     string      `json:"content"`
	Link        string      `json:"url"`
	Image       string      `json:"image"`
	PublishedAt APIRespTime `json:"publishedAt"`
	Source      string      `json:"-"`
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

type ArticleSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func ParseHTTPResponse(resp *http.Response, currPage int) (*Response, error) {
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

	apiResponse, err := ParseResponse(body, currPage)
	if err != nil {
		return nil, err
	}

	if apiResponse.HasNext() {
		u, _ := url.Parse(resp.Request.URL.String())
		v, _ := url.ParseQuery(resp.Request.URL.Query().Encode())
		v.Set(qPage, strconv.Itoa(currPage+1))
		u.RawQuery = v.Encode()
		apiResponse.nextPageURL = u.String()
	}
	return apiResponse, err
}

func ParseErrorResponse(b []byte) (*APIError, error) {
	var respErr *APIError
	if err := json.Unmarshal(b, &respErr); err != nil {
		return nil, fmt.Errorf("error while unmarshaling body: %v", err)
	}
	return respErr, nil
}

func ParseResponse(b []byte, currPage int) (*Response, error) {
	var apiResponse Response
	if err := json.Unmarshal(b, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}