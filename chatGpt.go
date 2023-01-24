package main

import (
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func requetOpenAI(chatGpt *ChatGpt, input string) (string, error) {
	var ret string
	chatGpt.history = append(chatGpt.history, input)
	inp := strings.Join(chatGpt.history, " ")
	req := gogpt.CompletionRequest{
		Model:            chatGpt.Model,
		MaxTokens:        int(chatGpt.MaxTokens),
		Prompt:           inp,
		Temperature:      float32(chatGpt.Temperature),
		TopP:             float32(chatGpt.TopP),
		FrequencyPenalty: float32(chatGpt.FrequencyPenalty),
		PresencePenalty:  float32(chatGpt.PresencePenalty),
	}
	resp, err := chatGpt.c.CreateCompletion(chatGpt.ctx, req)
	if err != nil {
		return ret, err
	}
	if len(resp.Choices) > 0 {
		ret = strings.TrimPrefix(resp.Choices[0].Text, "\n")
		chatGpt.history = append(chatGpt.history, ret)
	}
	return ret, nil
}
