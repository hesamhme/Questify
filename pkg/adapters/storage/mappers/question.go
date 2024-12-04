package mappers

import (
	"Questify/internal/question"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/fp"
	"time"
)

func QuestionEntityToDomain(entity entities.Question, questionChoices []entities.QuestionChoices) question.Question {

	choices := BatchQuestionChoiceEntityToDomain(questionChoices)
	
	return question.Question{
		ID:              entity.ID,
		Index:           entity.Index,
		SurveyId:        entity.SurveyID,
		Text:            entity.Text,
		Type:            question.QuestionType(entity.Type),
		IsMandatory:     entity.IsMandatory,
		MediaPath:       entity.MediaPath,
		QuestionChoices: &choices,
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
	for _, questionChoice := range *domainQuestion.QuestionChoices {
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


func AnswerEntityToDomain(entity entities.Answer) question.Answer {
	return question.Answer{
		ID:         entity.ID,
		QuestionID: entity.QuestionID, // Directly use the QuestionID field
		UserID:     entity.UserID,     // Directly use the UserID field
		Response:   entity.Response,
		CreatedAt:  entity.CreatedAt,
	}
}

func AnswerDomainToEntity(domain question.Answer) entities.Answer {
	return entities.Answer{
		ID:         domain.ID,
		QuestionID: domain.QuestionID, // Directly assign QuestionID
		UserID:     domain.UserID,     // Directly assign UserID
		Response:   domain.Response,
		CreatedAt:  domain.CreatedAt,
	}
}


func BatchAnswerEntityToDomain(entities []entities.Answer) []question.Answer {
	var domainAnswers []question.Answer
	for _, entity := range entities {
		domainAnswers = append(domainAnswers, AnswerEntityToDomain(entity))
	}
	return domainAnswers
}