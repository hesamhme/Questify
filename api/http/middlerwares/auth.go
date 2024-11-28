package middlerwares

import (
	"Questify/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates the JWT token and sets user context.
func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Missing Authorization header",
			})
		}

		// Parse and validate JWT
		claims, err := jwt.ParseToken(token, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token",
			})
		}

		// Set user ID in context for downstream handlers
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}
