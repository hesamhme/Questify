package question

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, question *Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*Question, error)

	// New methods for Answer
	CreateAnswer(ctx context.Context, answer *Answer) error
	GetAnswersByQuestion(ctx context.Context, questionID uuid.UUID) ([]Answer, error)
	GetAnswersByUser(ctx context.Context, userID uuid.UUID) ([]Answer, error)
}

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
	QuestionChoices *[]QuestionChoice // todo: make it pointer to handle nil ops
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
