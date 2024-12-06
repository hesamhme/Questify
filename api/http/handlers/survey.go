package handlers

import (
	"Questify/api/http/handlers/presenter"
	"Questify/internal/survey"
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

func CreateSurvey(surveyService *survey.Ops) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var surveyReq presenter.Survey
		if err := c.BodyParser(&surveyReq); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		domainSurvey := presenter.MapPresenterToSurvey(&surveyReq)

		err := surveyService.Create(c.Context(), domainSurvey)
		if err != nil {
			return presenter.InternalServerError(c, err)
		}

		return presenter.Created(c, "Created Survey with ID: ", domainSurvey.ID)
	}
}

func CreateAnswer(surveyService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the request body
		var req presenter.Answer
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
		}

		// Map presenter answer to domain model
		answer := presenter.MapPresenterToAnswer(&req)

		// Call the service to create the answer
		err := surveyService.CreateAnswer(c.Context(), answer)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "Answer submitted successfully.",
		})
	}
}

func GetQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		surveyID := c.Params("surveyId")
		if surveyID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Survey ID is required")
		}

		surveyUUID, err := uuid.Parse(surveyID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Survey ID format")
		}

		userID := c.Locals("userID").(string)
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("User not authenticated")
		}

		question, err := questionService.GetNextQuestion(c.Context(), surveyUUID, userID)
		if err != nil {
			switch err {
			case qt.ErrSurveyNotFound:
				return c.Status(fiber.StatusNotFound).SendString("Survey not found")
			case qt.ErrNoMoreQuestionsForThisSurvey:
				return c.Status(fiber.StatusNotFound).SendString("No more questions available")
			default:
				return presenter.InternalServerError(c, err)
			}
		}

		if question == nil {
			return c.Status(fiber.StatusNotFound).SendString("No more questions available")
		}

		presentedQuestion := presenter.MapQuestionToPresenter(question)
		return c.Status(fiber.StatusOK).JSON(presentedQuestion)
	}
}
