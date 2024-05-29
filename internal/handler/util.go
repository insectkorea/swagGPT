package handler

import (
	"bytes"
	"go/format"
	"go/token"

	"github.com/insectkorea/swagGPT/internal/api"
	"github.com/insectkorea/swagGPT/internal/scanner"
	"github.com/sirupsen/logrus"
)

// EstimateTotalTokens estimates the total number of tokens for all handlers in the given files.
func EstimateTotalTokens(files []string) int {
	totalTokens := 0
	for _, file := range files {
		handlers, _, err := scanner.ParseFile(file)
		if err != nil {
			logrus.Errorf("Error parsing file %s: %v", file, err)
			continue
		}
		for _, handler := range handlers {
			var buf bytes.Buffer
			if err := format.Node(&buf, token.NewFileSet(), handler); err != nil {
				logrus.Errorf("Failed to format handler %s: %v", handler.Name.Name, err)
				continue
			}
			handlerContent := buf.String()
			totalTokens += api.EstimateTokens(handlerContent)
		}
	}
	return totalTokens
}
