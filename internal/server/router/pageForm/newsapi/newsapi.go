package newsapi

import (
	"fmt"
	"net/url"
	"strings"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
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

const NewsAPIName = "NEWS API"

const (
	EPEverything   string = "Everything"
	EPTopHeadlines        = "Top Headlines"
	EPSources             = "Sources"
)

type NEWSAPIEverything struct {
	pageform.SearchIn
	pageform.TimeRange
	Keyword        string `form:"keyword"`
	Sources        string `form:"sources"`
	Domains        string `form:"domains"`
	ExcludeDomains string `form:"exclude-domains"`
	Language       string `form:"language" val:"newsapi_lang"`
}

func (f NEWSAPIEverything) Endpoint() string {
	return EPEverything
}

func (f NEWSAPIEverything) API() string {
	return NewsAPIName
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

type NEWSAPISources struct {
	Language string `form:"language" val:"newsapi_lang"`
	Country  string `form:"country"  val:"newsapi_ctry"`
	Category string `form:"category" val:"newsapi_cat"`
}

func (f NEWSAPISources) Endpoint() string {
	return EPSources
}

func (f NEWSAPISources) API() string {
	return NewsAPIName
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

type NEWSAPITopHeadlines struct {
	Keyword  string `form:"keyword"`
	Sources  string `form:"sources"`
	Country  string `form:"country"`
	Category string `form:"category"`
}

func (f NEWSAPITopHeadlines) Endpoint() string {
	return EPTopHeadlines
}

func (f NEWSAPITopHeadlines) API() string {
	return NewsAPIName
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

const (
	Business      string = "business"
	Entertainment string = "entertainment"
	General       string = "general"
	Health        string = "health"
	Science       string = "science"
	Sports        string = "sports"
	Technology    string = "technology"
)

var Category = map[string]string{
	Business:      "Business",
	Entertainment: "Entertainment",
	General:       "General",
	Health:        "Health",
	Science:       "Science",
	Sports:        "Sports",
	Technology:    "Technology",
}

var CategoryValidator = validator.NewEnmusFromMap(
	"newsapi_cat", Category, "key",
)

const (
	UnitedArabEmirates string = "ae"
	Argentina          string = "ar"
	Austria            string = "at"
	Australia          string = "au"
	Belgium            string = "be"
	Bulgaria           string = "bg"
	Brazil             string = "br"
	Canada             string = "ca"
	Switzerland        string = "ch"
	China              string = "cn"
	Colombia           string = "co"
	Cuba               string = "cu"
	Czechia            string = "cz"
	Germany            string = "de"
	Egypt              string = "eg"
	France             string = "fr"
	UnitedKingdom      string = "gb"
	Greece             string = "gr"
	HongKong           string = "hk"
	Hungary            string = "hu"
	Indonesia          string = "id"
	Ireland            string = "ie"
	Israel             string = "il"
	India              string = "in"
	Italy              string = "it"
	Japan              string = "jp"
	Korea              string = "kr"
	Lithuania          string = "lt"
	Latvia             string = "lv"
	Morocco            string = "ma"
	Mexico             string = "mx"
	Malaysia           string = "my"
	Nigeria            string = "ng"
	Netherlands        string = "nl"
	Norway             string = "no"
	NewZealand         string = "nz"
	Philippines        string = "ph"
	Poland             string = "pl"
	Portugal           string = "pt"
	Romania            string = "ro"
	Serbia             string = "rs"
	Russia             string = "ru"
	SaudiArabia        string = "sa"
	Sweden             string = "se"
	Singapore          string = "sg"
	Slovenia           string = "si"
	Slovakia           string = "sk"
	Thailand           string = "th"
	Turkey             string = "tr"
	Taiwan             string = "tw"
	Ukraine            string = "ua"
	UnitedStates       string = "us"
	Venezuela          string = "ve"
	SouthAfrica        string = "za"
)

var Country = map[string]string{
	UnitedArabEmirates: "UnitedArabEmirates",
	Argentina:          "Argentina",
	Austria:            "Austria",
	Australia:          "Australia",
	Belgium:            "Belgium",
	Bulgaria:           "Bulgaria",
	Brazil:             "Brazil",
	Canada:             "Canada",
	Switzerland:        "Switzerland",
	China:              "China",
	Colombia:           "Colombia",
	Cuba:               "Cuba",
	Czechia:            "Czechia",
	Germany:            "Germany",
	Egypt:              "Egypt",
	France:             "France",
	UnitedKingdom:      "UnitedKingdom",
	Greece:             "Greece",
	HongKong:           "HongKong",
	Hungary:            "Hungary",
	Indonesia:          "Indonesia",
	Ireland:            "Ireland",
	Israel:             "Israel",
	India:              "India",
	Italy:              "Italy",
	Japan:              "Japan",
	Korea:              "Korea",
	Lithuania:          "Lithuania",
	Latvia:             "Latvia",
	Morocco:            "Morocco",
	Mexico:             "Mexico",
	Malaysia:           "Malaysia",
	Nigeria:            "Nigeria",
	Netherlands:        "Netherlands",
	Norway:             "Norway",
	NewZealand:         "New Zealand",
	Philippines:        "Philippines",
	Poland:             "Poland",
	Portugal:           "Portugal",
	Romania:            "Romania",
	Serbia:             "Serbia",
	Russia:             "Russia",
	SaudiArabia:        "SaudiArabia",
	Sweden:             "Sweden",
	Singapore:          "Singapore",
	Slovenia:           "Slovenia",
	Slovakia:           "Slovakia",
	Thailand:           "Thailand",
	Turkey:             "Turkey",
	Taiwan:             "Taiwan",
	Ukraine:            "Ukraine",
	UnitedStates:       "UnitedStates",
	Venezuela:          "Venezuela",
	SouthAfrica:        "SouthAfrica",
}

var CountryValidator = validator.NewEnmusFromMap(
	"newsapi_ctry", Country, "key")

const (
	Arabic     string = "ar"
	German     string = "de"
	English    string = "en"
	Spanish    string = "es"
	French     string = "fr"
	Hebrew     string = "he"
	Italian    string = "it"
	Dutch      string = "nl"
	Norwegian  string = "no"
	Portuguese string = "pt"
	Russian    string = "ru"
	Swedish    string = "sv"
	Urdu       string = "ud" // the official language of Pakistan
	Chinese    string = "zh"
)

var Language = map[string]string{
	Arabic:     "Arabic",
	German:     "German",
	English:    "English",
	Spanish:    "Spanish",
	French:     "French",
	Hebrew:     "Hebrew",
	Italian:    "Italian",
	Dutch:      "Dutch",
	Norwegian:  "Norwegian",
	Portuguese: "Portuguese",
	Russian:    "Russian",
	Swedish:    "Swedish",
	Urdu:       "Urdu",
	Chinese:    "Chinese",
}

var LanguageValidator = validator.NewEnmusFromMap(
	"newsapi_lang", Language, "key",
)
