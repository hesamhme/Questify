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

	// Return the signup response with a sign-in message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "You have successfully signed up and are now signed in.",
		"token":         response.Token,
		"email":         response.Email,
		"national_code": response.NationalCode,
	})
}

// Signin handles user login
func Signin(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// Call the service layer for authentication
	token, err := service.UserService.Login(input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return the generated JWT token
	return c.JSON(fiber.Map{
		"message": "You are successfully signed in.",
		"token":   token,
	})
}
