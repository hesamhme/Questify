package http

import (
	"Questify/api/http/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupHTTP initializes the HTTP server and routes
func SetupHTTP(app *fiber.App) {
	api := app.Group("/api") // Base API group

	// User routes
	api.Post("/signup", handlers.Signup)
}
