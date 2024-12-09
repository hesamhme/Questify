package mappers

import (
	"Questify/internal/survey"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/fp"
)

func SurveyEntityToDomain(entity entities.Survey) survey.Survey {
	return survey.Survey{
		ID:                 entity.ID,
		Title:              entity.Title,
		OwnerID:            entity.OwnerID,
		StartTime:          entity.StartTime,
		EndTime:            entity.EndTime,
		IsRandom:           entity.IsRandom,
		IsCanceled:         entity.IsCanceled,
		AllowBack:          entity.AllowBack,
		ParticipationLimit: entity.ParticipationLimit,
		ResponseTimeLimit:  entity.ResponseTimeLimit,
	}
}

func BatchSurveyEntityToDomain(entities []entities.Survey) []survey.Survey {
	return fp.Map(entities, SurveyEntityToDomain)
}

func SurveyDomainToEntity(domain *survey.Survey) *entities.Survey {
	return &entities.Survey{
		ID:                 domain.ID,
		Title:              domain.Title,
		OwnerID:            domain.OwnerID,
		StartTime:          domain.StartTime,
		EndTime:            domain.EndTime,
		IsRandom:           domain.IsRandom,
		IsCanceled:         domain.IsCanceled,
		AllowBack:          domain.AllowBack,
		ParticipationLimit: domain.ParticipationLimit,
		ResponseTimeLimit:  domain.ResponseTimeLimit,
	}
}
