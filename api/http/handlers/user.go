package handlers

import (
	"Questify/service"

	"github.com/gofiber/fiber/v2"
)

// Signup handles the user registration process
func Signup(c *fiber.Ctx) error {
	var input service.UserInput

	// Parse JSON input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// Call the UserService to handle business logic
	err := service.UserService.Signup(input)
	if err != nil {
		// Handle specific errors (e.g., duplicate or invalid NID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}
