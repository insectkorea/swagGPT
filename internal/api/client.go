package api

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// Client is an interface representing the OpenAI client.
type Client interface {
	GenerateSwaggerComment(functionName, functionContent, model string, routeString string) (string, error)
}

// OpenAIClient is a struct that implements the Client interface.
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient initializes and returns an OpenAI client.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
	}
}

// GenerateSwaggerComment generates Swagger comments using OpenAI API.
func (c *OpenAIClient) GenerateSwaggerComment(functionName, functionContent, model string, routeString string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role: "user",
			Content: fmt.Sprintf(userPromptTemplate,
				functionContent,
				routeString,
			),
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	resp, err := c.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
