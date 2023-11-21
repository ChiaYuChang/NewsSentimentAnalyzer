package newsapi

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/parser"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
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

			if !parser.Has(u.Host) {
				continue
			}

			md5hash, _ := api.MD5Hash(
				resp.Articles[i].Title,
				resp.Articles[i].PublishedAt.ToTime(),
				resp.Articles[i].Content,
			)

			var q *parser.Query
			var req *service.NewsCreateRequest
			if val, ok := ctx.Value(api.QueryOriPageKey).(bool); ok && val {
				q = parser.ParseURL(u)
				req = q.ToNewsCreateParam(md5hash)
			} else {
				req = &service.NewsCreateRequest{
					Md5Hash:     md5hash,
					Guid:        parser.ToGUID(u),
					Author:      collection.NewCSL(resp.Articles[i].Author),
					Title:       resp.Articles[i].Title,
					Link:        resp.Articles[i].Link,
					Description: resp.Articles[i].Description,
					Language:    "",
					Content:     []string{resp.Articles[i].Content},
					Category:    "",
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

	page, err := extractCurrPageFromResp(resp)
	if err != nil {
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
