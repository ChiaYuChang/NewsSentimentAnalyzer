package newsdata

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
)

func init() {
	cat := make([]string, len(CatList))
	sort.Slice(CatList, func(i, j int) bool {
		if CatList[i][0] < CatList[j][0] {
			return true
		}
		return false
	})
	for i := 0; i < len(cat); i++ {
		cat[i] = CatList[i][1]
	}
	CatVal = validator.NewEnmus(
		VAL_TAG_CATEGORY,
		cat...,
	)

	ctry := make([]string, len(CtryList))
	sort.Slice(CtryList, func(i, j int) bool {
		if CtryList[i][0] < CtryList[j][0] {
			return true
		}
		return false
	})
	for i := 0; i < len(ctry); i++ {
		ctry[i] = CtryList[i][1]
	}
	CtryVal = validator.NewEnmus(
		VAL_TAG_COUNTRY,
		ctry...,
	)

	lang := make([]string, len(LangList))
	sort.Slice(LangList, func(i, j int) bool {
		if LangList[i][0] < LangList[j][0] {
			return true
		}
		return false
	})
	for i := 0; i < len(lang); i++ {
		lang[i] = LangList[i][1]
	}
	CtryVal = validator.NewEnmus(
		VAL_TAG_LANGUAGE,
		lang...,
	)
}

const (
	API_ROOT             = "https://newsdata.io/api"
	API_VERSION          = "1"
	API_METHOD           = http.MethodGet
	API_TIME_FORMAT      = "2006-01-02"
	API_MAX_NUM_DOMAIN   = 5
	API_MAX_NUM_COUNTRY  = 5
	API_MAX_NUM_CATEGORY = 5
	API_MAX_NUM_LANGUAGE = 5
)

var API_URL, _ = url.Parse(fmt.Sprintf("%s/%s", API_ROOT, API_VERSION))

const (
	VAL_TAG_DOMAIN   = "domain"
	VAL_TAG_CATEGORY = "cat"
	VAL_TAG_LANGUAGE = "lang"
	VAL_TAG_COUNTRY  = "ctry"
)

const (
	EPLatestNews  string = "news"
	EPCrypto      string = "crypto"
	EPNewsArchive string = "archive"
	EPNewsSources string = "sources"
)

type SelectOpts [2]string

var (
	CatBusiness      = SelectOpts{"Business", "business"}
	CatEntertainment = SelectOpts{"Entertainment", "entertainment"}
	CatEnvironment   = SelectOpts{"Environment", "environment"}
	CatFood          = SelectOpts{"Food", "food"}
	CatHealth        = SelectOpts{"Health", "health"}
	CatPolitics      = SelectOpts{"Politics", "politics"}
	CatScience       = SelectOpts{"Science", "science"}
	CatSports        = SelectOpts{"Sports", "sports"}
	CatTechnology    = SelectOpts{"Technology", "technology"}
	CatTop           = SelectOpts{"Top", "top"}
	CatTourism       = SelectOpts{"Tourism", "tourism"}
	CatWorld         = SelectOpts{"World", "world"}
)

var CatList = [][2]string{
	CatBusiness,
	CatEntertainment,
	CatEnvironment,
	CatFood,
	CatHealth,
	CatPolitics,
	CatScience,
	CatSports,
	CatTechnology,
	CatTop,
	CatTourism,
	CatWorld,
}

var CatVal validator.Enmus[string]

