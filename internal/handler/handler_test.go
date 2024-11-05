package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/insectkorea/swagGPT/internal/scanner"
	"github.com/insectkorea/swagGPT/internal/test"
)

func TestProcessFiles(t *testing.T) {
	// Setup temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "swagger-comment-adder")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test Go file with a handler function
	goFileContent := `
package example

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Helloworld handler function
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
`
	goFilePath := filepath.Join(tmpDir, "example.go")
	if err := os.WriteFile(goFilePath, []byte(goFileContent), 0644); err != nil {
		t.Fatalf("Failed to write test Go file: %v", err)
	}

	files, err := scanner.ScanDir(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	mockClient := &test.MockOpenAIClient{}

	err = ProcessFiles(files, mockClient, false, "test-model", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the contents of the modified file
	modifiedContent, err := os.ReadFile(goFilePath)
	if err != nil {
		t.Fatalf("Failed to read modified Go file: %v", err)
	}

	expectedComment := `// Helloworld godoc
// @Summary Helloworld summary
// @Description do Helloworld
// @Success 200 {string} string "OK"`

	if !strings.Contains(string(modifiedContent), expectedComment) {
		t.Fatalf("Expected comment not found in modified file:\n%s", string(modifiedContent))
	}
}
