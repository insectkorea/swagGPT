package handler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"sync"

	"github.com/insectkorea/swagGPT/internal/api"
	"github.com/insectkorea/swagGPT/internal/config"
	"github.com/insectkorea/swagGPT/internal/scanner"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

// ProcessFiles processes the given files to add Swagger comments to handler functions.
func ProcessFiles(files []string, client api.Client, dryRun bool, configPath string) error {
	config, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	bar := progressbar.Default(int64(len(files)))
	handlerChannel := make(chan *ast.FuncDecl)

	var wg sync.WaitGroup

	processFilesConcurrently(files, &wg, handlerChannel, bar)

	go func() {
		wg.Wait()
		close(handlerChannel)
	}()

	processHandlersConcurrently(handlerChannel, client, config, dryRun)

	return nil
}

// processFilesConcurrently processes files concurrently and sends handlers to a channel.
func processFilesConcurrently(files []string, wg *sync.WaitGroup, handlerChannel chan<- *ast.FuncDecl, bar *progressbar.ProgressBar) {
	for _, file := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			defer bar.Add(1)
			logrus.Infof("Processing file: %s", filename)
			handlers, err := scanner.ParseFile(filename)
			if err != nil {
				logrus.Errorf("Error parsing file %s: %v", filename, err)
				return
			}
			for _, handler := range handlers {
				handlerChannel <- handler
			}
		}(file)
	}
}

// processHandlersConcurrently processes handlers concurrently and generates Swagger comments.
func processHandlersConcurrently(handlerChannel <-chan *ast.FuncDecl, client api.Client, config *config.Config, dryRun bool) {
	var wg sync.WaitGroup
	for handler := range handlerChannel {
		wg.Add(1)
		go func(handler *ast.FuncDecl) {
			defer wg.Done()
			processHandler(handler, client, config, dryRun)
		}(handler)
	}
	wg.Wait()
}

// processHandler generates and inserts Swagger comments for a single handler.
func processHandler(handler *ast.FuncDecl, client api.Client, config *config.Config, dryRun bool) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), handler); err != nil {
		logrus.Errorf("Failed to format handler %s: %v", handler.Name.Name, err)
		return
	}
	handlerContent := buf.String()
	comment, err := client.GenerateSwaggerComment(handler.Name.Name, handlerContent, config)
	if err != nil {
		logrus.Errorf("Failed to generate comment for %s: %v", handler.Name.Name, err)
		return
	}

	newContent := fmt.Sprintf("// %s\n%s", handler.Name.Name, comment)

	if dryRun {
		logrus.Infof("Dry Run: Would insert comment for %s in file %s: \n%s", handler.Name.Name, handler.Name.Name+".go", comment)
	} else {
		if err := backupAndWriteFile(handler.Name.Name+".go", newContent); err != nil {
			logrus.Errorf("Failed to write file %s: %v", handler.Name.Name+".go", err)
		}
	}
}

// backupAndWriteFile creates a backup and writes new content to the file.
func backupAndWriteFile(filePath, newContent string) error {
	backupFilePath := filePath + ".bak"
	if err := backupFile(filePath, backupFilePath); err != nil {
		return fmt.Errorf("failed to backup file %s: %v", filePath, err)
	}
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// RestoreFiles restores original files from backups.
func RestoreFiles(dir string) error {
	files, err := scanner.ScanDir(dir)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file) == ".bak" {
			originalFile := file[:len(file)-len(".bak")]
			if err := os.Rename(file, originalFile); err != nil {
				logrus.Errorf("Failed to restore file %s: %v", originalFile, err)
			}
		}
	}

	return nil
}

// backupFile creates a backup of the specified file.
func backupFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
