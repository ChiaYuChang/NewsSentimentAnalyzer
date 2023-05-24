package newsapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type Response struct {
	Status      string    `json:"status"`
	TotalResult int       `json:"totalResults"`
	Articles    []Article `json:"articles"`
}

func (resp Response) String() string {
	sb := strings.Builder{}
	sb.WriteString("API Response:\n")
	sb.WriteString(fmt.Sprintf("Status      : %s\n", resp.Status))
	sb.WriteString(fmt.Sprintf("Total Result: %d\n", resp.TotalResult))
	sb.WriteString("Articles:\n")
	for _, a := range resp.Articles {
		sb.WriteString(strings.Repeat("=", 30) + "\n")
		sb.WriteString(a.String())
	}
	sb.WriteString(strings.Repeat("=", 30) + "\n")
	return sb.String()
}

type Article struct {
	ArticleSource `json:"source"`
	Author        string    `json:"author"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Url           string    `json:"url"`
	UrlToImage    string    `json:"urlToImage"`
	PublishedAt   time.Time `json:"publishedAt"`
	Content       string    `json:"content"`
}

func (a Article) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Title      : %s\n", a.Title))
	sb.WriteString(fmt.Sprintf("Author     : %s Published at: %s\n", a.Author, a.PublishedAt.Format(API_TIME_FORMAT)))

	sb.WriteString("Description:\n")
	if len(a.Description) > 100 {
		sb.WriteString("\t" + a.Description[:100] + "...\n")
	} else {
		sb.WriteString("\t" + a.Description + "\n")
	}

	sb.WriteString("Content    :\n")
	if len(a.Content) > 100 {
		sb.WriteString("\t" + a.Content[:100] + "...\n")
	} else {
		sb.WriteString("\t" + a.Content + "\n")
	}
	return sb.String()
}

type ArticleSource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (as ArticleSource) String() string {
	return fmt.Sprintf("Source: %s (%s)", as.Name, as.Id)
}

type APIError struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// See https://newsapi.org/docs/errors
func (apiErr APIError) ToError() error {
	var ecCode ec.ErrorCode
	switch apiErr.Code {
	case 200:
		return ec.MustGetErr(ec.Success)
	case 401:
		ecCode = ec.ECUnauthorized
	case 429:
		ecCode = ec.ECTooManyRequests
	case 500:
		ecCode = ec.ECServerError
	default:
		ecCode = ec.ECBadRequest
	}

	if apiErr.Message != "" {
		return ec.MustGetErr(ecCode).(*ec.Error).
			WithDetails(apiErr.Message)
	}
	return ec.MustGetErr(ecCode)
}

func (apiErr APIError) String() string {
	sb := strings.Builder{}
	sb.WriteString("Api Error:\n")
	sb.WriteString("\t- Status Code: " + strconv.Itoa(apiErr.Code) + "\n")
	sb.WriteString("\t- Message    : " + apiErr.Message + "\n")
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
		return &respObj, nil
	} else {
		var respErr APIError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshaling body: %w", err)
		}
		return nil, respErr.ToError()
	}
}
