package openai

import "fmt"

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
