package newsdata

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
	EPLatestNews  string = "Latest News"
	EPNewsArchive string = "News Archive"
	EPNewsSources string = "News Sources"
	EPCrypto      string = "Crypto"
)

type NEWSDATAIOLatestNews struct {
	// IncludeContent // currently not yet support
	Keyword  string   `form:"keyword"  validate:"max=512"`
	Domains  string   `form:"domains"  validate:"max=5"`
	Language []string `form:"language" validate:"max=5,newsdata_lang"`
	Country  []string `form:"country"  validate:"max=5,newsdata_ctry"`
	Category []string `form:"category" validate:"max=5,newsdata_cat"`
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

type NEWSDATAIONewsArchive struct {
	pageform.TimeRange
	Keyword  string   `form:"keyword"  validate:"max=512"`
	Domains  string   `form:"domains"  validate:"max=512"`
	Language []string `form:"language" validate:"max=5,newsdata_lang"`
	Country  []string `form:"country"  validate:"max=5,newsdata_ctry"`
	Category []string `form:"category" validate:"max=5,newsdata_cat"`
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

type IncludeContent struct {
	Image       bool `form:"image"`
	Video       bool `form:"video"`
	FullContent bool `form:"full-content"`
}
