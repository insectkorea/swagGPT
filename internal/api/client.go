package api

import (
	"context"
	"fmt"

	"github.com/insectkorea/swagGPT/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

// Client is an interface representing the OpenAI client.
type Client interface {
	GenerateSwaggerComment(functionName, functionContent string, config *config.Config) (string, error)
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
func (c *OpenAIClient) GenerateSwaggerComment(functionName, functionContent string, config *config.Config) (string, error) {
	prompt := fmt.Sprintf(`
Generate Swagger comments for the following function:
%s
Summary: %s
Description: %s
Tags: %v
Accept: %s
Produce: %s
Success Response: %s
Router: %s
`,
		functionContent,
		config.Swagger.SummaryTemplate,
		config.Swagger.DescriptionTemplate,
		config.Swagger.Tags,
		config.Swagger.Accept,
		config.Swagger.Produce,
		config.Swagger.SuccessResponse,
		config.Swagger.RouterTemplate,
	)

	req := openai.CompletionRequest{
		Model:     "gpt-3.5-turbo-instruct",
		Prompt:    prompt,
		MaxTokens: 150,
	}

	resp, err := c.client.CreateCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}
