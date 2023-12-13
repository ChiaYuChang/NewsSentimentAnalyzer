package openai

import (
	"fmt"
	"sync"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

const (
	API_SCHEME  = "https"
	API_HOST    = "api.openai.com"
	API_VERSION = "v1"
)

var API_URL = fmt.Sprintf("%s://%s/%s", API_SCHEME, API_HOST, API_VERSION)

// API Endpoints
const (
	EPCompletions     string = "completions"
	EPChatCompletions string = "chat/completions"
	EPEmbeddings      string = "embeddings"
)

var Modifier = struct {
	*mold.Transformer
	sync.Once
}{}

func SetModifier(m *mold.Transformer) {
	Modifier.Once.Do(func() {
		Modifier.Transformer = m
	})
}

func GetModifier() *mold.Transformer {
	Modifier.Do(func() {
		Modifier.Transformer = modifiers.New()
	})
	return Modifier.Transformer
}

var Validator = struct {
	*validator.Validate
	sync.Once
}{}

func SetValidator(v *validator.Validate) {
	Validator.Once.Do(func() {
		Validator.Validate = v
	})
}

func GetValidator() *validator.Validate {
	Validator.Do(func() {
		Validator.Validate = validator.New()
	})
	return Validator.Validate
}
