package test

import "github.com/insectkorea/swagGPT/internal/config"

// MockOpenAIClient is a mock implementation of the Client interface.
type MockOpenAIClient struct{}

func (m *MockOpenAIClient) GenerateSwaggerComment(functionName, functionContent string, config *config.Config) (string, error) {
	return `// ` + functionName + ` godoc
// @Summary ` + functionName + ` summary
// @Description do ` + functionName + `
`, nil
}
