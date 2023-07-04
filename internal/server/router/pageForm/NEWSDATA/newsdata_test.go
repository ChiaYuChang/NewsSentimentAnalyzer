package newsdata_test

import (
	"testing"

	newsdata "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm/NEWSDATA"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	val "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestNEWSDATAIOFormValidationStruct(t *testing.T) {
	val := val.New()
	err := validator.RegisterValidator(
		val,
		newsdata.LanguageValidator,
		newsdata.CountryValidator,
		newsdata.CategoryValidator,
	)
	require.NoError(t, err)

	type valCatStruct struct {
		Category string `validate:"newsdata_cat"`
	}

	type valCtryStruct struct {
		Country string `validate:"newsdata_ctry"`
	}

	type valLangStruct struct {
		Language string `validate:"newsdata_lang"`
	}

	require.NoError(t, val.Var("business", newsdata.CategoryValidator.Tag()))
	require.NoError(t, val.Struct(valCatStruct{Category: "business"}))
	require.Error(t, val.Var("xx", newsdata.CategoryValidator.Tag()))
	require.Error(t, val.Struct(valCatStruct{Category: "xx"}))

	require.NoError(t, val.Var("tw", newsdata.CountryValidator.Tag()))
	require.NoError(t, val.Struct(valCtryStruct{Country: "tw"}))
	require.Error(t, val.Var("xx", newsdata.CountryValidator.Tag()))
	require.Error(t, val.Struct(valCtryStruct{Country: "xx"}))

	require.NoError(t, val.Var("en", newsdata.LanguageValidator.Tag()))
	require.NoError(t, val.Struct(valLangStruct{Language: "en"}))
	require.Error(t, val.Var("xx", newsdata.LanguageValidator.Tag()))
	require.Error(t, val.Struct(valLangStruct{Language: "xx"}))
}
