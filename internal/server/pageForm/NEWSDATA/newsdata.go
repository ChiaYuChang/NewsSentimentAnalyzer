package newsdata

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
		NEWSDATAIOLatestNews{},
		NEWSDATAIONewsArchive{},
		NEWSDATAIONewsSources{},
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

const API_NAME = "NEWSDATA.IO"

const (
	VAL_TAG_DOMAIN   = "newsdata_domain"
	VAL_TAG_CATEGORY = "newsdata_cat"
	VAL_TAG_LANGUAGE = "newsdata_lang"
	VAL_TAG_COUNTRY  = "newsdata_ctry"
)

const (
	EPLatestNews  = "Latest News"
	EPNewsArchive = "News Archive"
	EPNewsSources = "News Sources"
	EPCrypto      = "Crypto"
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
		{
			OptMap:         Country,
			MaxDiv:         5,
			DefaultValue:   "",
			DefaultText:    "all",
			InsertButtonId: "insert-country-btn",
			DeleteButtonId: "delete-country-btn",
			PositionId:     "country",
			AlertMessage:   "You can only add up to 5 categories in a single query",
		},
	}
}

type NEWSDATAIOLatestNews struct {
	// IncludeContent // currently not yet support
	Keyword  string   `mod:"trim"  form:"keyword"  validate:"max=512"`
	Domains  string   `mod:"trim"  form:"domains"  validate:"max=5"`
	Language []string `            form:"language" validate:"max=5,newsdata_lang"`
	Country  []string `            form:"country"  validate:"max=5,newsdata_ctry"`
	Category []string `            form:"category" validate:"max=5,newsdata_cat"`
}

func (f NEWSDATAIOLatestNews) Endpoint() string {
	return EPLatestNews
}

func (f NEWSDATAIOLatestNews) API() string {
	return API_NAME
}

func (f NEWSDATAIOLatestNews) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[NEWSDATAIOLatestNews](decoder, val, postForm)
}

func (f NEWSDATAIOLatestNews) String() string {
	sb := strings.Builder{}
	sb.WriteString("NEWSDATAIOLatestNews:\n")
	sb.WriteString(fmt.Sprintf("\t- Keyword : %s\n", f.Keyword))
	sb.WriteString(fmt.Sprintf("\t- Domains : %s\n", f.Domains))
	sb.WriteString(fmt.Sprintf("\t- Category: %s\n", strings.Join(f.Category, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Country : %s\n", strings.Join(f.Country, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Language: %s\n", strings.Join(f.Language, ", ")))
	return sb.String()
}

func (f NEWSDATAIOLatestNews) Key() pageform.PageFormRepoKey {
	return pageform.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f NEWSDATAIOLatestNews) SelectionOpts() []object.SelectOpts {
	return SelectionOpts()
}

type NEWSDATAIONewsArchive struct {
	pageform.TimeRange
	Keyword  string   `mod:"trim"  form:"keyword"  validate:"max=512"`
	Domains  string   `mod:"trim"  form:"domains"  validate:"max=512"`
	Language []string `            form:"language" validate:"max=5,newsdata_lang"`
	Country  []string `            form:"country"  validate:"max=5,newsdata_ctry"`
	Category []string `            form:"category" validate:"max=5,newsdata_cat"`
}

func (f NEWSDATAIONewsArchive) Endpoint() string {
	return EPNewsArchive
}

func (f NEWSDATAIONewsArchive) API() string {
	return API_NAME
}

func (f NEWSDATAIONewsArchive) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[NEWSDATAIONewsArchive](decoder, val, postForm)
}

func (f NEWSDATAIONewsArchive) String() string {
	sb := strings.Builder{}
	sb.WriteString("EWSDATAIONewsArchive:\n")
	sb.WriteString(fmt.Sprintf("\t- Keyword : %s\n", f.Keyword))
	sb.WriteString(fmt.Sprintf("\t- Domains : %s\n", f.Domains))
	sb.WriteString(fmt.Sprintf("\t- Category: %s\n", strings.Join(f.Category, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Country : %s\n", strings.Join(f.Country, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Language: %s\n", strings.Join(f.Language, ", ")))
	sb.WriteString(f.TimeRange.ToString("\t"))
	return sb.String()
}

func (f NEWSDATAIONewsArchive) Key() pageform.PageFormRepoKey {
	return pageform.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f NEWSDATAIONewsArchive) SelectionOpts() []object.SelectOpts {
	return SelectionOpts()
}

type NEWSDATAIONewsSources struct {
	Language []string `form:"language" validate:"max=5,newsdata_lang"`
	Country  []string `form:"country"  validate:"max=5,newsdata_ctry"`
	Category []string `form:"category" validate:"max=5,newsdata_cat"`
}

func (f NEWSDATAIONewsSources) Endpoint() string {
	return EPNewsSources
}

func (f NEWSDATAIONewsSources) API() string {
	return API_NAME
}

func (f NEWSDATAIONewsSources) FormDecodeAndValidate(
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (pageform.PageForm, error) {
	return pageform.FormDecodeAndValidate[NEWSDATAIONewsSources](decoder, val, postForm)
}

func (f NEWSDATAIONewsSources) String() string {
	sb := strings.Builder{}
	sb.WriteString("EWSDATAIONewsArchive:\n")
	sb.WriteString(fmt.Sprintf("\t- Category: %s\n", strings.Join(f.Category, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Country : %s\n", strings.Join(f.Country, ", ")))
	sb.WriteString(fmt.Sprintf("\t- Language: %s\n", strings.Join(f.Language, ", ")))
	return sb.String()
}

func (f NEWSDATAIONewsSources) Key() pageform.PageFormRepoKey {
	return pageform.NewPageFormRepoKey(f.API(), f.Endpoint())
}

func (f NEWSDATAIONewsSources) SelectionOpts() []object.SelectOpts {
	return SelectionOpts()
}

type IncludeContent struct {
	Image       bool `form:"image"`
	Video       bool `form:"video"`
	FullContent bool `form:"full-content"`
}