var (
	CtryAfghanistan          SelectOpts = SelectOpts{"Afghanistan", "af"}
	CtryAlbania              SelectOpts = SelectOpts{"Albania", "al"}
	CtryAlgeria              SelectOpts = SelectOpts{"Algeria", "dz"}
	CtryAngola               SelectOpts = SelectOpts{"Angola", "ao"}
	CtryArgentina            SelectOpts = SelectOpts{"Argentina", "ar"}
	CtryAustralia            SelectOpts = SelectOpts{"Australia", "au"}
	CtryAustria              SelectOpts = SelectOpts{"Austria", "at"}
	CtryAzerbaijan           SelectOpts = SelectOpts{"Azerbaijan", "az"}
	CtryBahrain              SelectOpts = SelectOpts{"Bahrain", "bh"}
	CtryBangladesh           SelectOpts = SelectOpts{"Bangladesh", "bd"}
	CtryBarbados             SelectOpts = SelectOpts{"Barbados", "bb"}
	CtryBelarus              SelectOpts = SelectOpts{"Belarus", "by"}
	CtryBelgium              SelectOpts = SelectOpts{"Belgium", "be"}
	CtryBermuda              SelectOpts = SelectOpts{"Bermuda", "bm"}
	CtryBhutan               SelectOpts = SelectOpts{"Bhutan", "bt"}
	CtryBolivia              SelectOpts = SelectOpts{"Bolivia", "bo"}
	CtryBosniaAndHerzegovina SelectOpts = SelectOpts{"Bosnia And Herzegovina", "ba"}
	CtryBrazil               SelectOpts = SelectOpts{"Brazil", "br"}
	CtryBrunei               SelectOpts = SelectOpts{"Brunei", "bn"}
	CtryBulgaria             SelectOpts = SelectOpts{"Bulgaria", "bg"}
	CtryBurkinafasco         SelectOpts = SelectOpts{"Burkinafasco", "bf"}
	CtryCambodia             SelectOpts = SelectOpts{"Cambodia", "kh"}
	CtryCameroon             SelectOpts = SelectOpts{"Cameroon", "cm"}
	CtryCanada               SelectOpts = SelectOpts{"Canada", "ca"}
	CtryCapeVerde            SelectOpts = SelectOpts{"CapeVerde", "cv"}
	CtryCaymanIslands        SelectOpts = SelectOpts{"Cayman Islands", "ky"}
	CtryChile                SelectOpts = SelectOpts{"Chile", "cl"}
	CtryChina                SelectOpts = SelectOpts{"China", "cn"}
	CtryColombia             SelectOpts = SelectOpts{"Colombia", "co"}
	CtryComoros              SelectOpts = SelectOpts{"Comoros", "km"}
	CtryCostaRica            SelectOpts = SelectOpts{"Costa Rica", "cr"}
	CtryCotedIvoire          SelectOpts = SelectOpts{"Côte d'Ivoire", "ci"}
	CtryCroatia              SelectOpts = SelectOpts{"Croatia", "hr"}
	CtryCuba                 SelectOpts = SelectOpts{"Cuba", "cu"}
	CtryCyprus               SelectOpts = SelectOpts{"Cyprus", "cy"}
	CtryCzechRepublic        SelectOpts = SelectOpts{"Czech Republic", "cz"}
	CtryDenmark              SelectOpts = SelectOpts{"Denmark", "dk"}
	CtryDjibouti             SelectOpts = SelectOpts{"Djibouti", "dj"}
	CtryDominica             SelectOpts = SelectOpts{"Dominica", "dm"}
	CtryDominicanRepublic    SelectOpts = SelectOpts{"Dominican Republic", "do"}
	CtryDRCongo              SelectOpts = SelectOpts{"Democratic Republic of the Congo", "cd"}
	CtryEcuador              SelectOpts = SelectOpts{"Ecuador", "ec"}
	CtryEgypt                SelectOpts = SelectOpts{"Egypt", "eg"}
	CtryElSalvador           SelectOpts = SelectOpts{"ElSalvador", "sv"}
	CtryEstonia              SelectOpts = SelectOpts{"Estonia", "ee"}
	CtryEthiopia             SelectOpts = SelectOpts{"Ethiopia", "et"}
	CtryFiji                 SelectOpts = SelectOpts{"Fiji", "fj"}
	CtryFinland              SelectOpts = SelectOpts{"Finland", "fi"}
	CtryFrance               SelectOpts = SelectOpts{"France", "fr"}
	CtryFrenchPolynesia      SelectOpts = SelectOpts{"French Polynesia", "pf"}
	CtryGabon                SelectOpts = SelectOpts{"Gabon", "ga"}
	CtryGeorgia              SelectOpts = SelectOpts{"Georgia", "ge"}
	CtryGermany              SelectOpts = SelectOpts{"Germany", "de"}
	CtryGhana                SelectOpts = SelectOpts{"Ghana", "gh"}
	CtryGreece               SelectOpts = SelectOpts{"Greece", "gr"}
	CtryGuatemala            SelectOpts = SelectOpts{"Guatemala", "gt"}
	CtryGuinea               SelectOpts = SelectOpts{"Guinea", "gn"}
	CtryHaiti                SelectOpts = SelectOpts{"Haiti", "ht"}
	CtryHonduras             SelectOpts = SelectOpts{"Honduras", "hn"}
	CtryHongKong             SelectOpts = SelectOpts{"Hong Kong", "hk"}
	CtryHungary              SelectOpts = SelectOpts{"Hungary", "hu"}
	CtryIceland              SelectOpts = SelectOpts{"Iceland", "is"}
	CtryIndia                SelectOpts = SelectOpts{"India", "in"}
	CtryIndonesia            SelectOpts = SelectOpts{"Indonesia", "id"}
	CtryIraq                 SelectOpts = SelectOpts{"Iraq", "iq"}
	CtryIreland              SelectOpts = SelectOpts{"Ireland", "ie"}
	CtryIsrael               SelectOpts = SelectOpts{"Israel", "il"}
	CtryItaly                SelectOpts = SelectOpts{"Italy", "it"}
	CtryJamaica              SelectOpts = SelectOpts{"Jamaica", "jm"}
	CtryJapan                SelectOpts = SelectOpts{"Japan", "jp"}
	CtryJordan               SelectOpts = SelectOpts{"Jordan", "jo"}
	CtryKazakhstan           SelectOpts = SelectOpts{"Kazakhstan", "kz"}
	CtryKenya                SelectOpts = SelectOpts{"Kenya", "ke"}
	CtryKuwait               SelectOpts = SelectOpts{"Kuwait", "kw"}
	CtryKyrgyzstan           SelectOpts = SelectOpts{"Kyrgyzstan", "kg"}
	CtryLatvia               SelectOpts = SelectOpts{"Latvia", "lv"}
	CtryLebanon              SelectOpts = SelectOpts{"Lebanon", "lb"}
	CtryLibya                SelectOpts = SelectOpts{"Libya", "ly"}
	CtryLithuania            SelectOpts = SelectOpts{"Lithuania", "lt"}
	CtryLuxembourg           SelectOpts = SelectOpts{"Luxembourg", "lu"}
	CtryMacau                SelectOpts = SelectOpts{"Macau", "mo"}
	CtryMacedonia            SelectOpts = SelectOpts{"Macedonia", "mk"}
	CtryMadagascar           SelectOpts = SelectOpts{"Madagascar", "mg"}
	CtryMalawi               SelectOpts = SelectOpts{"Malawi", "mw"}
	CtryMalaysia             SelectOpts = SelectOpts{"Malaysia", "my"}
	CtryMaldives             SelectOpts = SelectOpts{"Maldives", "mv"}
	CtryMali                 SelectOpts = SelectOpts{"Mali", "ml"}
	CtryMalta                SelectOpts = SelectOpts{"Malta", "mt"}
	CtryMauritania           SelectOpts = SelectOpts{"Mauritania", "mr"}
	CtryMexico               SelectOpts = SelectOpts{"Mexico", "mx"}
	CtryMoldova              SelectOpts = SelectOpts{"Moldova", "md"}
	CtryMongolia             SelectOpts = SelectOpts{"Mongolia", "mn"}
	CtryMontenegro           SelectOpts = SelectOpts{"Montenegro", "me"}
	CtryMorocco              SelectOpts = SelectOpts{"Morocco", "ma"}
	CtryMozambique           SelectOpts = SelectOpts{"Mozambique", "mz"}
	CtryMyanmar              SelectOpts = SelectOpts{"Myanmar", "mm"}
	CtryNamibia              SelectOpts = SelectOpts{"Namibia", "na"}
	CtryNepal                SelectOpts = SelectOpts{"Nepal", "np"}
	CtryNetherland           SelectOpts = SelectOpts{"Netherland", "nl"}
	CtryNewzealand           SelectOpts = SelectOpts{"Newzealand", "nz"}
	CtryNiger                SelectOpts = SelectOpts{"Niger", "ne"}
	CtryNigeria              SelectOpts = SelectOpts{"Nigeria", "ng"}
	CtryNorthkorea           SelectOpts = SelectOpts{"North Korea", "kp"}
	CtryNorway               SelectOpts = SelectOpts{"Norway", "no"}
	CtryOman                 SelectOpts = SelectOpts{"Oman", "om"}
	CtryPakistan             SelectOpts = SelectOpts{"Pakistan", "pk"}
	CtryPanama               SelectOpts = SelectOpts{"Panama", "pa"}
	CtryParaguay             SelectOpts = SelectOpts{"Paraguay", "py"}
	CtryPeru                 SelectOpts = SelectOpts{"Peru", "pe"}
	CtryPhilippines          SelectOpts = SelectOpts{"Philippines", "ph"}
	CtryPoland               SelectOpts = SelectOpts{"Poland", "pl"}
	CtryPortugal             SelectOpts = SelectOpts{"Portugal", "pt"}
	CtryPuertorico           SelectOpts = SelectOpts{"Puertorico", "pr"}
	CtryRomania              SelectOpts = SelectOpts{"Romania", "ro"}
	CtryRussia               SelectOpts = SelectOpts{"Russia", "ru"}
	CtryRwanda               SelectOpts = SelectOpts{"Rwanda", "rw"}
	CtrySamoa                SelectOpts = SelectOpts{"Samoa", "ws"}
	CtrySanMarino            SelectOpts = SelectOpts{"SanMarino", "sm"}
	CtrySaudiarabia          SelectOpts = SelectOpts{"Saudiarabia", "sa"}
	CtrySenegal              SelectOpts = SelectOpts{"Senegal", "sn"}
	CtrySerbia               SelectOpts = SelectOpts{"Serbia", "rs"}
	CtrySingapore            SelectOpts = SelectOpts{"Singapore", "sg"}
	CtrySlovakia             SelectOpts = SelectOpts{"Slovakia", "sk"}
	CtrySlovenia             SelectOpts = SelectOpts{"Slovenia", "si"}
	CtrySolomonIslands       SelectOpts = SelectOpts{"Solomon Islands", "sb"}
	CtrySomalia              SelectOpts = SelectOpts{"Somalia", "so"}
	CtrySouthAfrica          SelectOpts = SelectOpts{"South Africa", "za"}
	CtrySouthKorea           SelectOpts = SelectOpts{"South Korea", "kr"}
	CtrySpain                SelectOpts = SelectOpts{"Spain", "es"}
	CtrySriLanka             SelectOpts = SelectOpts{"Sri Lanka", "lk"}
	CtrySudan                SelectOpts = SelectOpts{"Sudan", "sd"}
	CtrySweden               SelectOpts = SelectOpts{"Sweden", "se"}
	CtrySwitzerland          SelectOpts = SelectOpts{"Switzerland", "ch"}
	CtrySyria                SelectOpts = SelectOpts{"Syria", "sy"}
	CtryTaiwan               SelectOpts = SelectOpts{"Taiwan", "tw"}
	CtryTajikistan           SelectOpts = SelectOpts{"Tajikistan", "tj"}
	CtryTanzania             SelectOpts = SelectOpts{"Tanzania", "tz"}
	CtryThailand             SelectOpts = SelectOpts{"Thailand", "th"}
	CtryTonga                SelectOpts = SelectOpts{"Tonga", "to"}
	CtryTunisia              SelectOpts = SelectOpts{"Tunisia", "tn"}
	CtryTurkey               SelectOpts = SelectOpts{"Turkey", "tr"}
	CtryTurkmenistan         SelectOpts = SelectOpts{"Turkmenistan", "tm"}
	CtryUganda               SelectOpts = SelectOpts{"Uganda", "ug"}
	CtryUkraine              SelectOpts = SelectOpts{"Ukraine", "ua"}
	CtryUnitedArabemirates   SelectOpts = SelectOpts{"United Arabemirates", "ae"}
	CtryUnitedKingdom        SelectOpts = SelectOpts{"United Kingdom", "gb"}
	CtryUnitedStates         SelectOpts = SelectOpts{"United States", "us"}
	CtryUruguay              SelectOpts = SelectOpts{"Uruguay", "uy"}
	CtryUzbekistan           SelectOpts = SelectOpts{"Uzbekistan", "uz"}
	CtryVenezuela            SelectOpts = SelectOpts{"Venezuela", "ve"}
	CtryVietnam              SelectOpts = SelectOpts{"Vietnam", "vi"}
	CtryYemen                SelectOpts = SelectOpts{"Yemen", "ye"}
	CtryZambia               SelectOpts = SelectOpts{"Zambia", "zm"}
	CtryZimbabwe             SelectOpts = SelectOpts{"Zimbabwe", "zw"}
)

