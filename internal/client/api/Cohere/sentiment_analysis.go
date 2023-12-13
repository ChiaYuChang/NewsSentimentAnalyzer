package cohere

import "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"

const SentimentAnalysisPrompt = `As an AI specializing in language and emotion analysis, your task is to assess the sentiments conveyed in an article. The given article will be enclosed with the symbols [^] and [$]. Consider the overall tone, emotional nuances, and context within the statements. Classify the article into one of three categories: -1 negative, 0 for neutral, and 1 for positive. Please present your responses with a single number. Respond sequentially without repeating the provided sentences. For example, if the article is [^]I love this movie[$], your corresponding response should be 1.`

func NewSentimentAnalysisRequest(apikey, text string) Request[ChatRequestBody] {
	req := NewChatRequest(apikey, text)
	req.Body.AppendChatHistory(ChatRoleChatBot, SentimentAnalysisPrompt, "").
		SetTemperature(0.001).
		SetCitationQuality(CitationQualityAccurate)
	return req
}

type SentimentAnalysisObject ChatResponseBody

func (obj ChatResponseBody) Content() (int, error) {
	return convert.StrTo(obj.Text).Int()
}
