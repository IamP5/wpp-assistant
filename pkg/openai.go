package pkg

import (
	"context"
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
	log.Println("[OpenAI] - started CompleteChat method")

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
		log.Printf("[OpenAI] - ChatCompletion error: %s", err.Error())
		return "", err
	}

	log.Println("[OpenAI] - CompleteChat with success!")
	return resp.Choices[0].Message.Content, nil
}
