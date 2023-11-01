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
	Keyword  string   `form:"keyword"`
	Language []string `form:"language" validate:"gnews_lang"`
	Country  []string `form:"country"  validate:"gnews_ctry"`
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

const (
	General       string = "general"
	World                = "world"
	Nation               = "nation"
	Business             = "business"
	Technology           = "technology"
	Entertainment        = "entertainment"
	Sports               = "sports"
	Science              = "science"
	Health               = "health"
)

var Category = map[string]string{
	General:       "General",
	World:         "World",
	Nation:        "Nation",
	Business:      "Business",
	Technology:    "Technology",
	Entertainment: "Entertainment",
	Sports:        "Sports",
	Science:       "Science",
	Health:        "Health",
}

var CategoryValidator = validator.NewEnmusFromMap(
	VAL_TAG_CATEGORY, Category, "key",
)

const (
	Australia         = "au"
	Brazil            = "br"
	Canada            = "ca"
	China             = "cn"
	Egypt             = "eg"
	France            = "fr"
	Germany           = "de"
	Greece            = "gr"
	HongKong          = "hk"
	India             = "in"
	Ireland           = "ie"
	Israel            = "il"
	Italy             = "it"
	Japan             = "jp"
	Netherlands       = "nl"
	Norway            = "no"
	Pakistan          = "pk"
	Peru              = "pe"
	Philippines       = "ph"
	Portugal          = "pt"
	Romania           = "ro"
	RussianFederation = "ru"
	Singapore         = "sg"
	Spain             = "es"
	Sweden            = "se"
	Switzerland       = "ch"
	Taiwan            = "tw"
	Ukraine           = "ua"
	UnitedKingdom     = "gb"
	UnitedStates      = "us"
)

var Country = map[string]string{
	Australia:         "Australia",
	Brazil:            "Brazil",
	Canada:            "Canada",
	China:             "China",
	Egypt:             "Egypt",
	France:            "France",
	Germany:           "Germany",
	Greece:            "Greece",
	HongKong:          "Hong Kong",
	India:             "India",
	Ireland:           "Ireland",
	Israel:            "Israel",
	Italy:             "Italy",
	Japan:             "Japan",
	Netherlands:       "Netherlands",
	Norway:            "Norway",
	Pakistan:          "Pakistan",
	Peru:              "Peru",
	Philippines:       "Philippines",
	Portugal:          "Portugal",
	Romania:           "Romania",
	RussianFederation: "Russian Federation",
	Singapore:         "Singapore",
	Spain:             "Spain",
	Sweden:            "Sweden",
	Switzerland:       "Switzerland",
	Taiwan:            "Taiwan",
	Ukraine:           "Ukraine",
	UnitedKingdom:     "United Kingdom",
	UnitedStates:      "United States",
}

var CountryValidator = validator.NewEnmusFromMap(
	VAL_TAG_COUNTRY, Country, "key",
)

const (
	Arabic     = "ar"
	Chinese    = "zh"
	Dutch      = "nl"
	English    = "en"
	French     = "fr"
	German     = "de"
	Greek      = "el"
	Hebrew     = "he"
	Hindi      = "hi"
	Italian    = "it"
	Japanese   = "ja"
	Malayalam  = "ml"
	Marathi    = "mr"
	Norwegian  = "no"
	Portuguese = "pt"
	Romanian   = "ro"
	Russian    = "ru"
	Spanish    = "es"
	Swedish    = "sv"
	Tamil      = "ta"
	Telugu     = "te"
	Ukrainian  = "uk"
)

var Language = map[string]string{
	Arabic:     "Arabic",
	Chinese:    "Chinese",
	Dutch:      "Dutch",
	English:    "English",
	French:     "French",
	German:     "German",
	Greek:      "Greek",
	Hebrew:     "Hebrew",
	Hindi:      "Hindi",
	Italian:    "Italian",
	Japanese:   "Japanese",
	Malayalam:  "Malayalam",
	Marathi:    "Marathi",
	Norwegian:  "Norwegian",
	Portuguese: "Portuguese",
	Romanian:   "Romanian",
	Russian:    "Russian",
	Spanish:    "Spanish",
	Swedish:    "Swedish",
	Tamil:      "Tamil",
	Telugu:     "Telugu",
	Ukrainian:  "Ukrainian",
}

var LanguageValidator = validator.NewEnmusFromMap(
	VAL_TAG_LANGUAGE, Language, "key",
)
