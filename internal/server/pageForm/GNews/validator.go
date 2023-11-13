package gnews

import "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"

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
