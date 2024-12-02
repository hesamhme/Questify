package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/service"
	"github.com/gofiber/fiber/v2"
)

func CreateQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var question presenter.Question
		if err := c.BodyParser(&question); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		file, err := c.FormFile("media")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("No file uploaded or invalid file")
		}

		filePath := "./uploads/" + file.Filename
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
		}

		domainQuestion := presenter.MapPresenterToQuestion(&question, filePath)

		err = questionService.CreateQuestion(c.Context(), domainQuestion)
		if err != nil {
			// TODO validate domain errors and show correct errors
			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "Created With ID: ", domainQuestion.ID)
	}
}
