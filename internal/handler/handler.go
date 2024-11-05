package handler

import (
	"go/ast"
	"sync"

	"github.com/insectkorea/swagGPT/internal/api"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

// FileHandlerPair is a struct that holds a file path and a handler function.
type FileHandlerPair struct {
	FilePath string
	Handler  *ast.FuncDecl
}

// ProcessFiles processes the given files to add Swagger comments to handler functions.
func ProcessFiles(files []string, client api.Client, dryRun bool, model string, contextFilePath string) error {
	bar := progressbar.Default(int64(len(files)))

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			// nolint:errcheck
			defer bar.Add(1)
			if err := processFile(filename, client, dryRun, model, contextFilePath); err != nil {
				logrus.Errorf("Error processing file %s: %v", filename, err)
			}
		}(file)
	}

	wg.Wait()
	return nil
}
