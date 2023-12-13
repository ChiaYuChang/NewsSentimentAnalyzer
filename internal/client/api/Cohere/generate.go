package cohere

type GenerateRequestBody struct {
	Prompt            string             `json:"prompt"                                   validate:"required"`
	Model             string             `json:"model"              mod:"default=command" validate:"required,cohere_generate_model"`
	NumGenerations    int                `json:"num_generations"    mod:"default=1"       validate:"min=1,max=5"`
	Steam             bool               `json:"steam,omitempty"`
	MaxTokens         int                `json:"max_tokens"         mod:"default=256"`
	Truncate          string             `json:"truncate"           mod:"default=END"     validate:"cohere_truncate"`
	Temperature       float32            `json:"temperature"        mod:"default=0.75"    validate:"min=0.0,max=5.0"`
	Preset            string             `json:"preset,omitempty"`
	EndSequences      []string           `json:"end_sequences,omitempty"`
	StopSequences     []string           `json:"stop_sequences,omitempty"`
	K                 int                `json:"k,omitempty"                              validate:"min=0,max=500"`
	P                 float32            `json:"p,omitempty"                              validate:"min=0.01,max=0.99"`
	FrequencyPenalty  float32            `json:"frequency_penalty,omitempty"`
	PresencePenalty   float32            `json:"presence_penalty,omitempty"               validate:"min=0.0,max=1.0"`
	ReturnLikelihoods string             `json:"return_likelihoods" mod:"default:NONE"    validate:"oneof=NONE GENERATION ALL"`
	LogitBias         map[string]float32 `json:"logit_bias,omitempty"                     validate:"dive,min=-10.0,max=10.0"`
}

func (body GenerateRequestBody) Endpoint() string {
	return EPGenerate
}
