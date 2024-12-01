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
	SurveyId        uuid.UUID        `json:"survey_id"`
	Text            string           `json:"text"`
	Type            QuestionType     `json:"type"`
	IsMandatory     bool             `json:"is_mandatory"`
	QuestionChoices []QuestionChoice `json:"question_choices,omitempty"`
}

func MapPresenterToQuestion(presenterQuestion *Question, mediaPath string) *question.Question {
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

	return &question.Question{
		ID:              presenterQuestion.ID,
		Index:           presenterQuestion.Index,
		SurveyId:        presenterQuestion.SurveyId,
		Text:            presenterQuestion.Text,
		Type:            qType,
		IsMandatory:     presenterQuestion.IsMandatory,
		QuestionChoices: qChoices,
		MediaPath:       mediaPath,
	}
}
