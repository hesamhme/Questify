package main

import (
	"log"

	"Questify/api/http"
	"Questify/config"
	"Questify/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Read configuration
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Initialize app container
	appContainer, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatalf("Error initializing app container: %v", err)
	}
	defer func() {
		sqlDB, _ := appContainer.DB.DB()
		sqlDB.Close()
	}()

	// Initialize Fiber app
	app := fiber.New()

	// Setup HTTP routes with dependencies from the app container
	http.SetupHTTP(app, cfg.JWT.Secret, appContainer)

	// Start the Fiber server
	log.Println("Starting server on port 3000...")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
