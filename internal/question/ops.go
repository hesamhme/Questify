package question

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo: repo}
}

func (o *Ops) Create(ctx context.Context, question *Question) error {

	//TODO: Validate question

	err := o.repo.Create(ctx, question)
	if err != nil {
		return err
	}
	return nil
}

func (o *Ops) GetByID(ctx context.Context, id uuid.UUID) (*Question, error) {
	return o.repo.GetByID(ctx, id)
}

// CreateAnswer validates and adds a new answer to the database
func (o *Ops) CreateAnswer(ctx context.Context, answer *Answer) error {
	// TODO: Add validation logic if needed

	err := o.repo.CreateAnswer(ctx, answer)
	if err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}
	return nil
}

// GetAnswersByQuestion retrieves all answers for a specific question
func (o *Ops) GetAnswersByQuestion(ctx context.Context, questionID uuid.UUID) ([]Answer, error) {
	answers, err := o.repo.GetAnswersByQuestion(ctx, questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get answers by question: %w", err)
	}
	return answers, nil
}

// GetAnswersByUser retrieves all answers submitted by a specific user
func (o *Ops) GetAnswersByUser(ctx context.Context, userID uuid.UUID) ([]Answer, error) {
	answers, err := o.repo.GetAnswersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get answers by user: %w", err)
	}
	return answers, nil
}
