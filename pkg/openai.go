package pkg

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"os"
)

type OpenAI struct {
	Client *openai.Client
}

func MakeOpenAI() *OpenAI {
	return &OpenAI{
		Client: openai.NewClient(os.Getenv("OPENAI_AUTH_TOKEN")),
	}
}

func (o *OpenAI) CompleteChat(prompt string) (string, error) {
	resp, err := o.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	log.Println(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil
}
