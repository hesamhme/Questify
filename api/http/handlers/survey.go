package handlers

import (
	"Questify/api/http/handlers/presenter"
	qt "Questify/internal/question"
	"Questify/internal/survey"
	"Questify/pkg/jwt"
	"Questify/service"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		surveyId := c.Params("surveyId")
		if surveyId == "" {
			return presenter.BadRequest(c, errors.New("Survey ID is Required"))
		}

		surveyUUID, err := uuid.Parse(surveyId)
		if err != nil {
			return presenter.BadRequest(c, errors.New("Invalid Survey ID Format"))
		}

		var question presenter.Question
		if err := c.BodyParser(&question); err != nil {
			return presenter.BadRequest(c, errors.New("Invalid request body"))
		}

		choicesJSON := c.FormValue("question_choices")
		if choicesJSON != "" {
			var choices []presenter.QuestionChoice
			if err := json.Unmarshal([]byte(choicesJSON), &choices); err != nil {
				return presenter.BadRequest(c, errors.New("Invalid question_choices format"))
			}
			question.QuestionChoices = choices
		}

		var filePath string
		file, err := c.FormFile("media")
		if err == nil {
			filePath = "./uploads/" + file.Filename
			if err := c.SaveFile(file, filePath); err != nil {
				return presenter.InternalServerError(c, errors.New("Failed to save file"))
			}
		}

		domainQuestion := presenter.MapPresenterToQuestion(&question, filePath, surveyUUID)

		err = questionService.CreateQuestion(c.Context(), domainQuestion)
		if err != nil {
			switch err {
			case qt.ErrSurveyIdIsRequired, qt.ErrQuestionMultipleChoiceOptionsIsEmpty,
				qt.ErrQuestionDescriptionShouldNotHaveMultipleChoiceList,
				qt.ErrQuestionMultipleChoiceItemsCountGreaterThanOne,
				qt.ErrDuplicateValueForQuestionChoicesNotAllowed:
				return presenter.BadRequest(c, err)
			case qt.ErrSurveyNotFound:
				return presenter.NotFound(c, err)
			default:
				return presenter.InternalServerError(c, err)
			}
		}

		return presenter.Created(c, "Created With ID: ", domainQuestion.ID)
	}
}

func CreateSurvey(surveyService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var surveyReq presenter.Survey
		if err := c.BodyParser(&surveyReq); err != nil {
			return presenter.BadRequest(c, errors.New("Invalid request body"))
		}

		claims := c.Locals(jwt.UserClaimKey)
		userClaims, ok := claims.(*jwt.UserClaims)
		if !ok || userClaims.UserID == uuid.Nil {
			return presenter.Unauthorized(c, errors.New("User not authenticated"))
		}

		domainSurvey := presenter.MapPresenterToSurvey(&surveyReq, userClaims.UserID)

		err := surveyService.CreateSurvey(c.Context(), domainSurvey)
		if err != nil {
			switch {
			case errors.Is(err, survey.ErrInvalidTitle),
				errors.Is(err, survey.ErrInvalidTimeRange),
				errors.Is(err, survey.ErrPastStartTime),
				errors.Is(err, survey.ErrInvalidParticipationLimit),
				errors.Is(err, survey.ErrInvalidResponseTimeLimit):
				return presenter.BadRequest(c, err)
			default:
				return presenter.InternalServerError(c, err)
			}
		}

		return presenter.Created(c, "Created Survey with ID: ", domainSurvey.ID)
	}
}

func CreateAnswer(surveyService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		questionID := c.Params("questionId")
		if questionID == "" {
			return presenter.BadRequest(c, errors.New("Question ID is required"))
		}

		questionUUID, err := uuid.Parse(questionID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("Invalid Question ID Format"))
		}

		claims := c.Locals(jwt.UserClaimKey)
		userClaims, ok := claims.(*jwt.UserClaims)
		if !ok || userClaims.UserID == uuid.Nil {
			return presenter.Unauthorized(c, errors.New("User not authenticated"))
		}

		var req presenter.Answer
		if err := c.BodyParser(&req); err != nil {
			return presenter.BadRequest(c, errors.New("Invalid request payload"))
		}

		answer := presenter.MapPresenterToAnswer(&req, questionUUID, userClaims.UserID)

		err = surveyService.CreateAnswer(c.Context(), answer)
		if err != nil {
			switch {
			case errors.Is(err, qt.ErrUserIDRequired):
				return presenter.BadRequest(c, errors.New("User ID is required"))
			case errors.Is(err, qt.ErrQuestionIDRequired):
				return presenter.BadRequest(c, errors.New("Question ID is required"))
			case errors.Is(err, qt.ErrQuestionNotFound):
				return presenter.NotFound(c, errors.New("Question not found"))
			case errors.Is(err, qt.ErrInvalidAnswerForQuestionType):
				return presenter.BadRequest(c, errors.New("Invalid answer for the question type"))
			case errors.Is(err, qt.ErrUserAlreadyAnswered):
				return presenter.BadRequest(c, errors.New("User has already answered this question"))
			default:
				return presenter.InternalServerError(c, err)
			}
		}

		return presenter.Created(c, "Answer submitted successfully", nil)
	}
}

func GetNextQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		surveyID := c.Params("surveyId")
		if surveyID == "" {
			return presenter.BadRequest(c, errors.New("Survey ID is required"))
		}

		surveyUUID, err := uuid.Parse(surveyID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("Invalid Survey ID format"))
		}

		claims := c.Locals(jwt.UserClaimKey)
		userClaims, ok := claims.(*jwt.UserClaims)
		if !ok || userClaims.UserID == uuid.Nil {
			return presenter.Unauthorized(c, errors.New("User not authenticated"))
		}

		question, err := questionService.GetNextQuestion(c.Context(), surveyUUID, userClaims.UserID.String())
		if err != nil {
			switch err {
			case qt.ErrSurveyNotFound, qt.ErrNoMoreQuestionsForThisSurvey:
				return presenter.NotFound(c, err)
			default:
				return presenter.InternalServerError(c, err)
			}
		}

		if question == nil {
			return presenter.NotFound(c, errors.New("No more questions available"))
		}

		presentedQuestion := presenter.MapQuestionToPresenter(question)
		return presenter.OK(c, "Next question retrieved successfully", presentedQuestion)
	}
}

func GetPreviousQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		surveyID := c.Params("surveyId")
		if surveyID == "" {
			return presenter.BadRequest(c, errors.New("Survey ID is required"))
		}

		surveyUUID, err := uuid.Parse(surveyID)
		if err != nil {
			return presenter.BadRequest(c, errors.New("Invalid Survey ID format"))
		}

		claims := c.Locals(jwt.UserClaimKey)
		userClaims, ok := claims.(*jwt.UserClaims)
		if !ok || userClaims.UserID == uuid.Nil {
			return presenter.Unauthorized(c, errors.New("User not authenticated"))
		}

		question, err := questionService.GetPreviousQuestion(c.Context(), surveyUUID, userClaims.UserID.String())
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

func GetSurvey(surveyService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		surveyID := c.Params("surveyId")
		if surveyID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Survey ID is required")
		}

		surveyUUID, err := uuid.Parse(surveyID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Survey ID format")
		}

		survey, err := surveyService.GetSurvey(c.Context(), surveyUUID)
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

		presenterSurvey := presenter.MapSurveyToPresenter(survey)
		return c.Status(fiber.StatusOK).JSON(presenterSurvey)
	}
}

func UpdateQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		surveyId := c.Params("surveyId")
		if surveyId == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Survey ID is required")
		}

		surveyUUID, err := uuid.Parse(surveyId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Survey ID format")
		}

		questionID := c.Params("questionId")
		if questionID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Question ID is required")
		}

		questionUUID, err := uuid.Parse(questionID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Question ID format")
		}

		var question presenter.Question
		if err := c.BodyParser(&question); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		// Fetch the existing question to get the current media path
		existingQuestion, err := questionService.GetQuestion(c.Context(), questionUUID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Question not found")
		}

		var filePath string
		file, err := c.FormFile("media")
		if err == nil {
			filePath = "./uploads/" + file.Filename
			if err := c.SaveFile(file, filePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
			}
		} else if errors.Is(errors.New("there is no uploaded file associated with the given key"), err) {
			filePath = existingQuestion.MediaPath
		} else {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid file upload")
		}

		domainQuestion := presenter.MapPresenterToQuestion(&question, filePath, surveyUUID)

		domainQuestion.ID = questionUUID

		err = questionService.UpdateQuestion(c.Context(), domainQuestion, questionUUID)
		if err != nil {
			switch err {
			case qt.ErrQuestionNotFound:
				return c.Status(fiber.StatusNotFound).SendString("Question not found")
			case qt.ErrCannotChangeSurveyId:
				return c.Status(fiber.StatusBadRequest).SendString("Cannot change survey ID")
			case qt.ErrQuestionDescriptionShouldNotHaveMultipleChoiceList:
				return c.Status(fiber.StatusBadRequest).SendString("Description question should not contain a list of options")
			case qt.ErrQuestionMultipleChoiceOptionsIsEmpty:
				return c.Status(fiber.StatusBadRequest).SendString("Multiple Choice question should have a list of options")
			case qt.ErrQuestionMultipleChoiceItemsCountGreaterThanOne:
				return c.Status(fiber.StatusBadRequest).SendString("Question choices should be greater than 1")
			case qt.ErrDuplicateValueForQuestionChoicesNotAllowed:
				return c.Status(fiber.StatusBadRequest).SendString("Duplicate choice values are not allowed")
			default:
				return presenter.InternalServerError(c, err)
			}
		}

		return presenter.NoContent(c)
	}
}

func GetQuestion(questionService *service.SurveyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		questionID := c.Params("questionId")
		if questionID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Question ID is required")
		}

		questionUUID, err := uuid.Parse(questionID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid Question ID format")
		}

		claims := c.Locals(jwt.UserClaimKey)
		userClaims, ok := claims.(*jwt.UserClaims)
		if !ok || userClaims.UserID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).SendString("User not authenticated")
		}

		question, err := questionService.GetQuestion(c.Context(), questionUUID)
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
