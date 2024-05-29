package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/insectkorea/swagGPT/internal/api"
	"github.com/insectkorea/swagGPT/internal/handler"
	"github.com/insectkorea/swagGPT/internal/scanner"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

func main() {
	// Initialize logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	app := &cli.App{
		Name:  "Swagger Comment Adder",
		Usage: "Add Swagger comments to Gin handler functions",
		Commands: []*cli.Command{
			{
				Name:  "add-comments",
				Usage: "Add Swagger comments to handler functions",
				Action: func(c *cli.Context) error {
					dryRun := c.Bool("dry-run")
					model := c.String("model")
					skipPrompt := c.Bool("yes")

					dir := c.String("dir")
					if dir == "" {
						return cli.Exit("directory is required", 1)
					}

					apiKey := os.Getenv("OPENAI_API_KEY")
					if apiKey == "" {
						return cli.Exit("OpenAI API key is required", 1)
					}

					client := api.NewOpenAIClient(apiKey)

					files, err := scanner.ScanDir(dir)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}

					// Estimate total tokens and cost
					totalTokens := handler.EstimateTotalTokens(files)

					logrus.Infof(
						`
Estimated total tokens: %d
Estimated cost (approx): $%.2f
This approximation is based on gpt-4o model($5.00 / 1M tokens). 
Please check OpenAI's pricing for other models.`, totalTokens, float64(totalTokens)/1000000*5)
					// GPT-4o cost $5.00 / 1M tokens
					// Prompt user for confirmation unless --yes flag is provided
					if !skipPrompt {
						fmt.Print("Do you want to proceed? (y/N): ")
						reader := bufio.NewReader(os.Stdin)
						response, err := reader.ReadString('\n')
						if err != nil {
							return err
						}
						response = strings.TrimSpace(strings.ToLower(response))
						if response != "y" && response != "yes" {
							fmt.Println("Operation aborted.")
							return nil
						}
					}

					err = handler.ProcessFiles(files, client, dryRun, model)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dir",
						Usage:    "Directory to scan for Go files",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Preview changes without writing to files",
					},
					&cli.StringFlag{
						Name:  "model",
						Usage: "OpenAI model to use",
						Value: "gpt-4o",
					},
					&cli.BoolFlag{
						Name:  "yes",
						Usage: "Skip confirmation prompt",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
