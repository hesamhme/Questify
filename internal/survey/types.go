package survey

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Repo interface {
	Create(ctx context.Context, survey *Survey) error
	GetByID(ctx context.Context, id uuid.UUID) (*Survey, error)
}

type Survey struct {
	ID                 uuid.UUID
	Title              string
	OwnerID            uuid.UUID
	StartTime          time.Time
	EndTime            time.Time
	IsRandom           bool
	IsCanceled         bool
	AllowBack          bool
	ParticipationLimit uint
	ResponseTimeLimit  int64
}
