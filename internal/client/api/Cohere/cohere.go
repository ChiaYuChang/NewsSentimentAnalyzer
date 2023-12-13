package cohere

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
)

const (
	API_SCHEME  = "https"
	API_HOST    = "api.cohere.ai"
	API_VERSION = "v1"
	API_METHOD  = http.MethodPost
)

var API_URL = fmt.Sprintf("%s://%s/%s", API_SCHEME, API_HOST, API_VERSION)

const (
	EPCoEmbed  = "embed"
	EPChat     = "chat"
	EPGenerate = "generate"
)

const (
	InputTypeSearchDocument = "search_document"
	InputTypeSearchQuery    = "search_query"
	InputTypeClassification = "classification"
	InputTypeClustering     = "clustering"
)

const (
	TruncateNone  = "NONE"
	TruncateStart = "START"
	TruncateEnd   = "END"
)

const (
	EmbedModelEnglishv3           = "embed-english-v3.0"
	EmbedModelMultilingualv3      = "embed-multilingual-v3.0"
	EmbedModelEnglishLightv3      = "embed-english-light-v3.0"
	EmbedModelMultilingualLightv3 = "embed-multilingual-light-v3.0"
	EmbedModelEnglishv2           = "embed-english-v2.0"
	EmbedModelEnglishLightv2      = "embed-english-light-v2.0"
	EmbedModelMultilingualv2      = "embed-multilingual-v2.0"
)

const (
	GenerateModelCommand             = "command"
	GenerateModelCommandNightly      = "command-nightly"
	GenerateModelCommandLight        = "command-light"
	GenerateModelCommandLightNightly = "command-light-nightly"
)

func init() {
	EnmusEmbdedModel := validator.NewEnmus[string](
		"cohere_embed_model",
		EmbedModelEnglishv3,
		EmbedModelMultilingualv3,
		EmbedModelEnglishLightv3,
		EmbedModelMultilingualLightv3,
		EmbedModelEnglishv2,
		EmbedModelEnglishLightv2,
		EmbedModelMultilingualv2,
	)

	EnmusGenerateModel := validator.NewEnmus[string](
		"cohere_generate_model",
		GenerateModelCommand,
		GenerateModelCommandNightly,
		GenerateModelCommandLight,
		GenerateModelCommandLightNightly,
	)

	EnmusInputType := validator.NewEnmus[string](
		"cohere_embed_input_type",
		InputTypeSearchDocument,
		InputTypeSearchQuery,
		InputTypeClassification,
		InputTypeClustering,
	)

	EnmuTruncate := validator.NewEnmus[string](
		"cohere_truncate",
		TruncateNone,
		TruncateStart,
		TruncateEnd,
	)
	val, err := validator.GetDefaultValidate()
	if err != nil {
		fmt.Printf("error while GetDefaultValidate: %v", err)
		os.Exit(1)
	}

	err = validator.RegisterValidator(
		val,
		EnmusEmbdedModel,
		EnmusInputType,
		EnmuTruncate,
		EnmusGenerateModel,
	)
	if err != nil {
		fmt.Printf("error while RegisterValidator: %v", err)
		os.Exit(1)
	}
}
