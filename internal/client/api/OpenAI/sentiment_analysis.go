package openai

import "encoding/json"

const SentimentAnalysisPrompt = `As an AI specializing in language and emotion analysis, your task is to assess the sentiments conveyed in a set of statements. Each statement will be enclosed with the symbols [^] and [$]. Consider the overall tone, emotional nuances, and context within the statements. Classify each statement into one of three categories: -1 for negative, 0 for neutral, and 1 for positive. Please present your responses in a JSON list. Respond sequentially without repeating the provided sentences. For instance, if the sentence is [^]I love this movie[$] [^]I hate you[$], your corresponding response should be [1, -1].`

// See https://platform.openai.com/docs/api-reference/chat
func NewSentimentAnalysisRequest(apikey string, text string) Request[ChatCompletionsRequestBody] {
	req := Request[ChatCompletionsRequestBody]{
		Body:   ChatCompletionsRequestBody{},
		apikey: apikey,
	}
	req.Body.
		AppendSystemMessages(SentimentAnalysisPrompt, "").
		AppendUserMessages(text, "")
	return req
}

type SentimentAnalysisObject ChatCompletionsObject

func (obj SentimentAnalysisObject) Content() ([][]int, error) {
	content := make([][]int, len(obj.Choices))
	for i := range obj.Choices {
		err := json.Unmarshal([]byte(obj.Choices[i].Message.Content), &content[i])
		if err != nil {
			return nil, err
		}
	}
	return content, nil
}
