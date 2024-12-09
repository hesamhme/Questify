package survey

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)


type Repo interface {
	Create(ctx context.Context, survey *Survey) error
	GetByID(ctx context.Context, id uuid.UUID) (*Survey, error)
	GetUserSurveys(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Survey, int64, error)
}

var (
	ErrSurveyNotFound = errors.New("survey not found")
	ErrInvalidTitle              = errors.New("invalid title: must not be empty and should be less than 255 characters")
	ErrInvalidOwnerID            = errors.New("invalid owner ID")
	ErrInvalidTimeRange          = errors.New("end time must be after start time")
	ErrInvalidParticipationLimit = errors.New("participation limit must be greater than 0")
	ErrInvalidResponseTimeLimit  = errors.New("response time limit must be greater than 0")
	ErrPastStartTime             = errors.New("start time cannot be in the past")
)

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
