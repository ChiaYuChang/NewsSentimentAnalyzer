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
	"strings"
	"sync"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/parser"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
)

type Response struct {
	TotalArticles int       `json:"totalArticles"`
	Articles      []Article `json:"articles"`
	CurrPage      int       `json:"-"`
}

func (resp Response) ContentProcessFunc(c string) (string, error) {
	content := global.CLSToken + c
	content = strings.ReplaceAll(content, "。\n\r", "\n")
	content = strings.ReplaceAll(content, "。\n", "。"+global.SEPToken)
	content = strings.ReplaceAll(content, "\n", "")
	content = global.CLSToken + content
	return content, nil
}

func (resp Response) GetStatus() string {
	return "success"
}

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

func (resp Response) ToNewsItemList() (next api.NextPageToken, preview []api.NewsPreview) {
	if !resp.HasNext() {
		return api.IntLastPageToken, nil
	}
	preview = make([]api.NewsPreview, resp.Len())
	for i, article := range resp.Articles {
		content, err := resp.ContentProcessFunc(article.Content)
		if err != nil {
			global.Logger.Error().Err(err).Msg("content processing failed")
			continue
		}
		preview[i] = api.NewsPreview{
			Id:          i,
			Title:       article.Title,
			Link:        article.Link,
			Description: article.Description,
			Category:    "",
			Content:     content,
			PubDate:     time.Time(article.PublishedAt),
		}
	}
	return api.IntNextPageToken(resp.CurrPage + 1), preview
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
	Source      struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"source"`
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
		return nil, err
	}
	return &apiResponse, nil
}
