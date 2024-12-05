package question

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, question *Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*Question, error)
	GetBySurveyID(ctx context.Context, surveyID uuid.UUID) ([]*Question, error)
	GetMaxQuestionIndexBySurveyID(ctx context.Context, surveyId uuid.UUID) (uint, error)
}

var (
	ErrSurveyNotFound                                     = errors.New("survey not found")
	ErrSurveyIdIsRequired                                 = errors.New("Survey Id Is required")
	ErrQuestionMultipleChoiceOptionsIsEmpty               = errors.New("Multiple Choice question should has list of question options")
	ErrQuestionDescriptionShouldNotHaveMultipleChoiceList = errors.New("Descriptiob question should not contain list of question options")
	ErrQuestionMultipleChoiceItemsCountGreaterThanOne     = errors.New("Question Choices should be greater that 1")
	ErrDuplicateValueForQuestionChoicesNotAllowed         = errors.New("duplicate choice values are not allowed")
	ErrNoMoreQuestionsForThisSurvey                       = errors.New("No more questions available")
)

type QuestionType int

const (
	DESCRIPTION QuestionType = iota
	MULTIPLE_CHOICE
)

type Question struct {
	ID              uuid.UUID
	Index           uint
	SurveyId        uuid.UUID
	Text            string
	Type            QuestionType
	IsMandatory     bool
	MediaPath       string
	QuestionChoices *[]QuestionChoice
}

type QuestionChoice struct {
	ID       uint
	Value    string
	IsAnswer bool
}
