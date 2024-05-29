package handler

import (
	"bytes"
	"fmt"
	"go/ast"
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
	// Read the original file content
	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	handlers, fset, err := scanner.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", filePath, err)
	}

	// Process each handler in parallel
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

	// Wait for all handler processing to complete
	handlerWg.Wait()
	close(handlerResults)

	if len(handlers) == 0 {
		return nil
	}

	// Collect results and sort by start position
	var results []HandlerResult

	for result := range handlerResults {
		if result.Error != nil {
			logrus.Errorf("Error processing handler %s: %v", result.Handler.Name.Name, result.Error)
			continue
		}
		results = append(results, result)
	}

	// Sort the results by start position
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartPos < results[j].StartPos
	})

	// Collect results and update the file content
	var updatedContent bytes.Buffer
	lastPos := 0
	for _, result := range results {
		fn := result.Handler
		comment := result.Comment
		logrus.Infof("Processing handler %s", fn.Name.Name)
		startPos := result.StartPos
		endPos := result.EndPos

		logrus.Infof("Handler %s: startPos=%d, endPos=%d, lastPos=%d\n", fn.Name.Name, startPos, endPos, lastPos)

		// Ensure positions are within bounds
		if startPos < 0 || endPos < 0 || startPos > len(originalContent) || endPos > len(originalContent) {
			return fmt.Errorf("invalid byte positions for handler %s: start %d, end %d", fn.Name.Name, startPos, endPos)
		}

		// Write the content before the handler
		if _, err := updatedContent.Write(originalContent[lastPos:startPos]); err != nil {
			return fmt.Errorf("failed to write file %s: %v", filePath, err)
		}
		// Write the comment
		if _, err := updatedContent.WriteString(fmt.Sprintf("%s\n", comment)); err != nil {
			return fmt.Errorf("failed to write comment %s: %v", filePath, err)
		}
		// Write the handler itself
		handlerContent := originalContent[startPos:endPos]
		if _, err := updatedContent.Write(handlerContent); err != nil {
			return fmt.Errorf("failed to write handler %s: %v", fn.Name.Name, err)
		}
		lastPos = endPos
	}
	// Write the remaining content
	updatedContent.Write(originalContent[lastPos:])

	// Write the updated content back to the file if not in dry run mode
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
