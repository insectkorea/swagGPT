package api

import (
	"testing"

	"github.com/insectkorea/swagGPT/internal/config"
	"github.com/insectkorea/swagGPT/internal/test"
)

func TestGenerateSwaggerComment(t *testing.T) {
	client := &test.MockOpenAIClient{}
	config := &config.Config{
		Swagger: struct {
			SummaryTemplate     string   `yaml:"summary_template"`
			DescriptionTemplate string   `yaml:"description_template"`
			Tags                []string `yaml:"tags"`
			Accept              string   `yaml:"accept"`
			Produce             string   `yaml:"produce"`
			SuccessResponse     string   `yaml:"success_response"`
			RouterTemplate      string   `yaml:"router_template"`
		}{
			SummaryTemplate:     "Summary for {function_name}",
			DescriptionTemplate: "Description for {function_name}",
			Tags:                []string{"example"},
			Accept:              "json",
			Produce:             "json",
			SuccessResponse:     "200 {string} string \"OK\"",
			RouterTemplate:      "/example/{function_name} [get]",
		},
	}

	comment, err := client.GenerateSwaggerComment("Helloworld", "func %s(g *gin.Context) {", config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if comment == "" {
		t.Fatalf("Expected comment, got none")
	}
}
