package http

import (
	"Questify/api/http/handlers"
	"Questify/api/http/middlerwares"
	"Questify/service"
	"github.com/gofiber/fiber/v2"
)

// SetupHTTP sets up all routes and middleware for the app.
func SetupHTTP(app *fiber.App, jwtSecret string, container *service.AppContainer) {
	// Public routes
	app.Post("/signup", handlers.SignupHandler(container.AuthService))
	app.Post("/signin", handlers.LoginHandler(container.AuthService))

	// Protected routes
	protected := app.Group("/api", middlerwares.AuthMiddleware(jwtSecret))
	protected.Get("/profile", handlers.UserProfileHandler(container.UserService))
}
