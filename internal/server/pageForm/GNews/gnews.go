package gnews

import (
	"fmt"
	"net/url"
	"strings"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
)

func init() {
	pfs := []pageform.PageForm{
		GNewsHeadlines{},
		GNewsSearch{},
	}

	for _, pf := range pfs {
		pageform.Add(pf)
	}

	for _, v := range []validator.Enmus[string]{
		CategoryValidator,
		CountryValidator,
		LanguageValidator,
	} {
		validator.Validate.RegisterValidation(
			v.Tag(),
			v.ValFun(),
		)
	}
}

const API_NAME = "GNews"

const (
	VAL_TAG_CATEGORY = "gnews_cat"
	VAL_TAG_LANGUAGE = "gnews_lang"
	VAL_TAG_COUNTRY  = "gnews_ctry"
)

const (
	EPSearch       string = "Search"
	EPTopHeadlines string = "Top Headlines"
)

type GNewsHeadlines struct {
	pageform.TimeRange
	Keyword  string   `form:"keyword"`
	Language []string `form:"language" validate:"gnews_lang"`
	Country  []string `form:"country"  validate:"gnews_ctry"`
	Category []string `form:"category" validate:"gnews_cat"`
}

func (f GNewsHeadlines) Endpoint() string {
	return EPTopHeadlines
}

func (f GNewsHeadlines) API() string {
	return API_NAME
}

func (f GNewsHeadlines) String() string {
	sb := strings.Builder{}
	sb.WriteString("GNewsHeadlines:\n")
	sb.WriteString(fmt.Sprintf("\t- Keywords: %s\n", f.Keyword))
	sb.WriteString(fmt.Sprintf("\t- Category: %s\n", strings.Join(f.Category, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Country : %s\n", strings.Join(f.Country, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Language: %s\n", strings.Join(f.Language, ", ")))
	sb.WriteString(f.TimeRange.ToString("\t"))
	return sb.String()
}

func (f GNewsHeadlines) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[GNewsHeadlines](decoder, val, postForm)
}

type GNewsSearch struct {
	pageform.SearchIn
	pageform.TimeRange
	Keyword  string   `mod:"trim"  form:"keyword"`
	Language []string `            form:"language" validate:"gnews_lang"`
	Country  []string `            form:"country"  validate:"gnews_ctry"`
}

func (f GNewsSearch) Endpoint() string {
	return EPSearch
}

func (f GNewsSearch) API() string {
	return API_NAME
}

func (f GNewsSearch) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[GNewsSearch](decoder, val, postForm)
}

func (f GNewsSearch) String() string {
	sb := strings.Builder{}
	sb.WriteString("GNewsSearch:\n")
	sb.WriteString(fmt.Sprintf("\t- Keywords : %s\n", f.Keyword))
	sb.WriteString(fmt.Sprintf("\t- Search In: %v\n", f.SearchIn))
	sb.WriteString(fmt.Sprintf("\t- Country  : %s\n", strings.Join(f.Country, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Language : %s\n", strings.Join(f.Language, ", ")))
	sb.WriteString(f.TimeRange.ToString("\t"))
	return sb.String()
}
