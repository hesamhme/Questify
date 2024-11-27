package main

import (
	"log"

	"Questify/api/http"
	"Questify/config"
	"Questify/pkg/adapters/storage"
	"Questify/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Read configuration
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Setup database
	storage.SetupDatabase(&cfg.DB)

	// Initialize services
	service.InitializeServices(cfg)

	// Initialize Fiber
	app := fiber.New()

	// Setup HTTP routes with the JWT secret key
	http.SetupHTTP(app, cfg.JWT.Secret)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
