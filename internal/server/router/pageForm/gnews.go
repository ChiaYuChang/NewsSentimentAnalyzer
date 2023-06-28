package pageform

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
)

const (
	GNewsCatGeneral       = "general"
	GNewsCatWorld         = "world"
	GNewsCatNation        = "nation"
	GNewsCatBusiness      = "business"
	GNewsCatTechnology    = "technology"
	GNewsCatEntertainment = "entertainment"
	GNewsCatSports        = "sports"
	GNewsCatScience       = "science"
	GNewsCatHealth        = "health"
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

const (
	GNewsLangArabic     = "ar"
	GNewsLangChinese    = "zh"
	GNewsLangDutch      = "nl"
	GNewsLangEnglish    = "en"
	GNewsLangFrench     = "fr"
	GNewsLangGerman     = "de"
	GNewsLangGreek      = "el"
	GNewsLangHebrew     = "he"
	GNewsLangHindi      = "hi"
	GNewsLangItalian    = "it"
	GNewsLangJapanese   = "ja"
	GNewsLangMalayalam  = "ml"
	GNewsLangMarathi    = "mr"
	GNewsLangNorwegian  = "no"
	GNewsLangPortuguese = "pt"
	GNewsLangRomanian   = "ro"
	GNewsLangRussian    = "ru"
	GNewsLangSpanish    = "es"
	GNewsLangSwedish    = "sv"
	GNewsLangTamil      = "ta"
	GNewsLangTelugu     = "te"
	GNewsLangUkrainian  = "uk"
)

var GNewsLangVal = validator.NewEnmus(
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

const (
	GNewsCtryAustralia         = "au"
	GNewsCtryBrazil            = "br"
	GNewsCtryCanada            = "ca"
	GNewsCtryChina             = "cn"
	GNewsCtryEgypt             = "eg"
	GNewsCtryFrance            = "fr"
	GNewsCtryGermany           = "de"
	GNewsCtryGreece            = "gr"
	GNewsCtryHongKong          = "hk"
	GNewsCtryIndia             = "in"
	GNewsCtryIreland           = "ie"
	GNewsCtryIsrael            = "il"
	GNewsCtryItaly             = "it"
	GNewsCtryJapan             = "jp"
	GNewsCtryNetherlands       = "nl"
	GNewsCtryNorway            = "no"
	GNewsCtryPakistan          = "pk"
	GNewsCtryPeru              = "pe"
	GNewsCtryPhilippines       = "ph"
	GNewsCtryPortugal          = "pt"
	GNewsCtryRomania           = "ro"
	GNewsCtryRussianFederation = "ru"
	GNewsCtrySingapore         = "sg"
	GNewsCtrySpain             = "es"
	GNewsCtrySweden            = "se"
	GNewsCtrySwitzerland       = "ch"
	GNewsCtryTaiwan            = "tw"
	GNewsCtryUkraine           = "ua"
	GNewsCtryUnitedKingdom     = "gb"
	GNewsCtryUnitedStates      = "us"
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

type GNewsHeadlines struct {
	TimeRange
	Keyword  string `form:"keyword"`
	Language string `form:"language" validate:"gnews_lang"`
	Country  string `form:"country"  validate:"gnews_ctry"`
	Category string `form:"category" validate:"gnews_cat"`
}

type GNewsSearch struct {
	SearchIn
	TimeRange
	Keyword  string `form:"keyword"`
	Language string `form:"language" validate:"gnews_lang"`
	Country  string `form:"country"  validate:"gnews_ctry"`
}
