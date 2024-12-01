package question

import (
	"context"
	"github.com/google/uuid"
)

type Repo interface {
	Create(ctx context.Context, question *Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*Question, error)
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
	QuestionChoices []QuestionChoice
}

type QuestionChoice struct {
	ID       uint
	Value    string
	IsAnswer bool
}
