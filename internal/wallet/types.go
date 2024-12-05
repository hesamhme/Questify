package wallet

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID
	Credit    float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repo interface {
	Create(ctx context.Context, wallet *Wallet) (error, uuid.UUID)
	Update(ctx context.Context, wallet *Wallet) error
	GetById(ctx context.Context, Id uuid.UUID) (*Wallet, error)
	Delete(ctx context.Context, Id uuid.UUID)
}
