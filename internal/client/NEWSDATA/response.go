package newsdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type Response struct {
	Status      string    `json:"status"`
	TotalResult int       `json:"totalResults"`
	Articles    []Article `json:"results"`
	NextPage    string    `json:"nextPage"`
	NextPageUrl *url.URL  `json:"-"`
}

func (resp *Response) HasNextPage() bool {
	return resp.NextPage != ""
}

func (resp Response) GetStatus() string {
	return resp.Status
}

func (resp Response) Len() int {
	return len(resp.Articles)
}

func (resp Response) ToNews(c chan<- client.News) {
	news := make([]client.News, resp.Len())
	for i := 0; i < resp.Len(); i++ {
		news[i] = client.News{
			MD5Hash: client.MD5Hash(
				resp.Articles[i].Title,
				resp.Articles[i].PublishedAt.ToTime(),
			),
			Title:       resp.Articles[i].Title,
			Url:         resp.Articles[i].Url,
			Description: resp.Articles[i].Description,
			Content:     resp.Articles[i].Content,
			Source:      resp.Articles[i].SourceId,
			PublishAt:   resp.Articles[i].PublishedAt.ToTime().UTC(),
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

const APIRespTimeFormat = "2006-01-02 15:04:05"

func (respTime *APIRespTime) UnmarshalJSON(b []byte) error {
	s := bytes.Trim(b, "\"")
	t, err := time.Parse(APIRespTimeFormat, string(s))
	if err != nil {
		return err
	}
	*respTime = APIRespTime(t)
	return nil
}

func (tm APIRespTime) ToTime() time.Time {
	return time.Time(tm)
}

type APIError struct {
	Status string            `json:"status"`
	Code   int               `json:"-"`
	Result map[string]string `json:"results"`
}

// see https://newsdata.io/documentation/#http-response-codes
func (apiErr APIError) ToError() error {
	var ecCode ec.ErrorCode
	switch apiErr.Code {
	case 200:
		return ec.MustGetErr(ec.Success)
	case 401:
		ecCode = ec.ECUnauthorized
	case 403:
		ecCode = ec.ECForbidden
	case 409:
		ecCode = ec.ECConflict
	case 415:
		ecCode = ec.ECUnsupportedMediaType
	case 422:
		ecCode = ec.ECUnprocessableContent
	case 429:
		ecCode = ec.ECTooManyRequests
	case 500:
		ecCode = ec.ECServerError
	}

	if apiErr.Result["message"] != "" {
		return ec.MustGetErr(ecCode).(*ec.Error).
			WithDetails(apiErr.Result["message"])
	}
	return ec.MustGetErr(ecCode)
}

func (apiErr APIError) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	sb.WriteString("\t- Status Code: " + strconv.Itoa(apiErr.Code) + "\n")
	if len(apiErr.Result) > 0 {
		sb.WriteString("\t- Result:\n")
		for key, val := range apiErr.Result {
			sb.WriteString(fmt.Sprintf("\t  - %s: %s\n", key, val))
		}
	}
	return sb.String()
}

func ParseHTTPResponse(resp *http.Response) (*Response, error) {
	var respObj Response
	var err error

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		json.Unmarshal(body, &respObj)
		if respObj.HasNextPage() {
			resp.Request.URL.Query().Set("page", respObj.NextPage)
			respObj.NextPageUrl = resp.Request.URL
		}
		return &respObj, nil
	} else {
		var respErr APIError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshaling body: %w", err)
		}
		respErr.Code = resp.StatusCode
		return nil, respErr.ToError()
	}
}
