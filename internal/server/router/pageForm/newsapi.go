package pageform

import (
	"fmt"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
)

type NEWSAPILanguage string

const (
	NEWSAPILangAll        NEWSAPILanguage = "all"
	NEWSAPILangArabic     NEWSAPILanguage = "ar"
	NEWSAPILangGerman     NEWSAPILanguage = "de"
	NEWSAPILangEnglish    NEWSAPILanguage = "en"
	NEWSAPILangSpanish    NEWSAPILanguage = "es"
	NEWSAPILangFrench     NEWSAPILanguage = "fr"
	NEWSAPILangHebrew     NEWSAPILanguage = "he"
	NEWSAPILangItalian    NEWSAPILanguage = "it"
	NEWSAPILangDutch      NEWSAPILanguage = "nl"
	NEWSAPILangNorwegian  NEWSAPILanguage = "no"
	NEWSAPILangPortuguese NEWSAPILanguage = "pt"
	NEWSAPILangRussian    NEWSAPILanguage = "ru"
	NEWSAPILangSwedish    NEWSAPILanguage = "sv"
	NEWSAPILangUrdu       NEWSAPILanguage = "ud" // the official language of Pakistan
	NEWSAPILangChinese    NEWSAPILanguage = "zh"
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

func (e *NEWSAPILanguage) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NEWSAPILanguage(s)
	case string:
		*e = NEWSAPILanguage(s)
	default:
		return fmt.Errorf("unsupported scan type for ApiType: %T", src)
	}
	return nil
}

type NEWSAPICountry string

const (
	NEWSAPICtryAll                NEWSAPICountry = "all"
	NEWSAPICtryUnitedArabEmirates NEWSAPICountry = "ae"
	NEWSAPICtryArgentina          NEWSAPICountry = "ar"
	NEWSAPICtryAustria            NEWSAPICountry = "at"
	NEWSAPICtryAustralia          NEWSAPICountry = "au"
	NEWSAPICtryBelgium            NEWSAPICountry = "be"
	NEWSAPICtryBulgaria           NEWSAPICountry = "bg"
	NEWSAPICtryBrazil             NEWSAPICountry = "br"
	NEWSAPICtryCanada             NEWSAPICountry = "ca"
	NEWSAPICtrySwitzerland        NEWSAPICountry = "ch"
	NEWSAPICtryChina              NEWSAPICountry = "cn"
	NEWSAPICtryColombia           NEWSAPICountry = "co"
	NEWSAPICtryCuba               NEWSAPICountry = "cu"
	NEWSAPICtryCzechia            NEWSAPICountry = "cz"
	NEWSAPICtryGermany            NEWSAPICountry = "de"
	NEWSAPICtryEgypt              NEWSAPICountry = "eg"
	NEWSAPICtryFrance             NEWSAPICountry = "fr"
	NEWSAPICtryUnitedKingdom      NEWSAPICountry = "gb"
	NEWSAPICtryGreece             NEWSAPICountry = "gr"
	NEWSAPICtryHongKong           NEWSAPICountry = "hk"
	NEWSAPICtryHungary            NEWSAPICountry = "hu"
	NEWSAPICtryIndonesia          NEWSAPICountry = "id"
	NEWSAPICtryIreland            NEWSAPICountry = "ie"
	NEWSAPICtryIsrael             NEWSAPICountry = "il"
	NEWSAPICtryIndia              NEWSAPICountry = "in"
	NEWSAPICtryItaly              NEWSAPICountry = "it"
	NEWSAPICtryJapan              NEWSAPICountry = "jp"
	NEWSAPICtryKorea              NEWSAPICountry = "kr"
	NEWSAPICtryLithuania          NEWSAPICountry = "lt"
	NEWSAPICtryLatvia             NEWSAPICountry = "lv"
	NEWSAPICtryMorocco            NEWSAPICountry = "ma"
	NEWSAPICtryMexico             NEWSAPICountry = "mx"
	NEWSAPICtryMalaysia           NEWSAPICountry = "my"
	NEWSAPICtryNigeria            NEWSAPICountry = "ng"
	NEWSAPICtryNetherlands        NEWSAPICountry = "nl"
	NEWSAPICtryNorway             NEWSAPICountry = "no"
	NEWSAPICtryNewZealand         NEWSAPICountry = "nz"
	NEWSAPICtryPhilippines        NEWSAPICountry = "ph"
	NEWSAPICtryPoland             NEWSAPICountry = "pl"
	NEWSAPICtryPortugal           NEWSAPICountry = "pt"
	NEWSAPICtryRomania            NEWSAPICountry = "ro"
	NEWSAPICtrySerbia             NEWSAPICountry = "rs"
	NEWSAPICtryRussian            NEWSAPICountry = "ru"
	NEWSAPICtrySaudiArabia        NEWSAPICountry = "sa"
	NEWSAPICtrySweden             NEWSAPICountry = "se"
	NEWSAPICtrySingapore          NEWSAPICountry = "sg"
	NEWSAPICtrySlovenia           NEWSAPICountry = "si"
	NEWSAPICtrySlovakia           NEWSAPICountry = "sk"
	NEWSAPICtryThailand           NEWSAPICountry = "th"
	NEWSAPICtryTurkey             NEWSAPICountry = "tr"
	NEWSAPICtryTaiwan             NEWSAPICountry = "tw"
	NEWSAPICtryUkraine            NEWSAPICountry = "ua"
	NEWSAPICtryUnitedStates       NEWSAPICountry = "us"
	NEWSAPICtryVenezuela          NEWSAPICountry = "ve"
	NEWSAPICtrySouthAfrica        NEWSAPICountry = "za"
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

func (e *NEWSAPICountry) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NEWSAPICountry(s)
	case string:
		*e = NEWSAPICountry(s)
	default:
		return fmt.Errorf("unsupported scan type for ApiType: %T", src)
	}
	return nil
}

type NEWSAPICategory string

const (
	NEWSAPICatAll           NEWSAPICategory = "all"
	NEWSAPICatBusiness      NEWSAPICategory = "business"
	NEWSAPICatEntertainment NEWSAPICategory = "entertainment"
	NEWSAPICatGeneral       NEWSAPICategory = "general"
	NEWSAPICatHealth        NEWSAPICategory = "health"
	NEWSAPICatScience       NEWSAPICategory = "science"
	NEWSAPICatSports        NEWSAPICategory = "sports"
	NEWSAPICatTechnology    NEWSAPICategory = "technology"
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

func (e *NEWSAPICategory) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NEWSAPICategory(s)
	case string:
		*e = NEWSAPICategory(s)
	default:
		return fmt.Errorf("unsupported scan type for ApiType: %T", src)
	}
	return nil
}

type NEWSAPIEverything struct {
	SearchIn
	TimeRange
	Keyword        string          `form:"keyword"`
	Sources        string          `form:"sources"`
	Domains        string          `form:"domains"`
	ExcludeDomains string          `form:"exclude-domains"`
	Language       NEWSAPILanguage `form:"language" val:"newsapi_lang"`
}

type NEWSAPISources struct {
	Language NEWSAPILanguage `form:"language" val:"newsapi_lang"`
	Country  NEWSAPICountry  `form:"country"  val:"newsapi_ctry"`
	Category NEWSAPICategory `form:"category" val:"newsapi_cat"`
}

type NEWSAPITopHeadlines struct {
	Keyword  string          `form:"keyword"`
	Sources  string          `form:"sources"`
	Category NEWSAPICategory `form:"category"`
}
