package pageform

import (
	"fmt"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
)

type GNewsCategory string

const (
	GNewsCatGeneral       GNewsCategory = "general"
	GNewsCatWorld         GNewsCategory = "world"
	GNewsCatNation        GNewsCategory = "nation"
	GNewsCatBusiness      GNewsCategory = "business"
	GNewsCatTechnology    GNewsCategory = "technology"
	GNewsCatEntertainment GNewsCategory = "entertainment"
	GNewsCatSports        GNewsCategory = "sports"
	GNewsCatScience       GNewsCategory = "science"
	GNewsCatHealth        GNewsCategory = "health"
)

var GnewsCatVal = validator.NewEnmus(
	"gnews_cat",
	GNewsCatGeneral,
	GNewsCatWorld,
	GNewsCatNation,
	GNewsCatBusiness,
	GNewsCatTechnology,
	GNewsCatEntertainment,
	GNewsCatSports,
	GNewsCatScience,
	GNewsCatHealth,
)

func (e *GNewsCategory) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = GNewsCategory(s)
	case string:
		*e = GNewsCategory(s)
	default:
		return fmt.Errorf("unsupported scan type for ApiType: %T", src)
	}
	return nil
}

type GNewsLanguage string

const (
	GNewsLangArabic     GNewsLanguage = "ar"
	GNewsLangChinese    GNewsLanguage = "zh"
	GNewsLangDutch      GNewsLanguage = "nl"
	GNewsLangEnglish    GNewsLanguage = "en"
	GNewsLangFrench     GNewsLanguage = "fr"
	GNewsLangGerman     GNewsLanguage = "de"
	GNewsLangGreek      GNewsLanguage = "el"
	GNewsLangHebrew     GNewsLanguage = "he"
	GNewsLangHindi      GNewsLanguage = "hi"
	GNewsLangItalian    GNewsLanguage = "it"
	GNewsLangJapanese   GNewsLanguage = "ja"
	GNewsLangMalayalam  GNewsLanguage = "ml"
	GNewsLangMarathi    GNewsLanguage = "mr"
	GNewsLangNorwegian  GNewsLanguage = "no"
	GNewsLangPortuguese GNewsLanguage = "pt"
	GNewsLangRomanian   GNewsLanguage = "ro"
	GNewsLangRussian    GNewsLanguage = "ru"
	GNewsLangSpanish    GNewsLanguage = "es"
	GNewsLangSwedish    GNewsLanguage = "sv"
	GNewsLangTamil      GNewsLanguage = "ta"
	GNewsLangTelugu     GNewsLanguage = "te"
	GNewsLangUkrainian  GNewsLanguage = "uk"
)

var GnewsLangVal = validator.NewEnmus(
	"gnews_lang",
	GNewsLangArabic,
	GNewsLangChinese,
	GNewsLangDutch,
	GNewsLangEnglish,
	GNewsLangFrench,
	GNewsLangGerman,
	GNewsLangGreek,
	GNewsLangHebrew,
	GNewsLangHindi,
	GNewsLangItalian,
	GNewsLangJapanese,
	GNewsLangMalayalam,
	GNewsLangMarathi,
	GNewsLangNorwegian,
	GNewsLangPortuguese,
	GNewsLangRomanian,
	GNewsLangRussian,
	GNewsLangSpanish,
	GNewsLangSwedish,
	GNewsLangTamil,
	GNewsLangTelugu,
	GNewsLangUkrainian,
)

func (e *GNewsLanguage) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = GNewsLanguage(s)
	case string:
		*e = GNewsLanguage(s)
	default:
		return fmt.Errorf("unsupported scan type for ApiType: %T", src)
	}
	return nil
}

type GNewsCountry string

const (
	GNewsCtryAustralia         GNewsCountry = "au"
	GNewsCtryBrazil            GNewsCountry = "br"
	GNewsCtryCanada            GNewsCountry = "ca"
	GNewsCtryChina             GNewsCountry = "cn"
	GNewsCtryEgypt             GNewsCountry = "eg"
	GNewsCtryFrance            GNewsCountry = "fr"
	GNewsCtryGermany           GNewsCountry = "de"
	GNewsCtryGreece            GNewsCountry = "gr"
	GNewsCtryHongKong          GNewsCountry = "hk"
	GNewsCtryIndia             GNewsCountry = "in"
	GNewsCtryIreland           GNewsCountry = "ie"
	GNewsCtryIsrael            GNewsCountry = "il"
	GNewsCtryItaly             GNewsCountry = "it"
	GNewsCtryJapan             GNewsCountry = "jp"
	GNewsCtryNetherlands       GNewsCountry = "nl"
	GNewsCtryNorway            GNewsCountry = "no"
	GNewsCtryPakistan          GNewsCountry = "pk"
	GNewsCtryPeru              GNewsCountry = "pe"
	GNewsCtryPhilippines       GNewsCountry = "ph"
	GNewsCtryPortugal          GNewsCountry = "pt"
	GNewsCtryRomania           GNewsCountry = "ro"
	GNewsCtryRussianFederation GNewsCountry = "ru"
	GNewsCtrySingapore         GNewsCountry = "sg"
	GNewsCtrySpain             GNewsCountry = "es"
	GNewsCtrySweden            GNewsCountry = "se"
	GNewsCtrySwitzerland       GNewsCountry = "ch"
	GNewsCtryTaiwan            GNewsCountry = "tw"
	GNewsCtryUkraine           GNewsCountry = "ua"
	GNewsCtryUnitedKingdom     GNewsCountry = "gb"
	GNewsCtryUnitedStates      GNewsCountry = "us"
)

var GNewsCtryVal = validator.NewEnmus(
	"gnews_ctry",
	GNewsCtryAustralia,
	GNewsCtryBrazil,
	GNewsCtryCanada,
	GNewsCtryChina,
	GNewsCtryEgypt,
	GNewsCtryFrance,
	GNewsCtryGermany,
	GNewsCtryGreece,
	GNewsCtryHongKong,
	GNewsCtryIndia,
	GNewsCtryIreland,
	GNewsCtryIsrael,
	GNewsCtryItaly,
	GNewsCtryJapan,
	GNewsCtryNetherlands,
	GNewsCtryNorway,
	GNewsCtryPakistan,
	GNewsCtryPeru,
	GNewsCtryPhilippines,
	GNewsCtryPortugal,
	GNewsCtryRomania,
	GNewsCtryRussianFederation,
	GNewsCtrySingapore,
	GNewsCtrySpain,
	GNewsCtrySweden,
	GNewsCtrySwitzerland,
	GNewsCtryTaiwan,
	GNewsCtryUkraine,
	GNewsCtryUnitedKingdom,
	GNewsCtryUnitedStates,
)

func (e *GNewsCountry) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = GNewsCountry(s)
	case string:
		*e = GNewsCountry(s)
	default:
		return fmt.Errorf("unsupported scan type for ApiType: %T", src)
	}
	return nil
}

type GNewsHeadlines struct {
	TimeRange
	Keyword  string        `form:"keyword"`
	Language GNewsLanguage `form:"language" val:"gnews_lang"`
	Country  GNewsCountry  `form:"country"  val:"gnews_ctry"`
	Category GNewsCategory `form:"category" val:"gnews_cat"`
}

type GNewsSearch struct {
	SearchIn
	TimeRange
	Keyword  string        `form:"keyword"`
	Language GNewsLanguage `form:"language" val:"gnews_lang"`
	Country  GNewsCountry  `form:"country"  val:"gnews_ctry"`
}
