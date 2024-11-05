package handler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"

	"github.com/insectkorea/swagGPT/internal/api"
	"github.com/insectkorea/swagGPT/internal/matcher"
	"github.com/insectkorea/swagGPT/internal/model"
)

// processHandler processes a single handler to generate a Swagger comment.
func processHandler(handler *ast.FuncDecl, client api.Client, model string, routes []model.Route) (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), handler); err != nil {
		return "", fmt.Errorf("failed to format handler %s: %v", handler.Name.Name, err)
	}
	handlerContent := buf.String()

	routeString, err := matcher.MatchHandlerToRoute(handler.Name.Name, routes)
	if err != nil {
		return "", fmt.Errorf("failed to match handler to route for %s: %v", handler.Name.Name, err)
	}

	comment, err := client.GenerateSwaggerComment(handler.Name.Name, handlerContent, model, routeString)
	if err != nil {
		return "", fmt.Errorf("failed to generate comment for %s: %v", handler.Name.Name, err)
	}
	fmt.Println(fmt.Sprintf("comment: %s", comment))

	return comment, nil
}
