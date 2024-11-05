package handler

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/insectkorea/swagGPT/internal/model"
)

type ContextFileHandler struct {
}

func (ch *ContextFileHandler) ExtractRoutes(contextFilePath string) ([]model.Route, error) {
	if contextFilePath == "" {
		return nil, nil
	}

	// Ensure the file exists
	if _, err := os.Stat(contextFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("context file does not exist: %s", contextFilePath)
	}

	// Create a new token file set
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, contextFilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	// Initialize the RouteExtractor
	extractor := &RouteExtractor{
		Routes:     []model.Route{},
		groupStack: []string{},
	}

	// Traverse the AST with the extractor
	ast.Walk(extractor, node)

	// Output the extracted routes
	for _, route := range extractor.Routes {
		println("Method:", route.Method, "Path:", route.Path, "Pattern:", route.Pattern)
	}

	return extractor.Routes, nil
}

type RouteExtractor struct {
	Routes     []model.Route
	groupStack []string
}

func (re *RouteExtractor) Visit(node ast.Node) ast.Visitor {
	switch x := node.(type) {
	case *ast.CallExpr:
		selExpr, ok := x.Fun.(*ast.SelectorExpr)
		if !ok {
			return re
		}

		methodName := selExpr.Sel.Name

		// List of HTTP methods to look for
		httpMethods := map[string]bool{
			"GET":     true,
			"POST":    true,
			"PUT":     true,
			"DELETE":  true,
			"PATCH":   true,
			"OPTIONS": true,
			"HEAD":    true,
		}

		// Check if it's a route registration method
		if httpMethods[methodName] {
			// Ensure there is at least one argument for the path
			if len(x.Args) < 1 {
				return re
			}

			// Extract the path argument
			pathArg, ok := x.Args[0].(*ast.BasicLit)
			if !ok || pathArg.Kind != token.STRING {
				return re
			}

			// Remove quotes from the path string
			path := strings.Trim(pathArg.Value, `"'`)

			// Build the full path considering the group prefixes
			fullPath := ""
			if len(re.groupStack) > 0 {
				fullPath = strings.Join(re.groupStack, "") + path
			} else {
				fullPath = path
			}

			// Append the route to the slice
			re.Routes = append(re.Routes, model.Route{
				Method:  methodName,
				Path:    path,
				Pattern: fullPath,
			})
		}

		// Handle Group() calls to track prefixes
		if selExpr.Sel.Name == "Group" && len(x.Args) >= 1 {
			groupPathArg, ok := x.Args[0].(*ast.BasicLit)
			if !ok || groupPathArg.Kind != token.STRING {
				return re
			}

			groupPath := strings.Trim(groupPathArg.Value, `"'`)

			// Push the group path to the stack
			re.groupStack = append(re.groupStack, groupPath)

			// After the group is processed, pop the stack
			// This simplistic approach assumes no concurrent groups
			defer func() {
				if len(re.groupStack) > 0 {
					re.groupStack = re.groupStack[:len(re.groupStack)-1]
				}
			}()
		}
	}
	return re
}
