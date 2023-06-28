package pageform

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
)

const (
	NEWSAPILangAll        string = "all"
	NEWSAPILangArabic     string = "ar"
	NEWSAPILangGerman     string = "de"
	NEWSAPILangEnglish    string = "en"
	NEWSAPILangSpanish    string = "es"
	NEWSAPILangFrench     string = "fr"
	NEWSAPILangHebrew     string = "he"
	NEWSAPILangItalian    string = "it"
	NEWSAPILangDutch      string = "nl"
	NEWSAPILangNorwegian  string = "no"
	NEWSAPILangPortuguese string = "pt"
	NEWSAPILangRussian    string = "ru"
	NEWSAPILangSwedish    string = "sv"
	NEWSAPILangUrdu       string = "ud" // the official language of Pakistan
	NEWSAPILangChinese    string = "zh"
)

var NEWSAPILangVal = validator.NewEnmus(
	"newsapi_lang",
	NEWSAPILangAll,
	NEWSAPILangArabic,
	NEWSAPILangGerman,
	NEWSAPILangEnglish,
	NEWSAPILangSpanish,
	NEWSAPILangFrench,
	NEWSAPILangHebrew,
	NEWSAPILangItalian,
	NEWSAPILangDutch,
	NEWSAPILangNorwegian,
	NEWSAPILangPortuguese,
	NEWSAPILangRussian,
	NEWSAPILangSwedish,
	NEWSAPILangUrdu,
	NEWSAPILangChinese,
)

const (
	NEWSAPICtryAll                string = "all"
	NEWSAPICtryUnitedArabEmirates string = "ae"
	NEWSAPICtryArgentina          string = "ar"
	NEWSAPICtryAustria            string = "at"
	NEWSAPICtryAustralia          string = "au"
	NEWSAPICtryBelgium            string = "be"
	NEWSAPICtryBulgaria           string = "bg"
	NEWSAPICtryBrazil             string = "br"
	NEWSAPICtryCanada             string = "ca"
	NEWSAPICtrySwitzerland        string = "ch"
	NEWSAPICtryChina              string = "cn"
	NEWSAPICtryColombia           string = "co"
	NEWSAPICtryCuba               string = "cu"
	NEWSAPICtryCzechia            string = "cz"
	NEWSAPICtryGermany            string = "de"
	NEWSAPICtryEgypt              string = "eg"
	NEWSAPICtryFrance             string = "fr"
	NEWSAPICtryUnitedKingdom      string = "gb"
	NEWSAPICtryGreece             string = "gr"
	NEWSAPICtryHongKong           string = "hk"
	NEWSAPICtryHungary            string = "hu"
	NEWSAPICtryIndonesia          string = "id"
	NEWSAPICtryIreland            string = "ie"
	NEWSAPICtryIsrael             string = "il"
	NEWSAPICtryIndia              string = "in"
	NEWSAPICtryItaly              string = "it"
	NEWSAPICtryJapan              string = "jp"
	NEWSAPICtryKorea              string = "kr"
	NEWSAPICtryLithuania          string = "lt"
	NEWSAPICtryLatvia             string = "lv"
	NEWSAPICtryMorocco            string = "ma"
	NEWSAPICtryMexico             string = "mx"
	NEWSAPICtryMalaysia           string = "my"
	NEWSAPICtryNigeria            string = "ng"
	NEWSAPICtryNetherlands        string = "nl"
	NEWSAPICtryNorway             string = "no"
	NEWSAPICtryNewZealand         string = "nz"
	NEWSAPICtryPhilippines        string = "ph"
	NEWSAPICtryPoland             string = "pl"
	NEWSAPICtryPortugal           string = "pt"
	NEWSAPICtryRomania            string = "ro"
	NEWSAPICtrySerbia             string = "rs"
	NEWSAPICtryRussian            string = "ru"
	NEWSAPICtrySaudiArabia        string = "sa"
	NEWSAPICtrySweden             string = "se"
	NEWSAPICtrySingapore          string = "sg"
	NEWSAPICtrySlovenia           string = "si"
	NEWSAPICtrySlovakia           string = "sk"
	NEWSAPICtryThailand           string = "th"
	NEWSAPICtryTurkey             string = "tr"
	NEWSAPICtryTaiwan             string = "tw"
	NEWSAPICtryUkraine            string = "ua"
	NEWSAPICtryUnitedStates       string = "us"
	NEWSAPICtryVenezuela          string = "ve"
	NEWSAPICtrySouthAfrica        string = "za"
)

