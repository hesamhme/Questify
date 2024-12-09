package survey

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo: repo}
}

func (o *Ops) Create(ctx context.Context, survey *Survey) error {
	if err := validateSurvey(survey); err != nil {
		return err
	}

	return o.repo.Create(ctx, survey)
}

func validateSurvey(s *Survey) error {
	if s.Title == "" || len(s.Title) > 255 {
		return ErrInvalidTitle
	}

	if s.OwnerID == uuid.Nil {
		return ErrInvalidOwnerID
	}

	now := time.Now()
	if s.StartTime.Before(now) {
		return ErrPastStartTime
	}
	if !s.EndTime.After(s.StartTime) {
		return ErrInvalidTimeRange
	}

	if s.ParticipationLimit == 0 {
		return ErrInvalidParticipationLimit
	}

	if s.ResponseTimeLimit <= 0 {
		return ErrInvalidResponseTimeLimit
	}

	return nil
}

func (o *Ops) GetByID(ctx context.Context, id uuid.UUID) (*Survey, error) {

	s, err := o.repo.GetByID(ctx, id)

	if err!=nil{
		return nil, err
	}
	if s == nil{
		return nil, ErrSurveyNotFound
	}
	return s, nil
}
