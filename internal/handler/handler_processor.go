package handler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"

	"github.com/insectkorea/swagGPT/internal/api"
)

// processHandler processes a single handler to generate a Swagger comment.
func processHandler(handler *ast.FuncDecl, client api.Client, model string) (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), handler); err != nil {
		return "", fmt.Errorf("failed to format handler %s: %v", handler.Name.Name, err)
	}
	handlerContent := buf.String()

	comment, err := client.GenerateSwaggerComment(handler.Name.Name, handlerContent, model)
	if err != nil {
		return "", fmt.Errorf("failed to generate comment for %s: %v", handler.Name.Name, err)
	}

	return comment, nil
}
