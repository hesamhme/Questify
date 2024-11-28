package main

import (
	"log"

	"Questify/api/http"
	"Questify/config"
	"Questify/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load Configuration
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize AppContainer
	appContainer, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatalf("Error initializing services: %v", err)
	}

	// Initialize Fiber
	app := fiber.New()

	// Setup Routes
	http.SetupHTTP(app, cfg.JWT.Secret, appContainer.AuthService)

	// Start Server
	log.Println("Starting server on port 3000...")
	log.Fatal(app.Listen(":3000"))
}