var CtryList = [][2]string{
	CtryAfghanistan,
	CtryAlbania,
	CtryAlgeria,
	CtryAngola,
	CtryArgentina,
	CtryAustralia,
	CtryAustria,
	CtryAzerbaijan,
	CtryBahrain,
	CtryBangladesh,
	CtryBarbados,
	CtryBelarus,
	CtryBelgium,
	CtryBermuda,
	CtryBhutan,
	CtryBolivia,
	CtryBosniaAndHerzegovina,
	CtryBrazil,
	CtryBrunei,
	CtryBulgaria,
	CtryBurkinafasco,
	CtryCambodia,
	CtryCameroon,
	CtryCanada,
	CtryCapeVerde,
	CtryCaymanIslands,
	CtryChile,
	CtryChina,
	CtryColombia,
	CtryComoros,
	CtryCostaRica,
	CtryCotedIvoire,
	CtryCroatia,
	CtryCuba,
	CtryCyprus,
	CtryCzechRepublic,
	CtryDenmark,
	CtryDjibouti,
	CtryDominica,
	CtryDominicanRepublic,
	CtryDRCongo,
	CtryEcuador,
	CtryEgypt,
	CtryElSalvador,
	CtryEstonia,
	CtryEthiopia,
	CtryFiji,
	CtryFinland,
	CtryFrance,
	CtryFrenchPolynesia,
	CtryGabon,
	CtryGeorgia,
	CtryGermany,
	CtryGhana,
	CtryGreece,
	CtryGuatemala,
	CtryGuinea,
	CtryHaiti,
	CtryHonduras,
	CtryHongKong,
	CtryHungary,
	CtryIceland,
	CtryIndia,
	CtryIndonesia,
	CtryIraq,
	CtryIreland,
	CtryIsrael,
	CtryItaly,
	CtryJamaica,
	CtryJapan,
	CtryJordan,
	CtryKazakhstan,
	CtryKenya,
	CtryKuwait,
	CtryKyrgyzstan,
	CtryLatvia,
	CtryLebanon,
	CtryLibya,
	CtryLithuania,
	CtryLuxembourg,
	CtryMacau,
	CtryMacedonia,
	CtryMadagascar,
	CtryMalawi,
	CtryMalaysia,
	CtryMaldives,
	CtryMali,
	CtryMalta,
	CtryMauritania,
	CtryMexico,
	CtryMoldova,
	CtryMongolia,
	CtryMontenegro,
	CtryMorocco,
	CtryMozambique,
	CtryMyanmar,
	CtryNamibia,
	CtryNepal,
	CtryNetherland,
	CtryNewzealand,
	CtryNiger,
	CtryNigeria,
	CtryNorthkorea,
	CtryNorway,
	CtryOman,
	CtryPakistan,
	CtryPanama,
	CtryParaguay,
	CtryPeru,
	CtryPhilippines,
	CtryPoland,
	CtryPortugal,
	CtryPuertorico,
	CtryRomania,
	CtryRussia,
	CtryRwanda,
	CtrySamoa,
	CtrySanMarino,
	CtrySaudiarabia,
	CtrySenegal,
	CtrySerbia,
	CtrySingapore,
	CtrySlovakia,
	CtrySlovenia,
	CtrySolomonIslands,
	CtrySomalia,
	CtrySouthAfrica,
	CtrySouthKorea,
	CtrySpain,
	CtrySriLanka,
	CtrySudan,
	CtrySweden,
	CtrySwitzerland,
	CtrySyria,
	CtryTaiwan,
	CtryTajikistan,
	CtryTanzania,
	CtryThailand,
	CtryTonga,
	CtryTunisia,
	CtryTurkey,
	CtryTurkmenistan,
	CtryUganda,
	CtryUkraine,
	CtryUnitedArabemirates,
	CtryUnitedKingdom,
	CtryUnitedStates,
	CtryUruguay,
	CtryUzbekistan,
	CtryVenezuela,
	CtryVietnam,
	CtryYemen,
	CtryZambia,
	CtryZimbabwe,
}

