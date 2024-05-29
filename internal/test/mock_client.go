package test

// MockOpenAIClient is a mock implementation of the Client interface.
type MockOpenAIClient struct{}

func (m *MockOpenAIClient) GenerateSwaggerComment(functionName, functionContent, model string) (string, error) {
	return `// ` + functionName + ` godoc
// @Summary ` + functionName + ` summary
// @Description do ` + functionName + `
// @Success 200 {string} string "OK"
// @Router /example/` + functionName + ` [get]
`, nil
}
