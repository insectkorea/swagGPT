package main

import (
	"log"
	"os"

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
					configPath := c.String("config")

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

					err = handler.ProcessFiles(files, client, dryRun, configPath)
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
						Name:  "config",
						Usage: "Path to configuration file",
						Value: "configs/config.yaml",
					},
				},
			},
			{
				Name:  "undo",
				Usage: "Restore original files from backups",
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					if dir == "" {
						return cli.Exit("directory is required", 1)
					}

					err := handler.RestoreFiles(dir)
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
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
