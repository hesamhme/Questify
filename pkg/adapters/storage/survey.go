package storage

import (
	"Questify/internal/survey"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/adapters/storage/mappers"
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type surveyRepo struct {
	db *gorm.DB
}

func NewSurveyRepo(db *gorm.DB) survey.Repo {
	return &surveyRepo{db: db}
}

func (r *surveyRepo) Create(ctx context.Context, survey *survey.Survey) error {
	entitySurvey := mappers.SurveyDomainToEntity(survey)
	err := r.db.Create(&entitySurvey).Error
	if err != nil {
		return err
	}
	survey.ID = entitySurvey.ID
	return nil
}

func (r *surveyRepo) GetByID(ctx context.Context, id uuid.UUID) (*survey.Survey, error) {
	var entitySurvey entities.Survey
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entitySurvey).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, err
	}
	domainSurvey := mappers.SurveyEntityToDomain(entitySurvey)
	return &domainSurvey, nil
}
