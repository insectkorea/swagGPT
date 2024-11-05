package api

import (
	"testing"

	"github.com/insectkorea/swagGPT/internal/test"
)

func TestGenerateSwaggerComment(t *testing.T) {
	client := &test.MockOpenAIClient{}

	comment, err := client.GenerateSwaggerComment("Helloworld", "func %s(g *gin.Context) {", "test", "/api/v1/organizations/:organization_id/bundles [get]")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if comment == "" {
		t.Fatalf("Expected comment, got none")
	}
}
