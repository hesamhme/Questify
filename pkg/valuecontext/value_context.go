package valuecontext

import "github.com/gofiber/fiber/v2"

const (
	UserIDKey = "userID"
	RoleKey   = "role"
)

// SetUserID stores the user ID in the request context.
func SetUserID(c *fiber.Ctx, userID string) {
	c.Locals(UserIDKey, userID)
}

// GetUserID retrieves the user ID from the request context.
func GetUserID(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals(UserIDKey).(string)
	return userID, ok
}

// SetRole stores the user role in the request context.
func SetRole(c *fiber.Ctx, role string) {
	c.Locals(RoleKey, role)
}

// GetRole retrieves the user role from the request context.
func GetRole(c *fiber.Ctx) (string, bool) {
	role, ok := c.Locals(RoleKey).(string)
	return role, ok
}
