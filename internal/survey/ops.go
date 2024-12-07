package survey

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

func (o *Ops) Create(ctx context.Context, survey *Survey) error {
	// TODO: Validate survey
	return o.repo.Create(ctx, survey)
}

func (o *Ops) GetByID(ctx context.Context, id uuid.UUID) (*Survey, error) {
	return o.repo.GetByID(ctx, id)
}
