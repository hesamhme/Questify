package main

import (
	"log"

	"Questify/config"
	"Questify/pkg/adapters/storage"
)

func main() {
	// Read configuration
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Setup database
	storage.SetupDatabase(&cfg.DB)
	// yadam bashe ino badan pak konm ...
	log.Println("Server is running... (database initialized)")

}
