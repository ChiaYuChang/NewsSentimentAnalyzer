package pageform_test

// import (
// 	"testing"

// 	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
// 	val "github.com/go-playground/validator/v10"
// 	"github.com/stretchr/testify/require"
// )

// func TestNEWSDATAIOFormValidationStruct(t *testing.T) {
// 	val := val.New()
// 	err := validator.RegisterValidator(
// 		val,
// 		pageform.NEWSDATAIOCatVal,
// 		pageform.NEWSDATAIOCtryVal,
// 		pageform.NEWSDATAIOLangVal,
// 	)
// 	require.NoError(t, err)

// 	type valCatStruct struct {
// 		Category string `validate:"newsdataio_cat"`
// 	}

// 	type valCtryStruct struct {
// 		Country string `validate:"newsdataio_ctry"`
// 	}

// 	type valLangStruct struct {
// 		Language string `validate:"newsdataio_lang"`
// 	}

// 	require.NoError(t, val.Var("business", pageform.NEWSDATAIOCatVal.Tag()))
// 	require.NoError(t, val.Struct(valCatStruct{Category: "business"}))
// 	require.Error(t, val.Var("xx", pageform.NEWSDATAIOCatVal.Tag()))
// 	require.Error(t, val.Struct(valCatStruct{Category: "xx"}))

// 	require.NoError(t, val.Var("tw", pageform.NEWSDATAIOCtryVal.Tag()))
// 	require.NoError(t, val.Struct(valCtryStruct{Country: "tw"}))
// 	require.Error(t, val.Var("xx", pageform.NEWSDATAIOCtryVal.Tag()))
// 	require.Error(t, val.Struct(valCtryStruct{Country: "xx"}))

// 	require.NoError(t, val.Var("en", pageform.NEWSDATAIOLangVal.Tag()))
// 	require.NoError(t, val.Struct(valLangStruct{Language: "en"}))
// 	require.Error(t, val.Var("xx", pageform.NEWSDATAIOLangVal.Tag()))
// 	require.Error(t, val.Struct(valLangStruct{Language: "xx"}))
// }
