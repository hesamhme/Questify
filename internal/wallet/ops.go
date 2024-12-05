package wallet

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{
		repo: repo,
	}
}

func (o *Ops) Create(ctx context.Context, wallet *Wallet) (error, uuid.UUID) {
	err, id := o.repo.Create(ctx, wallet)
	if err != nil {
		log.Fatal(err)
		return err, id
	}
	return nil, id
}

func (o *Ops) Delete(ctx context.Context, id uuid.UUID) {
	o.repo.Delete(ctx, id)
}
