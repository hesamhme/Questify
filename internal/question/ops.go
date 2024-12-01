package question

import (
	"context"
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