var CtryVal validator.Enmus[string]

var (
	LangAfrikaans      SelectOpts = SelectOpts{"Afrikaans", "af"}
	LangAlbanian       SelectOpts = SelectOpts{"Albanian", "sq"}
	LangAmharic        SelectOpts = SelectOpts{"Amharic", "am"}
	LangArabic         SelectOpts = SelectOpts{"Arabic", "ar"}
	LangAssamese       SelectOpts = SelectOpts{"Assamese", "as"}
	LangAzerbaijani    SelectOpts = SelectOpts{"Azerbaijani", "az"}
	LangBelarusian     SelectOpts = SelectOpts{"Belarusian", "be"}
	LangBengali        SelectOpts = SelectOpts{"Bengali", "bn"}
	LangBosnian        SelectOpts = SelectOpts{"Bosnian", "bs"}
	LangBulgarian      SelectOpts = SelectOpts{"Bulgarian", "bg"}
	LangBurmese        SelectOpts = SelectOpts{"Burmese", "my"}
	LangCatalan        SelectOpts = SelectOpts{"Catalan", "ca"}
	LangCentralKurdish SelectOpts = SelectOpts{"CentralKurdish", "ckb"}
	LangChinese        SelectOpts = SelectOpts{"Chinese", "zh"}
	LangCroatian       SelectOpts = SelectOpts{"Croatian", "hr"}
	LangCzech          SelectOpts = SelectOpts{"Czech", "cs"}
	LangDanish         SelectOpts = SelectOpts{"Danish", "da"}
	LangDutch          SelectOpts = SelectOpts{"Dutch", "nl"}
	LangEnglish        SelectOpts = SelectOpts{"English", "en"}
	LangEstonian       SelectOpts = SelectOpts{"Estonian", "et"}
	LangFilipino       SelectOpts = SelectOpts{"Filipino", "pi"}
	LangFinnish        SelectOpts = SelectOpts{"Finnish", "fi"}
	LangFrench         SelectOpts = SelectOpts{"French", "fr"}
	LangGeorgian       SelectOpts = SelectOpts{"Georgian", "ka"}
	LangGerman         SelectOpts = SelectOpts{"German", "de"}
	LangGreek          SelectOpts = SelectOpts{"Greek", "el"}
	LangGujarati       SelectOpts = SelectOpts{"Gujarati", "gu"}
	LangHebrew         SelectOpts = SelectOpts{"Hebrew", "he"}
	LangHindi          SelectOpts = SelectOpts{"Hindi", "hi"}
	LangHungarian      SelectOpts = SelectOpts{"Hungarian", "hu"}
	LangIcelandic      SelectOpts = SelectOpts{"Icelandic", "is"}
	LangIndonesian     SelectOpts = SelectOpts{"Indonesian", "id"}
	LangItalian        SelectOpts = SelectOpts{"Italian", "it"}
	LangJapanese       SelectOpts = SelectOpts{"Japanese", "jp"}
	LangKhmer          SelectOpts = SelectOpts{"Khmer", "kh"}
	LangKinyarwanda    SelectOpts = SelectOpts{"Kinyarwanda", "rw"}
	LangKorean         SelectOpts = SelectOpts{"Korean", "ko"}
	LangLatvian        SelectOpts = SelectOpts{"Latvian", "lv"}
	LangLithuanian     SelectOpts = SelectOpts{"Lithuanian", "lt"}
	LangLuxembourgish  SelectOpts = SelectOpts{"Luxembourgish", "lb"}
	LangMacedonian     SelectOpts = SelectOpts{"Macedonian", "mk"}
	LangMalay          SelectOpts = SelectOpts{"Malay", "ms"}
	LangMalayalam      SelectOpts = SelectOpts{"Malayalam", "ml"}
	LangMaltese        SelectOpts = SelectOpts{"Maltese", "mt"}
	LangMaori          SelectOpts = SelectOpts{"Maori", "mi"}
	LangMarathi        SelectOpts = SelectOpts{"Marathi", "mr"}
	LangMongolian      SelectOpts = SelectOpts{"Mongolian", "mn"}
	LangNepali         SelectOpts = SelectOpts{"Nepali", "ne"}
	LangNorwegian      SelectOpts = SelectOpts{"Norwegian", "no"}
	LangOriya          SelectOpts = SelectOpts{"Oriya", "or"}
	LangPashto         SelectOpts = SelectOpts{"Pashto", "ps"}
	LangPersian        SelectOpts = SelectOpts{"Persian", "fa"}
	LangPolish         SelectOpts = SelectOpts{"Polish", "pl"}
	LangPortuguese     SelectOpts = SelectOpts{"Portuguese", "pt"}
	LangPunjabi        SelectOpts = SelectOpts{"Punjabi", "pa"}
	LangRomanian       SelectOpts = SelectOpts{"Romanian", "ro"}
	LangRussian        SelectOpts = SelectOpts{"Russian", "ru"}
	LangSamoan         SelectOpts = SelectOpts{"Samoan", "sm"}
	LangSerbian        SelectOpts = SelectOpts{"Serbian", "sr"}
	LangShona          SelectOpts = SelectOpts{"Shona", "sn"}
	LangSinhala        SelectOpts = SelectOpts{"Sinhala", "si"}
	LangSlovak         SelectOpts = SelectOpts{"Slovak", "sk"}
	LangSlovenian      SelectOpts = SelectOpts{"Slovenian", "sl"}
	LangSomali         SelectOpts = SelectOpts{"Somali", "so"}
	LangSpanish        SelectOpts = SelectOpts{"Spanish", "es"}
	LangSwahili        SelectOpts = SelectOpts{"Swahili", "sw"}
	LangSwedish        SelectOpts = SelectOpts{"Swedish", "sv"}
	LangTajik          SelectOpts = SelectOpts{"Tajik", "tg"}
	LangTamil          SelectOpts = SelectOpts{"Tamil", "ta"}
	LangTelugu         SelectOpts = SelectOpts{"Telugu", "te"}
	LangThai           SelectOpts = SelectOpts{"Thai", "th"}
	LangTurkish        SelectOpts = SelectOpts{"Turkish", "tr"}
	LangTurkmen        SelectOpts = SelectOpts{"Turkmen", "tk"}
	LangUkrainian      SelectOpts = SelectOpts{"Ukrainian", "uk"}
	LangUrdu           SelectOpts = SelectOpts{"Urdu", "ur"}
	LangUzbek          SelectOpts = SelectOpts{"Uzbek", "uz"}
	LangVietnamese     SelectOpts = SelectOpts{"Vietnamese", "vi"}
)