var NEWSAPICtryVal = validator.NewEnmus(
	"newsapi_ctry",
	NEWSAPICtryAll,
	NEWSAPICtryUnitedArabEmirates,
	NEWSAPICtryArgentina,
	NEWSAPICtryAustria,
	NEWSAPICtryAustralia,
	NEWSAPICtryBelgium,
	NEWSAPICtryBulgaria,
	NEWSAPICtryBrazil,
	NEWSAPICtryCanada,
	NEWSAPICtrySwitzerland,
	NEWSAPICtryChina,
	NEWSAPICtryColombia,
	NEWSAPICtryCuba,
	NEWSAPICtryCzechia,
	NEWSAPICtryGermany,
	NEWSAPICtryEgypt,
	NEWSAPICtryFrance,
	NEWSAPICtryUnitedKingdom,
	NEWSAPICtryGreece,
	NEWSAPICtryHongKong,
	NEWSAPICtryHungary,
	NEWSAPICtryIndonesia,
	NEWSAPICtryIreland,
	NEWSAPICtryIsrael,
	NEWSAPICtryIndia,
	NEWSAPICtryItaly,
	NEWSAPICtryJapan,
	NEWSAPICtryKorea,
	NEWSAPICtryLithuania,
	NEWSAPICtryLatvia,
	NEWSAPICtryMorocco,
	NEWSAPICtryMexico,
	NEWSAPICtryMalaysia,
	NEWSAPICtryNigeria,
	NEWSAPICtryNetherlands,
	NEWSAPICtryNorway,
	NEWSAPICtryNewZealand,
	NEWSAPICtryPhilippines,
	NEWSAPICtryPoland,
	NEWSAPICtryPortugal,
	NEWSAPICtryRomania,
	NEWSAPICtrySerbia,
	NEWSAPICtryRussian,
	NEWSAPICtrySaudiArabia,
	NEWSAPICtrySweden,
	NEWSAPICtrySingapore,
	NEWSAPICtrySlovenia,
	NEWSAPICtrySlovakia,
	NEWSAPICtryThailand,
	NEWSAPICtryTurkey,
	NEWSAPICtryTaiwan,
	NEWSAPICtryUkraine,
	NEWSAPICtryUnitedStates,
	NEWSAPICtryVenezuela,
	NEWSAPICtrySouthAfrica,
)

const (
	NEWSAPICatAll           string = "all"
	NEWSAPICatBusiness      string = "business"
	NEWSAPICatEntertainment string = "entertainment"
	NEWSAPICatGeneral       string = "general"
	NEWSAPICatHealth        string = "health"
	NEWSAPICatScience       string = "science"
	NEWSAPICatSports        string = "sports"
	NEWSAPICatTechnology    string = "technology"
)

var NEWSAPICatVal = validator.NewEnmus(
	"newsapi_cat",
	NEWSAPICatAll,
	NEWSAPICatBusiness,
	NEWSAPICatEntertainment,
	NEWSAPICatGeneral,
	NEWSAPICatHealth,
	NEWSAPICatScience,
	NEWSAPICatSports,
	NEWSAPICatTechnology,
)

type NEWSAPIEverything struct {
	SearchIn
	TimeRange
	Keyword        string `form:"keyword"`
	Sources        string `form:"sources"`
	Domains        string `form:"domains"`
	ExcludeDomains string `form:"exclude-domains"`
	Language       string `form:"language" val:"newsapi_lang"`
}

type NEWSAPISources struct {
	Language string `form:"language" val:"newsapi_lang"`
	Country  string `form:"country"  val:"newsapi_ctry"`
	Category string `form:"category" val:"newsapi_cat"`
}

type NEWSAPITopHeadlines struct {
	Keyword  string `form:"keyword"`
	Sources  string `form:"sources"`
	Category string `form:"category"`
}
