package main

import (
	"encoding/json"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
)

type APIList []APIItem

func (a APIList) N() int {
	return len(a) + 1
}

type APIItem struct {
	Id          int           `json:"id"            mod:"trim"`
	Name        string        `json:"name"          mod:"trim"  validate:"required"`
	Type        model.ApiType `json:"type"          mod:"trim"  validate:"required"`
	Image       string        `json:"image"         mod:"trim"  validate:"required"`
	Icon        string        `json:"icon"          mod:"trim"  validate:"required"`
	DocumentURL string        `json:"document_url"  mod:"trim"  validate:"required"`
	CreatedAt   time.Time     `json:"created_at"    mod:"default=2000-01-01T00:00:00+00:00"`
	UpdatedAt   time.Time     `json:"updated_at"    mod:"default=2000-01-01T00:00:00+00:00"`
	Probability float64       `json:"probability"   mod:"trim"  validate:"gte=0,lte=1"`
}

func (api APIItem) String() string {
	b, _ := json.MarshalIndent(api, "", "   ")
	return string(b)
}

var APIs = APIList{
	{
		Id:          1,
		Name:        "NEWSDATA.IO",
		Type:        APITypeSource,
		Image:       "logo_NEWSDATA.IO.png",
		Icon:        "favicon_NEWSDATA.IO.png",
		DocumentURL: "https://newsdata.io/documentation/",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		Probability: 0.5,
	},
	{
		Id:          2,
		Name:        "GNews",
		Type:        APITypeSource,
		Image:       "logo_GNews.png",
		Icon:        "favicon_GNews.ico",
		DocumentURL: "https://gnews.io/docs/v4",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		Probability: 0.5,
	},
	{
		Id:          3,
		Name:        "NEWS API",
		Type:        APITypeSource,
		Image:       "logo_NEWS_API.png",
		Icon:        "favicon_NEWS_API.ico",
		DocumentURL: "https://newsapi.org/docs/",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		Probability: 0.5,
	},
	{
		Id:          4,
		Name:        "Google API",
		Type:        APITypeSource,
		Image:       "logo_Google_Custom_Search.png",
		Icon:        "favicon-Google.png",
		DocumentURL: "https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		Probability: 0.5,
	},
	{
		Id:          5,
		Name:        "OpenAI",
		Type:        APITypeLLM,
		Image:       "logo_ChatGPT.svg",
		Icon:        "favicon_ChatGPT.ico",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		DocumentURL: "https://openai.com/blog/introducing-chatgpt-and-whisper-apis",
		Probability: 1.0,
	},
}
