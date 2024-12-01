package mappers

import (
	"Questify/internal/question"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/fp"
	"time"
)

func QuestionEntityToDomain(entity entities.Question, questionChoices []entities.QuestionChoices) question.Question {
	return question.Question{
		ID:              entity.ID,
		Index:           entity.Index,
		SurveyId:        entity.SurveyID,
		Text:            entity.Text,
		Type:            question.QuestionType(entity.Type),
		IsMandatory:     entity.IsMandatory,
		MediaPath:       entity.MediaPath,
		QuestionChoices: BatchQuestionChoiceEntityToDomain(questionChoices),
	}
}

func QuestionChoiceEntityToDomain(entity entities.QuestionChoices) question.QuestionChoice {
	return question.QuestionChoice{
		ID:       entity.ID,
		Value:    entity.Value,
		IsAnswer: entity.IsAnswer,
	}
}

func BatchQuestionChoiceEntityToDomain(entities []entities.QuestionChoices) []question.QuestionChoice {
	return fp.Map(entities, QuestionChoiceEntityToDomain)
}

func QuestionDomainToEntity(domainQuestion *question.Question) (*entities.Question, *[]entities.QuestionChoices) {
	questionEntity := entities.Question{
		ID:          domainQuestion.ID,
		Index:       domainQuestion.Index,
		SurveyID:    domainQuestion.SurveyId,
		Survey:      entities.Survey{},
		Text:        domainQuestion.Text,
		Type:        uint8(domainQuestion.Type),
		IsMandatory: domainQuestion.IsMandatory,
		MediaPath:   domainQuestion.MediaPath,
		CreatedAt:   time.Now(),
	}

	questionChoiceEntities := make([]entities.QuestionChoices, 0)
	for _, questionChoice := range domainQuestion.QuestionChoices {
		questionChoiceEntities = append(questionChoiceEntities, QuestionChoiceDomainToEntity(questionChoice, &questionEntity))
	}
	
	return &questionEntity, &questionChoiceEntities
}

func QuestionChoiceDomainToEntity(domain question.QuestionChoice, entityQuestion *entities.Question) entities.QuestionChoices {
	return entities.QuestionChoices{
		ID:         domain.ID,
		QuestionID: entityQuestion.ID,
		Question:   *entityQuestion,
		Value:      domain.Value,
		IsAnswer:   domain.IsAnswer,
	}
}
