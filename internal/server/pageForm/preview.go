package pageform

import "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"

type PreviewResponse struct {
	Error   *PreviewError     `form:"error"    json:"error,omitempty"`
	HasNext bool              `form:"has_next" json:"has_next"`
	Items   []api.NewsPreview `form:"items"    json:"items"`
}

type PreviewError struct {
	Code        int      `form:"code"    json:"code,omitempty"`
	Message     string   `form:"message" json:"message,omitempty"`
	Detail      []string `form:"details" json:"details,omitempty"`
	RedirectURL string   `form:"url"     json:"url,omitempty"`
}

type PreviewPostForm struct {
	SelectAll bool     `form:"select_all" json:"select_all"`
	Item      []string `form:"item"       json:"item"`
}

type PreviewPostResp struct {
	Error       *PreviewError `form:"error" json:"error,omitempty"`
	RedirectURL string        `form:"url"   json:"url,omitempty"`
}