var LangList = [][2]string{
	LangAfrikaans,
	LangAlbanian,
	LangAmharic,
	LangArabic,
	LangAssamese,
	LangAzerbaijani,
	LangBelarusian,
	LangBengali,
	LangBosnian,
	LangBulgarian,
	LangBurmese,
	LangCatalan,
	LangCentralKurdish,
	LangChinese,
	LangCroatian,
	LangCzech,
	LangDanish,
	LangDutch,
	LangEnglish,
	LangEstonian,
	LangFilipino,
	LangFinnish,
	LangFrench,
	LangGeorgian,
	LangGerman,
	LangGreek,
	LangGujarati,
	LangHebrew,
	LangHindi,
	LangHungarian,
	LangIcelandic,
	LangIndonesian,
	LangItalian,
	LangJapanese,
	LangKhmer,
	LangKinyarwanda,
	LangKorean,
	LangLatvian,
	LangLithuanian,
	LangLuxembourgish,
	LangMacedonian,
	LangMalay,
	LangMalayalam,
	LangMaltese,
	LangMaori,
	LangMarathi,
	LangMongolian,
	LangNepali,
	LangNorwegian,
	LangOriya,
	LangPashto,
	LangPersian,
	LangPolish,
	LangPortuguese,
	LangPunjabi,
	LangRomanian,
	LangRussian,
	LangSamoan,
	LangSerbian,
	LangShona,
	LangSinhala,
	LangSlovak,
	LangSlovenian,
	LangSomali,
	LangSpanish,
	LangSwahili,
	LangSwedish,
	LangTajik,
	LangTamil,
	LangTelugu,
	LangThai,
	LangTurkish,
	LangTurkmen,
	LangUkrainian,
	LangUrdu,
	LangUzbek,
	LangVietnamese,
}

var LangVal validator.Enmus[string]
