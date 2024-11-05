package handler

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/insectkorea/swagGPT/internal/model"
	"github.com/insectkorea/swagGPT/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestProcessFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.go")
	originalContent := `package main

import "github.com/gin-gonic/gin"

// TestHandler handles a test request
func TestHandler(c *gin.Context) {
	c.JSON(200, "Hello World")
}
`
	err := os.WriteFile(filePath, []byte(originalContent), 0644)
	assert.NoError(t, err)

	client := &test.MockOpenAIClient{}
	err = processFile(filePath, client, true, "test-model", "")
	assert.NoError(t, err)
}

func TestReadFileAndParse(t *testing.T) {
	filePath := createTempGoFile(t, `package main

import "github.com/gin-gonic/gin"

// TestHandler handles a test request
func TestHandler(c *gin.Context) {
	c.JSON(200, "Hello World")
}
`)

	originalContent, handlers, fset, err := readFileAndParse(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, originalContent)
	assert.NotNil(t, handlers)
	assert.NotNil(t, fset)
	assert.Equal(t, 1, len(handlers))
}

func TestProcessHandlers(t *testing.T) {
	handlers := parseHandlersFromContent(t, `package main

import "github.com/gin-gonic/gin"

// TestHandler handles a test request
func TestHandler(c *gin.Context) {
	c.JSON(200, "Hello World")
}
`)

	client := &test.MockOpenAIClient{}
	results, err := processHandlers(handlers, client, "test-model", token.NewFileSet(), []model.Route{
		{
			Path:    "/example/TestHandler",
			Method:  "GET",
			Pattern: "/example/TestHandler",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "// TestHandler godoc\n// @Summary TestHandler summary\n// @Description do TestHandler\n// @Success 200 {string} string \"OK\"\n// @Router /example/TestHandler [get]\n", results[0].Comment)
}

func TestCollectAndSortResults(t *testing.T) {
	handlerResults := make(chan HandlerResult, 2)
	handlerResults <- HandlerResult{StartPos: 10, Comment: "// Comment 1"}
	handlerResults <- HandlerResult{StartPos: 5, Comment: "// Comment 2"}
	close(handlerResults)

	results := collectAndSortResults(handlerResults)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, "// Comment 2", results[0].Comment)
	assert.Equal(t, "// Comment 1", results[1].Comment)
}

func TestUpdateFileContent(t *testing.T) {
	originalContent := []byte(`package main

import "github.com/gin-gonic/gin"

// TestHandler handles a test request
func TestHandler(c *gin.Context) {
	c.JSON(200, "Hello World")
}
`)
	results := []HandlerResult{
		{
			Handler:  &ast.FuncDecl{Name: &ast.Ident{Name: "TestHandler"}},
			StartPos: 56,
			EndPos:   120,
			Comment:  "// TestHandler godoc\n// @Summary TestHandler summary\n// @Description do TestHandler\n// @Success 200 {string} string \"OK\"\n// @Router /example/TestHandler [get]\n",
		},
	}

	err := updateFileContent("test.go", originalContent, results, true)
	assert.NoError(t, err)
}

func TestWriteHandlerContent(t *testing.T) {
	var updatedContent bytes.Buffer
	originalContent := []byte(`package main

import "github.com/gin-gonic/gin"

// TestHandler handles a test request
func TestHandler(c *gin.Context) {
	c.JSON(200, "Hello World")
}
`)
	result := HandlerResult{
		Handler:  &ast.FuncDecl{Name: &ast.Ident{Name: "TestHandler"}},
		Comment:  "// Swagger Comment",
		StartPos: 56,
		EndPos:   120,
	}
	lastPos := 0

	err := writeHandlerContent(&updatedContent, "test.go", originalContent, result, &lastPos)
	assert.NoError(t, err)
	assert.Contains(t, updatedContent.String(), "// Swagger Comment")
	assert.Equal(t, 120, lastPos)
}

// Helper functions
func createTempGoFile(t *testing.T, content string) string {
	t.Helper()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.go")
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)
	return filePath
}

func parseHandlersFromContent(t *testing.T, content string) []*ast.FuncDecl {
	t.Helper()
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	assert.NoError(t, err)

	var handlers []*ast.FuncDecl
	for _, decl := range node.Decls {
		if fn, isFn := decl.(*ast.FuncDecl); isFn {
			handlers = append(handlers, fn)
		}
	}
	return handlers
}
