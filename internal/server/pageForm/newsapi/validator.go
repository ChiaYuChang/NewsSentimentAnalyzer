package newsapi

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
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
