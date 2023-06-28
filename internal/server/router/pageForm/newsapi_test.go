package pageform_test

import (
	"testing"

	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	val "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestNEWSAPIFormValidationStruct(t *testing.T) {
	val := val.New()
	err := validator.RegisterValidator(
		val,
		pageform.NEWSAPICatVal,
		pageform.NEWSAPICtryVal,
		pageform.NEWSAPILangVal,
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

	require.NoError(t, val.Var("business", pageform.NEWSAPICatVal.Tag()))
	require.NoError(t, val.Struct(valCatStruct{Category: "business"}))
	require.Error(t, val.Var("xx", pageform.NEWSAPICatVal.Tag()))
	require.Error(t, val.Struct(valCatStruct{Category: "xx"}))

	require.NoError(t, val.Var("tw", pageform.NEWSAPICtryVal.Tag()))
	require.NoError(t, val.Struct(valCtryStruct{Country: "tw"}))
	require.Error(t, val.Var("xx", pageform.NEWSAPICtryVal.Tag()))
	require.Error(t, val.Struct(valCtryStruct{Country: "xx"}))

	require.NoError(t, val.Var("en", pageform.NEWSAPILangVal.Tag()))
	require.NoError(t, val.Struct(valLangStruct{Language: "en"}))
	require.Error(t, val.Var("xx", pageform.NEWSAPILangVal.Tag()))
	require.Error(t, val.Struct(valLangStruct{Language: "xx"}))
}
