package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Swagger struct {
		SummaryTemplate     string   `yaml:"summary_template"`
		DescriptionTemplate string   `yaml:"description_template"`
		Tags                []string `yaml:"tags"`
		Accept              string   `yaml:"accept"`
		Produce             string   `yaml:"produce"`
		SuccessResponse     string   `yaml:"success_response"`
		RouterTemplate      string   `yaml:"router_template"`
	} `yaml:"swagger"`
}

func LoadConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
