package backend

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"os"
)

var client *openai.Client

func CallAI(model, prompt string) (openai.ChatCompletionResponse, error) {
	return client.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Answer the following question by choosing one or many answers, return only answer number"},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
}

func InitializeAPI() {
	key := os.Getenv("OPENAI_KEY")
	if key == "" {
		panic("OPENAI_KEY is not set")
	}

	client = openai.NewClient(key)
}
