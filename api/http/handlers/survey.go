package handlers

import (
	"Questify/api/http/handlers/presenter"
	qt "Questify/internal/question"
	"Questify/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		surveyId := c.Params("surveyId")
		if surveyId == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Survey ID is required")
		}

		var question presenter.Question
		if err := c.BodyParser(&question); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		surveyUUID, err := uuid.Parse(surveyId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Survey ID format")
		}

		file, err := c.FormFile("media")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("No file uploaded or invalid file")
		}

		filePath := "./uploads/" + file.Filename
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
		}

		domainQuestion := presenter.MapPresenterToQuestion(&question, filePath, surveyUUID)

		err = questionService.CreateQuestion(c.Context(), domainQuestion)
		if err != nil {
			switch err {
			case qt.ErrSurveyIdIsRequired:
				return c.Status(fiber.StatusBadRequest).SendString("Survey id is required!")
			case qt.ErrSurveyNotFound:
				return c.Status(fiber.StatusNotFound).SendString("Survey not found")
			case qt.ErrQuestionMultipleChoiceOptionsIsEmpty:
				return c.Status(fiber.StatusBadRequest).SendString("Multiple Choice question should have a list of options")
			case qt.ErrQuestionDescriptionShouldNotHaveMultipleChoiceList:
				return c.Status(fiber.StatusBadRequest).SendString("Description question should not contain a list of options")
			case qt.ErrQuestionMultipleChoiceItemsCountGreaterThanOne:
				return c.Status(fiber.StatusBadRequest).SendString("Question choices should be greater than 1")
			case qt.ErrDuplicateValueForQuestionChoicesNotAllowed:
				return c.Status(fiber.StatusBadRequest).SendString("Duplicate choice values are not allowed")
			default:
				return presenter.InternalServerError(c, err)
			}
		}

		return presenter.Created(c, "Created With ID: ", domainQuestion.ID)
	}
}
