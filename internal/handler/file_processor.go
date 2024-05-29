package handler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"sort"
	"sync"

	"github.com/insectkorea/swagGPT/internal/api"
	"github.com/insectkorea/swagGPT/internal/scanner"

	"github.com/sirupsen/logrus"
)

// HandlerResult holds the result of processing a single handler.
type HandlerResult struct {
	Handler  *ast.FuncDecl
	Comment  string
	Error    error
	StartPos int
	EndPos   int
}

// processFile processes a single file to add Swagger comments to its handler functions.
func processFile(filePath string, client api.Client, dryRun bool, model string) error {
	originalContent, handlers, fset, err := readFileAndParse(filePath)
	if err != nil {
		return err
	}

	handlerResults, err := processHandlers(handlers, client, model, fset)
	if err != nil {
		return err
	}

	return updateFileContent(filePath, originalContent, handlerResults, dryRun)
}

func readFileAndParse(filePath string) ([]byte, []*ast.FuncDecl, *token.FileSet, error) {
	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	handlers, fset, err := scanner.ParseFile(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse file %s: %v", filePath, err)
	}

	return originalContent, handlers, fset, nil
}

func processHandlers(handlers []*ast.FuncDecl, client api.Client, model string, fset *token.FileSet) ([]HandlerResult, error) {
	var handlerWg sync.WaitGroup
	handlerResults := make(chan HandlerResult, len(handlers))

	for _, handler := range handlers {
		handlerWg.Add(1)
		go func(handler *ast.FuncDecl) {
			defer handlerWg.Done()
			comment, err := processHandler(handler, client, model)
			startPos := fset.Position(handler.Pos()).Offset
			endPos := fset.Position(handler.End()).Offset
			handlerResults <- HandlerResult{Handler: handler, Comment: comment, Error: err, StartPos: startPos, EndPos: endPos}
		}(handler)
	}

	handlerWg.Wait()
	close(handlerResults)

	if len(handlers) == 0 {
		return nil, nil
	}

	return collectAndSortResults(handlerResults), nil
}

func collectAndSortResults(handlerResults chan HandlerResult) []HandlerResult {
	var results []HandlerResult

	for result := range handlerResults {
		if result.Error != nil {
			logrus.Errorf("Error processing handler %s: %v", result.Handler.Name.Name, result.Error)
			continue
		}
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].StartPos < results[j].StartPos
	})

	return results
}

func updateFileContent(filePath string, originalContent []byte, results []HandlerResult, dryRun bool) error {
	var updatedContent bytes.Buffer
	lastPos := 0
	for _, result := range results {
		err := writeHandlerContent(&updatedContent, filePath, originalContent, result, &lastPos)
		if err != nil {
			return err
		}
	}

	updatedContent.Write(originalContent[lastPos:])

	if !dryRun {
		logrus.Infof("Writing updated content to file %s", filePath)
		if err := os.WriteFile(filePath, updatedContent.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", filePath, err)
		}
	} else {
		logrus.Infof("Dry Run: Would update the file %s with \n\n%s", filePath, updatedContent.String())
	}

	return nil
}

func writeHandlerContent(updatedContent *bytes.Buffer, filePath string, originalContent []byte, result HandlerResult, lastPos *int) error {
	fn := result.Handler
	comment := result.Comment
	logrus.Infof("Processing handler %s", fn.Name.Name)
	startPos := result.StartPos
	endPos := result.EndPos

	logrus.Infof("Handler %s: startPos=%d, endPos=%d, lastPos=%d\n", fn.Name.Name, startPos, endPos, *lastPos)

	if startPos < 0 || endPos < 0 || startPos > len(originalContent) || endPos > len(originalContent) {
		return fmt.Errorf("invalid byte positions for handler %s: start %d, end %d", fn.Name.Name, startPos, endPos)
	}

	if _, err := updatedContent.Write(originalContent[*lastPos:startPos]); err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}

	if _, err := updatedContent.WriteString(fmt.Sprintf("%s\n", comment)); err != nil {
		return fmt.Errorf("failed to write comment %s: %v", filePath, err)
	}

	handlerContent := originalContent[startPos:endPos]
	if _, err := updatedContent.Write(handlerContent); err != nil {
		return fmt.Errorf("failed to write handler %s: %v", fn.Name.Name, err)
	}

	*lastPos = endPos

	return nil
}
