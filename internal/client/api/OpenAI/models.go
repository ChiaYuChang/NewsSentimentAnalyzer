package openai

import (
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/nullable"
)

type Model struct {
	Id      string
	Object  string
	Created time.Time
	OwnedBy string
}

type Premission struct {
	Id                 string                  `json:"id"`
	Object             string                  `json:"object"`
	Created            time.Time               `json:"created"`
	AllowCreateEngine  bool                    `json:"allow_create_engine"`
	AllowSampling      bool                    `json:"allow_sampling"`
	AllowLogprobs      bool                    `json:"allow_logprobs"`
	AllowSearchIndices bool                    `json:"allow_search_indices"`
	AllowView          bool                    `json:"allow_view"`
	AllowFineTuning    bool                    `json:"allow_fine_tuning"`
	Organization       string                  `json:"organization"`
	Group              nullable.String[string] `json:"group"`
	IsBlocking         bool                    `json:"is_blocking"`
}
