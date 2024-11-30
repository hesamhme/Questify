package main

import (
	"flag"
	"log"
	"os"

	http_server"Questify/api/http"
	"Questify/config"
	"Questify/service"
)

var configPath = flag.String("config", "", "configuration path")

func main() {
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http_server.Run(cfg, app)
}

// func readConfig() config.Config {
// 	flag.Parse()

// 	if cfgPathEnv := os.Getenv("APP_CONFIG_PATH"); len(cfgPathEnv) > 0 {
// 		*configPath = cfgPathEnv
// 	}

// 	if len(*configPath) == 0 {
// 		log.Fatal("configuration file not found")
// 	}

// 	cfg, err := config.ReadStandard(*configPath)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return cfg
// }

func readConfig() config.Config {
	flag.Parse()

	if cfgPathEnv := os.Getenv("APP_CONFIG_PATH"); len(cfgPathEnv) > 0 {
		*configPath = cfgPathEnv
	}

	// Set a default path if no config path is provided
	if len(*configPath) == 0 {
		*configPath = "config.yaml"
	}

	cfg, err := config.ReadStandard(*configPath)

	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	return cfg
}