package handlers

import (
	"Questify/service"

	"github.com/gofiber/fiber/v2"
)

// Signup handles user registration
func Signup(c *fiber.Ctx) error {
	var input service.UserInput

	// Parse JSON input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// Call the service layer for business logic
	response, err := service.UserService.Signup(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}