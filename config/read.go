package config

import (
  "os"

  "gopkg.in/yaml.v3"
)

// ReadConfig reads and parses the configuration file.
func ReadConfig() (*Config, error) {
  file, err := os.Open("config.yaml")
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var cfg Config
  decoder := yaml.NewDecoder(file)
  if err := decoder.Decode(&cfg); err != nil {
    return nil, err
  }

  return &cfg, nil
}
