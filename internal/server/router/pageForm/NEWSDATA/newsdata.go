package newsdata

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
		NEWSDATAIOLatestNews{},
		NEWSDATAIONewsArchive{},
		NEWSDATAIONewsSources{},
	}

	for _, pf := range pfs {
		pageform.PageFormRepo.Add(pf)
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

const (
	Business      string = "business"
	Entertainment        = "entertainment"
	Environment          = "environment"
	Food                 = "food"
	Health               = "health"
	Politics             = "politics"
	Science              = "science"
	Sports               = "sports"
	Technology           = "technology"
	Top                  = "top"
	Tourism              = "tourism"
	World                = "world"
)

var Category = map[string]string{
	Business:      "Business",
	Entertainment: "Entertainment",
	Environment:   "Environment",
	Food:          "Food",
	Health:        "Health",
	Politics:      "Politics",
	Science:       "Science",
	Sports:        "Sports",
	Technology:    "Technology",
	Top:           "Top",
	Tourism:       "Tourism",
	World:         "World",
}

var CategoryValidator = validator.NewEnmusFromMap(
	VAL_TAG_CATEGORY, Category, "key",
)

const (
	Afghanistan          string = "af"
	Albania                     = "al"
	Algeria                     = "dz"
	Angola                      = "ao"
	Argentina                   = "ar"
	Australia                   = "au"
	Austria                     = "at"
	Azerbaijan                  = "az"
	Bahrain                     = "bh"
	Bangladesh                  = "bd"
	Barbados                    = "bb"
	Belarus                     = "by"
	Belgium                     = "be"
	Bermuda                     = "bm"
	Bhutan                      = "bt"
	Bolivia                     = "bo"
	BosniaAndHerzegovina        = "ba"
	Brazil                      = "br"
	Brunei                      = "bn"
	Bulgaria                    = "bg"
	Burkinafasco                = "bf"
	Cambodia                    = "kh"
	Cameroon                    = "cm"
	Canada                      = "ca"
	CapeVerde                   = "cv"
	CaymanIslands               = "ky"
	Chile                       = "cl"
	China                       = "cn"
	Colombia                    = "co"
	Comoros                     = "km"
	CostaRica                   = "cr"
	CotedIvoire                 = "ci"
	Croatia                     = "hr"
	Cuba                        = "cu"
	Cyprus                      = "cy"
	CzechRepublic               = "cz"
	Denmark                     = "dk"
	Djibouti                    = "dj"
	Dominica                    = "dm"
	DominicanRepublic           = "do"
	DRCongo                     = "cd"
	Ecuador                     = "ec"
	Egypt                       = "eg"
	ElSalvador                  = "sv"
	Estonia                     = "ee"
	Ethiopia                    = "et"
	Fiji                        = "fj"
	Finland                     = "fi"
	France                      = "fr"
	FrenchPolynesia             = "pf"
	Gabon                       = "ga"
	Georgia                     = "ge"
	Germany                     = "de"
	Ghana                       = "gh"
	Greece                      = "gr"
	Guatemala                   = "gt"
	Guinea                      = "gn"
	Haiti                       = "ht"
	Honduras                    = "hn"
	HongKong                    = "hk"
	Hungary                     = "hu"
	Iceland                     = "is"
	India                       = "in"
	Indonesia                   = "id"
	Iraq                        = "iq"
	Ireland                     = "ie"
	Israel                      = "il"
	Italy                       = "it"
	Jamaica                     = "jm"
	Japan                       = "jp"
	Jordan                      = "jo"
	Kazakhstan                  = "kz"
	Kenya                       = "ke"
	Kuwait                      = "kw"
	Kyrgyzstan                  = "kg"
	Latvia                      = "lv"
	Lebanon                     = "lb"
	Libya                       = "ly"
	Lithuania                   = "lt"
	Luxembourg                  = "lu"
	Macau                       = "mo"
	Macedonia                   = "mk"
	Madagascar                  = "mg"
	Malawi                      = "mw"
	Malaysia                    = "my"
	Maldives                    = "mv"
	Mali                        = "ml"
	Malta                       = "mt"
	Mauritania                  = "mr"
	Mexico                      = "mx"
	Moldova                     = "md"
	Mongolia                    = "mn"
	Montenegro                  = "me"
	Morocco                     = "ma"
	Mozambique                  = "mz"
	Myanmar                     = "mm"
	Namibia                     = "na"
	Nepal                       = "np"
	Netherland                  = "nl"
	Newzealand                  = "nz"
	Niger                       = "ne"
	Nigeria                     = "ng"
	Northkorea                  = "kp"
	Norway                      = "no"
	Oman                        = "om"
	Pakistan                    = "pk"
	Panama                      = "pa"
	Paraguay                    = "py"
	Peru                        = "pe"
	Philippines                 = "ph"
	Poland                      = "pl"
	Portugal                    = "pt"
	Puertorico                  = "pr"
	Romania                     = "ro"
	Russia                      = "ru"
	Rwanda                      = "rw"
	Samoa                       = "ws"
	SanMarino                   = "sm"
	Saudiarabia                 = "sa"
	Senegal                     = "sn"
	Serbia                      = "rs"
	Singapore                   = "sg"
	Slovakia                    = "sk"
	Slovenia                    = "si"
	SolomonIslands              = "sb"
	Somalia                     = "so"
	SouthAfrica                 = "za"
	SouthKorea                  = "kr"
	Spain                       = "es"
	SriLanka                    = "lk"
	Sudan                       = "sd"
	Sweden                      = "se"
	Switzerland                 = "ch"
	Syria                       = "sy"
	Taiwan                      = "tw"
	Tajikistan                  = "tj"
	Tanzania                    = "tz"
	Thailand                    = "th"
	Tonga                       = "to"
	Tunisia                     = "tn"
	Turkey                      = "tr"
	Turkmenistan                = "tm"
	Uganda                      = "ug"
	Ukraine                     = "ua"
	UnitedArabemirates          = "ae"
	UnitedKingdom               = "gb"
	UnitedStates                = "us"
	Uruguay                     = "uy"
	Uzbekistan                  = "uz"
	Venezuela                   = "ve"
	Vietnam                     = "vi"
	Yemen                       = "ye"
	Zambia                      = "zm"
	Zimbabwe                    = "zw"
)

var Country = map[string]string{
	Afghanistan:          "Afghanistan",
	Albania:              "Albania",
	Algeria:              "Algeria",
	Angola:               "Angola",
	Argentina:            "Argentina",
	Australia:            "Australia",
	Austria:              "Austria",
	Azerbaijan:           "Azerbaijan",
	Bahrain:              "Bahrain",
	Bangladesh:           "Bangladesh",
	Barbados:             "Barbados",
	Belarus:              "Belarus",
	Belgium:              "Belgium",
	Bermuda:              "Bermuda",
	Bhutan:               "Bhutan",
	Bolivia:              "Bolivia",
	BosniaAndHerzegovina: "Bosnia And Herzegovina",
	Brazil:               "Brazil",
	Brunei:               "Brunei",
	Bulgaria:             "Bulgaria",
	Burkinafasco:         "Burkinafasco",
	Cambodia:             "Cambodia",
	Cameroon:             "Cameroon",
	Canada:               "Canada",
	CapeVerde:            "CapeVerde",
	CaymanIslands:        "Cayman Islands",
	Chile:                "Chile",
	China:                "China",
	Colombia:             "Colombia",
	Comoros:              "Comoros",
	CostaRica:            "Costa Rica",
	CotedIvoire:          "Côte d'Ivoire",
	Croatia:              "Croatia",
	Cuba:                 "Cuba",
	Cyprus:               "Cyprus",
	CzechRepublic:        "Czech Republic",
	Denmark:              "Denmark",
	Djibouti:             "Djibouti",
	Dominica:             "Dominica",
	DominicanRepublic:    "Dominican Republic",
	DRCongo:              "Democratic Republic of the Congo",
	Ecuador:              "Ecuador",
	Egypt:                "Egypt",
	ElSalvador:           "ElSalvador",
	Estonia:              "Estonia",
	Ethiopia:             "Ethiopia",
	Fiji:                 "Fiji",
	Finland:              "Finland",
	France:               "France",
	FrenchPolynesia:      "French Polynesia",
	Gabon:                "Gabon",
	Georgia:              "Georgia",
	Germany:              "Germany",
	Ghana:                "Ghana",
	Greece:               "Greece",
	Guatemala:            "Guatemala",
	Guinea:               "Guinea",
	Haiti:                "Haiti",
	Honduras:             "Honduras",
	HongKong:             "Hong Kong",
	Hungary:              "Hungary",
	Iceland:              "Iceland",
	India:                "India",
	Indonesia:            "Indonesia",
	Iraq:                 "Iraq",
	Ireland:              "Ireland",
	Israel:               "Israel",
	Italy:                "Italy",
	Jamaica:              "Jamaica",
	Japan:                "Japan",
	Jordan:               "Jordan",
	Kazakhstan:           "Kazakhstan",
	Kenya:                "Kenya",
	Kuwait:               "Kuwait",
	Kyrgyzstan:           "Kyrgyzstan",
	Latvia:               "Latvia",
	Lebanon:              "Lebanon",
	Libya:                "Libya",
	Lithuania:            "Lithuania",
	Luxembourg:           "Luxembourg",
	Macau:                "Macau",
	Macedonia:            "Macedonia",
	Madagascar:           "Madagascar",
	Malawi:               "Malawi",
	Malaysia:             "Malaysia",
	Maldives:             "Maldives",
	Mali:                 "Mali",
	Malta:                "Malta",
	Mauritania:           "Mauritania",
	Mexico:               "Mexico",
	Moldova:              "Moldova",
	Mongolia:             "Mongolia",
	Montenegro:           "Montenegro",
	Morocco:              "Morocco",
	Mozambique:           "Mozambique",
	Myanmar:              "Myanmar",
	Namibia:              "Namibia",
	Nepal:                "Nepal",
	Netherland:           "Netherland",
	Newzealand:           "Newzealand",
	Niger:                "Niger",
	Nigeria:              "Nigeria",
	Northkorea:           "North Korea",
	Norway:               "Norway",
	Oman:                 "Oman",
	Pakistan:             "Pakistan",
	Panama:               "Panama",
	Paraguay:             "Paraguay",
	Peru:                 "Peru",
	Philippines:          "Philippines",
	Poland:               "Poland",
	Portugal:             "Portugal",
	Puertorico:           "Puertorico",
	Romania:              "Romania",
	Russia:               "Russia",
	Rwanda:               "Rwanda",
	Samoa:                "Samoa",
	SanMarino:            "SanMarino",
	Saudiarabia:          "Saudiarabia",
	Senegal:              "Senegal",
	Serbia:               "Serbia",
	Singapore:            "Singapore",
	Slovakia:             "Slovakia",
	Slovenia:             "Slovenia",
	SolomonIslands:       "Solomon Islands",
	Somalia:              "Somalia",
	SouthAfrica:          "South Africa",
	SouthKorea:           "South Korea",
	Spain:                "Spain",
	SriLanka:             "Sri Lanka",
	Sudan:                "Sudan",
	Sweden:               "Sweden",
	Switzerland:          "Switzerland",
	Syria:                "Syria",
	Taiwan:               "Taiwan",
	Tajikistan:           "Tajikistan",
	Tanzania:             "Tanzania",
	Thailand:             "Thailand",
	Tonga:                "Tonga",
	Tunisia:              "Tunisia",
	Turkey:               "Turkey",
	Turkmenistan:         "Turkmenistan",
	Uganda:               "Uganda",
	Ukraine:              "Ukraine",
	UnitedArabemirates:   "United Arabemirates",
	UnitedKingdom:        "United Kingdom",
	UnitedStates:         "United States",
	Uruguay:              "Uruguay",
	Uzbekistan:           "Uzbekistan",
	Venezuela:            "Venezuela",
	Vietnam:              "Vietnam",
	Yemen:                "Yemen",
	Zambia:               "Zambia",
	Zimbabwe:             "Zimbabwe",
}

var CountryValidator = validator.NewEnmusFromMap(
	VAL_TAG_COUNTRY, Country, "key",
)

const (
	Afrikaans      = "af"
	Albanian       = "sq"
	Amharic        = "am"
	Arabic         = "ar"
	Assamese       = "as"
	Azerbaijani    = "az"
	Belarusian     = "be"
	Bengali        = "bn"
	Bosnian        = "bs"
	Bulgarian      = "bg"
	Burmese        = "my"
	Catalan        = "ca"
	CentralKurdish = "ckb"
	Chinese        = "zh"
	Croatian       = "hr"
	Czech          = "cs"
	Danish         = "da"
	Dutch          = "nl"
	English        = "en"
	Estonian       = "et"
	Filipino       = "pi"
	Finnish        = "fi"
	French         = "fr"
	Georgian       = "ka"
	German         = "de"
	Greek          = "el"
	Gujarati       = "gu"
	Hebrew         = "he"
	Hindi          = "hi"
	Hungarian      = "hu"
	Icelandic      = "is"
	Indonesian     = "id"
	Italian        = "it"
	Japanese       = "jp"
	Khmer          = "kh"
	Kinyarwanda    = "rw"
	Korean         = "ko"
	Latvian        = "lv"
	Lithuanian     = "lt"
	Luxembourgish  = "lb"
	Macedonian     = "mk"
	Malay          = "ms"
	Malayalam      = "ml"
	Maltese        = "mt"
	Maori          = "mi"
	Marathi        = "mr"
	Mongolian      = "mn"
	Nepali         = "ne"
	Norwegian      = "no"
	Oriya          = "or"
	Pashto         = "ps"
	Persian        = "fa"
	Polish         = "pl"
	Portuguese     = "pt"
	Punjabi        = "pa"
	Romanian       = "ro"
	Russian        = "ru"
	Samoan         = "sm"
	Serbian        = "sr"
	Shona          = "sn"
	Sinhala        = "si"
	Slovak         = "sk"
	Slovenian      = "sl"
	Somali         = "so"
	Spanish        = "es"
	Swahili        = "sw"
	Swedish        = "sv"
	Tajik          = "tg"
	Tamil          = "ta"
	Telugu         = "te"
	Thai           = "th"
	Turkish        = "tr"
	Turkmen        = "tk"
	Ukrainian      = "uk"
	Urdu           = "ur"
	Uzbek          = "uz"
	Vietnamese     = "vi"
)

var Language = map[string]string{
	Afrikaans:      "Afrikaans",
	Albanian:       "Albanian",
	Amharic:        "Amharic",
	Arabic:         "Arabic",
	Assamese:       "Assamese",
	Azerbaijani:    "Azerbaijani",
	Belarusian:     "Belarusian",
	Bengali:        "Bengali",
	Bosnian:        "Bosnian",
	Bulgarian:      "Bulgarian",
	Burmese:        "Burmese",
	Catalan:        "Catalan",
	CentralKurdish: "CentralKurdish",
	Chinese:        "Chinese",
	Croatian:       "Croatian",
	Czech:          "Czech",
	Danish:         "Danish",
	Dutch:          "Dutch",
	English:        "English",
	Estonian:       "Estonian",
	Filipino:       "Filipino",
	Finnish:        "Finnish",
	French:         "French",
	Georgian:       "Georgian",
	German:         "German",
	Greek:          "Greek",
	Gujarati:       "Gujarati",
	Hebrew:         "Hebrew",
	Hindi:          "Hindi",
	Hungarian:      "Hungarian",
	Icelandic:      "Icelandic",
	Indonesian:     "Indonesian",
	Italian:        "Italian",
	Japanese:       "Japanese",
	Khmer:          "Khmer",
	Kinyarwanda:    "Kinyarwanda",
	Korean:         "Korean",
	Latvian:        "Latvian",
	Lithuanian:     "Lithuanian",
	Luxembourgish:  "Luxembourgish",
	Macedonian:     "Macedonian",
	Malay:          "Malay",
	Malayalam:      "Malayalam",
	Maltese:        "Maltese",
	Maori:          "Maori",
	Marathi:        "Marathi",
	Mongolian:      "Mongolian",
	Nepali:         "Nepali",
	Norwegian:      "Norwegian",
	Oriya:          "Oriya",
	Pashto:         "Pashto",
	Persian:        "Persian",
	Polish:         "Polish",
	Portuguese:     "Portuguese",
	Punjabi:        "Punjabi",
	Romanian:       "Romanian",
	Russian:        "Russian",
	Samoan:         "Samoan",
	Serbian:        "Serbian",
	Shona:          "Shona",
	Sinhala:        "Sinhala",
	Slovak:         "Slovak",
	Slovenian:      "Slovenian",
	Somali:         "Somali",
	Spanish:        "Spanish",
	Swahili:        "Swahili",
	Swedish:        "Swedish",
	Tajik:          "Tajik",
	Tamil:          "Tamil",
	Telugu:         "Telugu",
	Thai:           "Thai",
	Turkish:        "Turkish",
	Turkmen:        "Turkmen",
	Ukrainian:      "Ukrainian",
	Urdu:           "Urdu",
	Uzbek:          "Uzbek",
	Vietnamese:     "Vietnamese",
}

var LanguageValidator = validator.NewEnmusFromMap(
	VAL_TAG_LANGUAGE, Language, "key",
)
