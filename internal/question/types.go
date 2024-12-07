package question

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, question *Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*Question, error)
	Update(ctx context.Context, question *Question) error
	CreateAnswer(ctx context.Context, answer *Answer) error
	GetAnswersByQuestion(ctx context.Context, questionID uuid.UUID, limit, offset int) ([]Answer, error)   // Updated
	GetAnswersByUser(ctx context.Context, userID, surveyID uuid.UUID, limit, offset int) ([]Answer, error) // Updated
	GetBySurveyID(ctx context.Context, surveyID uuid.UUID) ([]*Question, error)
	GetMaxQuestionIndexBySurveyID(ctx context.Context, surveyId uuid.UUID) (uint, error)
	GetAnswerByUserAndQuestion(ctx context.Context, userID, questionID uuid.UUID) (*Answer, error)
}

var (
	ErrSurveyNotFound                                     = errors.New("survey not found")
	ErrSurveyIdIsRequired                                 = errors.New("survey Id Is required")
	ErrQuestionMultipleChoiceOptionsIsEmpty               = errors.New("multiple Choice question should has list of question options")
	ErrQuestionDescriptionShouldNotHaveMultipleChoiceList = errors.New("descriptiob question should not contain list of question options")
	ErrQuestionMultipleChoiceItemsCountGreaterThanOne     = errors.New("Question Choices should be greater that 1")
	ErrDuplicateValueForQuestionChoicesNotAllowed         = errors.New("duplicate choice values are not allowed")
	ErrNoMoreQuestionsForThisSurvey                       = errors.New("no more questions available")
	ErrQuestionNotFound                                   = errors.New("Question not found")
	ErrCannotChangeSurveyId                               = errors.New("Can not change survey id")
	ErrAnswerNotFound                                     = errors.New("Answer not found")
	ErrUserIDRequired                                     = errors.New("user ID is required")
	ErrQuestionIDRequired                                 = errors.New("question ID is required")
	ErrInvalidAnswerForQuestionType                       = errors.New("invalid answer for the question type")
	ErrUserAlreadyAnswered                                = errors.New("user has already answered this question")
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

// Answer represents an answer to a question
type Answer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	UserID     uuid.UUID
	Response   string
	CreatedAt  time.Time
}
