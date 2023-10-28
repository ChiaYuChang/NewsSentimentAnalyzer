package googlecse

import (
	"fmt"
	"net/url"
	"strings"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
)

func init() {
	pfs := []pageform.PageForm{GoogleCSE{}}

	for _, pf := range pfs {
		pageform.Add(pf)
	}
}

const API_NAME = "Google API"
const API_VERSION = "v1"

const EPCustomSearch string = "Custom Search"

type GoogleCSE struct {
	Keyword           string `form:"keyword"`
	SearchEngineID    string `form:"search-engine-id"     validate:"required"`
	DateRestrictValue int    `form:"date-restrict-value"  validate:"gte=0"`
	DateRestrictUnit  string `form:"date-restrict-unit"`
}

func (f GoogleCSE) DateRestrict() string {
	if f.DateRestrictValue == 0 {
		return ""
	}
	return fmt.Sprintf("%d%s", f.DateRestrictValue, f.DateRestrictUnit)
}

func (f GoogleCSE) Endpoint() string {
	return EPCustomSearch
}

func (f GoogleCSE) API() string {
	return API_NAME
}

func (f GoogleCSE) String() string {
	sb := strings.Builder{}
	sb.WriteString("GoogleCSE:")
	sb.WriteString(fmt.Sprintf("\t- Keywords        : %s\n", f.Keyword))
	sb.WriteString(fmt.Sprintf("\t- Search Engine ID: %s\n", f.Keyword))
	if f.DateRestrictValue > 0 {
		sb.WriteString(fmt.Sprintf("\t- DateRestrict: %d%s\n", f.DateRestrictValue, f.DateRestrictUnit))
	}
	return sb.String()
}

func (f GoogleCSE) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[GoogleCSE](decoder, val, postForm)
}
