package http

import (
	"Questify/api/http/handlers"
	"Questify/service"
	"github.com/gofiber/fiber/v2"
)

// SetupHTTP sets up all routes and middleware for the app.
func SetupHTTP(app *fiber.App, jwtSecret string, authService *service.AuthService) {
	// Public routes
	app.Post("/signup", handlers.SignupHandler(authService))
	app.Post("/signin", handlers.LoginHandler(authService))

	// Protected routes
	// protected := app.Group("/api", middlerwares.AuthMiddleware(jwtSecret))
	// protected.Get("/profile", handlers.UserProfileHandler(authService)) // Example: Protected route
}
