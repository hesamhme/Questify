package main

import (
	"log"

	"Questify/api/http"
	"Questify/config"
	"Questify/pkg/adapters/storage"
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
	// yadam bashe ino badan pak konm ...
	log.Println("Server is running... (database initialized)")

	// Initialize Fiber
	app := fiber.New()

	// Setup HTTP routes
	http.SetupHTTP(app)

	// Start server
	log.Fatal(app.Listen(":3000"))

}
