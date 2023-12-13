package cohere

import "github.com/google/uuid"

func NewChatRequest(apikey string, message string) Request[ChatRequestBody] {
	return Request[ChatRequestBody]{
		apikey: apikey,
		Body: ChatRequestBody{
			Message: message,
		},
	}
}

const (
	ChatRoleChatBot = "CHATBOT"
	ChatRoleUser    = "USER"
)

const (
	CitationQualityFast     = "fast"
	CitationQualityAccurate = "accurate"
)

type ChatRequestBody struct {
	Message          string              `json:"message"                                    validate:"required"`
	Model            string              `json:"model"               mod:"default=command"  validate:"cohere_generate_model"`
	Stream           bool                `json:"stream,omitempty"`
	PreambleOverride string              `json:"preamble_override,omitempty"`
	ChatHistory      []ChatHistory       `json:"chat_history,omitempty"                     validate:"dive"`
	ConversationId   string              `json:"conversation_id,omitempty"`
	PromptTruncation string              `json:"prompt_truncation"   mod:"default=OFF"      validate:"oneof=OFF AUTO"`
	Connectors       []Connector         `json:"connectors,omitempty"                       validate:"dive"`
	SearchQueryOnly  bool                `json:"search_query_only,omitempty"`
	Documents        []map[string]string `json:"documents,omitempty"`
	CitationQuality  string              `json:"citation_quality"    mod:"default=accurate" validate:"oneof=fast accurate"`
	Temperature      float32             `json:"temperature"         mod:"default=0.3"      validate:"gt=0.0"`
}

func (body ChatRequestBody) Endpoint() string {
	return EPChat
}

// role and message should not be empty
func (body *ChatRequestBody) AppendChatHistory(role, message, username string) *ChatRequestBody {
	body.ChatHistory = append(body.ChatHistory, ChatHistory{
		Role: role, Message: message, UserName: username,
	})
	return body
}

// A list of relevant documents that the model can use to enrich its reply.
// see https://docs.cohere.com/docs/retrieval-augmented-generation-rag#document-mode
func (body *ChatRequestBody) AppendDocument(document map[string]string) *ChatRequestBody {
	body.Documents = append(body.Documents, document)
	return body
}

// Current only support id = "web-search", option value are domains to restrict to.
// PromptTruncation will be set to "AUTO" if it has not yet set.
func (body *ChatRequestBody) AppendConnector(id string, option map[string]string) *ChatRequestBody {
	body.Connectors = append(body.Connectors, Connector{
		Id:     id,
		Option: option,
	})

	if body.PromptTruncation == "" {
		body.SetPromptTruncation("AUTO")
	}
	return body
}

// should be GenerateModelCommand ("command"), GenerateModelCommandNightly ("command-nightly"),
// GenerateModelCommandLight ("command-light"), or GenerateModelCommandLightNightly ("command-light-nightly")
func (body *ChatRequestBody) SetModel(model string) *ChatRequestBody {
	body.Model = model
	return body
}

// should be "OFF" or "AUTO"
func (body *ChatRequestBody) SetPromptTruncation(truncation string) *ChatRequestBody {
	body.PromptTruncation = truncation
	return body
}

// should be "fast" or "accurate"
func (body *ChatRequestBody) SetCitationQuality(quality string) *ChatRequestBody {
	body.CitationQuality = quality
	return body
}

// default to 0.3, should be a non-negative float32
func (body *ChatRequestBody) SetTemperature(temperature float32) *ChatRequestBody {
	body.Temperature = temperature
	return body
}

type ChatHistory struct {
	Role     string `json:"role"    validate:"required,oneof=CHATBOT USER"`
	Message  string `json:"message" validate:"required"`
	UserName string `json:"username,omitempty"`
}

type Connector struct {
	Id     string            `json:"id"      validate:"required,oneof=web-search"`
	Option map[string]string `json:"options" validate:"dive,required"`
}

type ChatResponseBody struct {
	ResponseId   uuid.UUID `json:"response_id"`
	Text         string    `json:"text"`
	GenerationId uuid.UUID `json:"generation_id"`
	Citations    []struct {
		Start       int      `json:"start"`
		End         int      `json:"end"`
		Text        string   `json:"text"`
		DocumentIds []string `json:"document_ids"`
	} `json:"citations"`
	Documents     []map[string]string `json:"documents"`
	SearchQueries []SearchQuery       `json:"search_queries"`
	SearchResults []struct {
		SearchQuery SearchQuery `json:"search_query"`
		Connector   Connector   `json:"connector"`
		DocumentIds []string    `json:"document_ids"`
	} `json:"search_results"`
	TokenCount TokenCount `json:"token_count"`
	Meta       Meta       `json:"meta"`
}

type TokenCount struct {
	PromptTokens   int `json:"prompt_tokens"`
	ResponseTokens int `json:"response_tokens"`
	TotalTokens    int `json:"total_tokens"`
	BilledTokens   int `json:"billed_tokens"`
}

type SearchQuery struct {
	Text         string `json:"text"`
	GenerationId string `json:"generation_id"`
}
