package llm

import (
	"context"
	"errors"
	"log"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	oai            *openai.Client
	chatModel      string
	embeddingModel string
}

func New(apiKey string, chatModel string, embeddingModel string) *Client {
	return &Client{
		oai:            openai.NewClient(apiKey),
		chatModel:      chatModel,
		embeddingModel: embeddingModel,
	}
}

func (c *Client) Chat(ctx context.Context, system string, user string) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := c.oai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: c.chatModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: system,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: user,
				},
			},
			MaxTokens:   1024,
			Temperature: 0.3,
		})

		if err == nil {
			if len(resp.Choices) == 0 {
				return "", errors.New("empty chat completion response")
			}

			return resp.Choices[0].Message.Content, nil
		}

		lastErr = err
		log.Printf("WARN: LLM attempt %d failed: %v", attempt, err)
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}

	return "", errors.New("LLM unavailable after retries: " + lastErr.Error())
}

func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
	resp, err := c.oai.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(c.embeddingModel),
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("empty embedding response")
	}

	return resp.Data[0].Embedding, nil
}