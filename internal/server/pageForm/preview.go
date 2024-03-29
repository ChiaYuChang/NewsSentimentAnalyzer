package pageform

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
)

type PreviewResponse struct {
	Error   *PreviewError     `form:"error"    json:"error,omitempty"`
	HasNext bool              `form:"has_next" json:"has_next"`
	Items   []api.NewsPreview `form:"items"    json:"items"`
}

type PreviewError struct {
	Code        int      `form:"code"     json:"code,omitempty"`
	PgxCode     string   `form:"pgx_code" json:"pgx_code,omitempty"`
	Message     string   `form:"message"  json:"message,omitempty"`
	Detail      []string `form:"details"  json:"details,omitempty"`
	RedirectURL string   `form:"url"      json:"url,omitempty"`
}

type PreviewPostForm struct {
	SelectAll bool     `form:"select_all" json:"select_all"`
	Item      []string `form:"item"       json:"item"`
}

type PreviewPostResp struct {
	Error       *PreviewError `form:"error" json:"error,omitempty"`
	RedirectURL string        `form:"url"   json:"url,omitempty"`
}

func (resp *PreviewPostResp) WithEcError(ecErr *ec.Error) *PreviewPostResp {
	resp.Error = &PreviewError{
		Code:    ecErr.HttpStatusCode,
		PgxCode: ecErr.PgxCode,
		Message: ecErr.Message,
		Detail:  ecErr.Details,
	}
	return resp
}

func (resp PreviewPostResp) HttpStatusCode() int {
	return resp.Error.Code
}

func (resp *PreviewPostResp) WithRedirectURL(url string) *PreviewPostResp {
	resp.RedirectURL = url
	return resp
}

func (resp *PreviewPostResp) WithOutDetails() *PreviewPostResp {
	resp.Error.Detail = nil
	return resp
}

func (resp *PreviewPostResp) WithOutMessage() *PreviewPostResp {
	resp.Error.Message = ""
	return resp
}
