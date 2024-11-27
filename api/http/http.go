package http

import (
	"Questify/api/http/handlers"
	"Questify/api/http/middlerwares"

	"github.com/gofiber/fiber/v2"
)

// SetupHTTP initializes the HTTP server and routes
func SetupHTTP(app *fiber.App, secretKey string) {
	api := app.Group("/api") // Base API group

	// User routes
	api.Post("/signup", handlers.Signup)
	api.Post("/signin", handlers.Signin)

	// Protected route
	api.Get("/protected", middlerwares.AuthMiddleware(secretKey), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		return c.JSON(fiber.Map{
			"message": "You have access to this protected route.",
			"user_id": userID,
		})
	})
}
