package gnews

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/oklog/ulid/v2"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
)

type Response struct {
	TotalArticles int       `json:"totalArticles"`
	Articles      []Article `json:"articles"`
	Max           int       `json:"-"`
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
	return (resp.Len() > 0 && resp.Len() == resp.Max)
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
	preview = make([]api.NewsPreview, resp.Len())
	for i, article := range resp.Articles {
		content, err := resp.ContentProcessFunc(article.Content)
		if err != nil {
			global.Logger.Error().Err(err).Msg("content processing failed")
			continue
		}
		id, _ := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
		preview[i] = api.NewsPreview{
			Id:          id,
			Title:       article.Title,
			Link:        article.Link,
			Description: article.Description,
			Category:    "",
			Content:     content,
			PubDate:     time.Time(article.PublishedAt),
		}
	}

	if !resp.HasNext() {
		return api.IntLastPageToken, preview
	}
	return api.IntNextPageToken(resp.CurrPage + 1), preview
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

func extractMaxFromResp(resp *http.Response) (int, error) {
	maxStr := resp.Request.URL.Query().Get(string(Max))
	if maxStr == "" {
		maxStr = "10"
	}

	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return 0, fmt.Errorf("error while converting max string to int: %w", err)
	}
	return max, nil
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

	max, err := extractMaxFromResp(resp)
	if err != nil {
		return nil, err
	}
	apiResponse.Max = max
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
