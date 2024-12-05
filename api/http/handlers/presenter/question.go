package presenter

import (
	"Questify/internal/question"

	"github.com/google/uuid"
)

type QuestionType string

const (
	TextQuestion   QuestionType = "text"
	ChoiceQuestion QuestionType = "choice"
)

type QuestionChoice struct {
	ChoiceText string `json:"choice_text"`
}

type Question struct {
	ID              uuid.UUID        `json:"id"`
	Index           uint             `json:"index"`
	Text            string           `json:"text"`
	Type            QuestionType     `json:"type"`
	IsMandatory     bool             `json:"is_mandatory"`
	QuestionChoices []QuestionChoice `json:"question_choices,omitempty"`
	MediaURL        string           `json:"media_url,omitempty"`
}

func MapPresenterToQuestion(presenterQuestion *Question, mediaPath string, surveyId uuid.UUID) *question.Question {
	var qType question.QuestionType
	switch presenterQuestion.Type {
	case TextQuestion:
		qType = question.DESCRIPTION
	case ChoiceQuestion:
		qType = question.MULTIPLE_CHOICE
	default:
		qType = question.DESCRIPTION
	}

	var qChoices []question.QuestionChoice
	for _, qc := range presenterQuestion.QuestionChoices {
		qChoices = append(qChoices, question.QuestionChoice{
			Value: qc.ChoiceText,
		})
	}

	// Fix: Wrap qChoices as a pointer
	return &question.Question{
		ID:          presenterQuestion.ID,
		Index:       presenterQuestion.Index,
		SurveyId:    surveyId,
		Text:        presenterQuestion.Text,
		Type:        qType,
		IsMandatory: presenterQuestion.IsMandatory,

		QuestionChoices: &qChoices,
		MediaPath:       mediaPath,
	}
}

type Answer struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	UserID     uuid.UUID `json:"user_id"`
	Response   string    `json:"response"`
}

func MapPresenterToAnswer(presenterAnswer *Answer) *question.Answer {
	return &question.Answer{
		ID:         presenterAnswer.ID,
		QuestionID: presenterAnswer.QuestionID,
		UserID:     presenterAnswer.UserID,
		Response:   presenterAnswer.Response,
	}
}

func MapQuestionToPresenter(q *question.Question) Question {
	presentedQuestion := Question{
		ID:          q.ID,
		Index:       q.Index,
		Text:        q.Text,
		IsMandatory: q.IsMandatory,
	}

	switch q.Type {
	case question.DESCRIPTION:
		presentedQuestion.Type = TextQuestion
	case question.MULTIPLE_CHOICE:
		presentedQuestion.Type = ChoiceQuestion
		if q.QuestionChoices != nil {
			presentedQuestion.QuestionChoices = make([]QuestionChoice, len(*q.QuestionChoices))
			for i, choice := range *q.QuestionChoices {
				presentedQuestion.QuestionChoices[i] = QuestionChoice{
					ChoiceText: choice.Value,
				}
			}
		}
	}

	// Check if MediaPath exists and set it to MediaURL
	if q.MediaPath != "" {
		presentedQuestion.MediaURL = q.MediaPath
	}

	return presentedQuestion
}
