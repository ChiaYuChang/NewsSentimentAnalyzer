package newsapi

import (
	"fmt"
	"net/url"
	"strings"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/go-playground/form"

	val "github.com/go-playground/validator/v10"
)

func init() {
	pfs := []pageform.PageForm{
		NEWSAPITopHeadlines{},
		NEWSAPIEverything{},
		NEWSAPISources{},
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

const API_NAME = "NEWS API"

const (
	EPEverything   string = "Everything"
	EPTopHeadlines        = "Top Headlines"
	EPSources             = "Sources"
)

func SelectionOpts() []object.SelectOpts {
	return []object.SelectOpts{
		{
			OptMap:         Category,
			MaxDiv:         5,
			DefaultValue:   "",
			DefaultText:    "all",
			InsertButtonId: "insert-category-btn",
			DeleteButtonId: "delete-category-btn",
			PositionId:     "category",
			AlertMessage:   "You can only add up to 5 categories in a single query",
		},
		{
			OptMap:         Country,
			MaxDiv:         5,
			DefaultValue:   "",
			DefaultText:    "all",
			InsertButtonId: "insert-country-btn",
			DeleteButtonId: "delete-country-btn",
			PositionId:     "country",
			AlertMessage:   "You can only add up to 5 countries in a single query",
		},
		{
			OptMap:         Language,
			MaxDiv:         5,
			DefaultValue:   "",
			DefaultText:    "all",
			InsertButtonId: "insert-lang-btn",
			DeleteButtonId: "delete-lang-btn",
			PositionId:     "language",
			AlertMessage:   "You can only add up to 5 languages in a single query",
		},
	}
}

type NEWSAPIEverything struct {
	pageform.SearchIn
	pageform.TimeRange
	Keyword        string `mod:"trim" form:"keyword"`
	Sources        string `mod:"trim" form:"sources"`
	Domains        string `mod:"trim" form:"domains"`
	ExcludeDomains string `mod:"trim" form:"exclude-domains"`
	Language       string `mod:"trim" form:"language" val:"newsapi_lang"`
}

func (f NEWSAPIEverything) Endpoint() string {
	return EPEverything
}

func (f NEWSAPIEverything) API() string {
	return API_NAME
}

func (f NEWSAPIEverything) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[NEWSAPIEverything](decoder, val, postForm)
}

func (f NEWSAPIEverything) String() string {
	sb := strings.Builder{}
	sb.WriteString("NEWSAPISources:\n")
	sb.WriteString(fmt.Sprintf("\t- Keyword  : %s\n", f.Keyword))
	sb.WriteString(fmt.Sprintf("\t- Sources  : %s\n", f.Sources))
	sb.WriteString(fmt.Sprintf("\t- Search in: %v\n", f.SearchIn))
	sb.WriteString(fmt.Sprintf("\t- Domains  : %s\n", f.Domains))
	sb.WriteString(fmt.Sprintf("\t- eDomains : %s\n", f.ExcludeDomains))
	sb.WriteString(fmt.Sprintf("\t- Language : %s\n", f.Language))
	sb.WriteString(f.TimeRange.ToString("\t"))
	return sb.String()
}

func (f NEWSAPIEverything) Key() pageform.PageFormRepoKey {
	return pageform.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f NEWSAPIEverything) SelectionOpts() []object.SelectOpts {
	return SelectionOpts()
}

type NEWSAPISources struct {
	Language string `mod:"trim" form:"language" val:"newsapi_lang"`
	Country  string `mod:"trim" form:"country"  val:"newsapi_ctry"`
	Category string `mod:"trim" form:"category" val:"newsapi_cat"`
}

func (f NEWSAPISources) Endpoint() string {
	return EPSources
}

func (f NEWSAPISources) API() string {
	return API_NAME
}

func (f NEWSAPISources) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[NEWSAPISources](decoder, val, postForm)
}

func (f NEWSAPISources) String() string {
	sb := strings.Builder{}
	sb.WriteString("NEWSAPISources:\n")
	sb.WriteString(fmt.Sprintf("\t- Category: %s\n", f.Category))
	sb.WriteString(fmt.Sprintf("\t- Country : %s\n", f.Country))
	sb.WriteString(fmt.Sprintf("\t- Language: %s\n", f.Language))
	return sb.String()
}

func (f NEWSAPISources) Key() pageform.PageFormRepoKey {
	return pageform.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f NEWSAPISources) SelectionOpts() []object.SelectOpts {
	return SelectionOpts()
}

type NEWSAPITopHeadlines struct {
	Keyword  string `mod:"trim" form:"keyword"`
	Sources  string `mod:"trim" form:"sources"`
	Country  string `mod:"trim" form:"country"`
	Category string `mod:"trim" form:"category"`
}

func (f NEWSAPITopHeadlines) Endpoint() string {
	return EPTopHeadlines
}

func (f NEWSAPITopHeadlines) API() string {
	return API_NAME
}

func (f NEWSAPITopHeadlines) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[NEWSAPITopHeadlines](decoder, val, postForm)
}

func (f NEWSAPITopHeadlines) String() string {
	sb := strings.Builder{}
	sb.WriteString("NEWSAPISources:\n")
	sb.WriteString(fmt.Sprintf("\t- Keyword : %s\n", f.Keyword))
	sb.WriteString(fmt.Sprintf("\t- Sources : %s\n", f.Sources))
	sb.WriteString(fmt.Sprintf("\t- Category: %s\n", f.Category))
	sb.WriteString(fmt.Sprintf("\t- Country : %s\n", f.Country))
	return sb.String()
}

func (f NEWSAPITopHeadlines) Key() pageform.PageFormRepoKey {
	return pageform.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f NEWSAPITopHeadlines) SelectionOpts() []object.SelectOpts {
	return SelectionOpts()
}
