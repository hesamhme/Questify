package handlers

import (
	"Questify/service"
	"github.com/gofiber/fiber/v2"
)

// SignupHandler handles user signup requests.
func SignupHandler(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			NID      string `json:"nid"`
		}

		// Parse and validate the request body
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Call the signup logic
		token, err := authService.SignUp(req.Email, req.Password, req.NID)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return the token
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"token": token,
		})
	}
}


func LoginHandler(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		token, err := authService.SignIn(req.Email, req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
	}
}


// UserProfileHandler fetches the profile of the authenticated user.
func UserProfileHandler(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the user ID from the context
		userID, exists := c.Locals("userID").(string)
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		// Fetch the user details
		user, err := userService.GetUserByID(userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.JSON(user)
	}
}