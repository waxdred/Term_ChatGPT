package main

import (
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func requetOpenAI(chatGpt *ChatGpt, curr Session, input string) {
	chatGpt.Lock()
	chatGpt.routine = true
	inp := curr.Msg
	chatGpt.Unlock()
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
		chatGpt.Lock()
		chatGpt.rep = err.Error()
		chatGpt.Unlock()
		return
	}
	if len(resp.Choices) > 0 {
		chatGpt.Lock()
		ret := strings.TrimPrefix(resp.Choices[0].Text, "\n")
		chatGpt.history = append(chatGpt.history, ret)
		if ret == "" {
			chatGpt.rep = "?"
		} else {
			chatGpt.rep = ret
		}
		chatGpt.Unlock()
	}
}
