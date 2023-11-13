package newsapi_test

import (
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/newsapi"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	val "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestNEWSAPIFormValidationStruct(t *testing.T) {
	val := val.New()
	err := validator.RegisterValidator(
		val,
		newsapi.CategoryValidator,
		newsapi.CountryValidator,
		newsapi.LanguageValidator,
	)
	require.NoError(t, err)

	type valCatStruct struct {
		Category string `validate:"newsapi_cat"`
	}

	type valCtryStruct struct {
		Country string `validate:"newsapi_ctry"`
	}

	type valLangStruct struct {
		Language string `validate:"newsapi_lang"`
	}

	require.NoError(t, val.Var("business", newsapi.CategoryValidator.Tag()))
	require.NoError(t, val.Struct(valCatStruct{Category: "business"}))
	require.Error(t, val.Var("xx", newsapi.CategoryValidator.Tag()))
	require.Error(t, val.Struct(valCatStruct{Category: "xx"}))

	require.NoError(t, val.Var("tw", newsapi.CountryValidator.Tag()))
	require.NoError(t, val.Struct(valCtryStruct{Country: "tw"}))
	require.Error(t, val.Var("xx", newsapi.CountryValidator.Tag()))
	require.Error(t, val.Struct(valCtryStruct{Country: "xx"}))

	require.NoError(t, val.Var("en", newsapi.LanguageValidator.Tag()))
	require.NoError(t, val.Struct(valLangStruct{Language: "en"}))
	require.Error(t, val.Var("xx", newsapi.LanguageValidator.Tag()))
	require.Error(t, val.Struct(valLangStruct{Language: "xx"}))
}
