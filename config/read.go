package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ReadConfig reads the configuration from config.yaml
func ReadConfig() (*Config, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config.yaml: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config.yaml: %w", err)
	}

	// Debug: Print the loaded config
	fmt.Printf("Loaded config: %+v\n", cfg)

	return &cfg, nil
}
