package storage

import (
	"Questify/internal/question"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/adapters/storage/mappers"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type questionRepo struct {
	db *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) question.Repo {
	return &questionRepo{
		db: db,
	}
}

func (r *questionRepo) Create(ctx context.Context, question *question.Question) error {
	newQuestion, newQuestionChoices := mappers.QuestionDomainToEntity(question)
	err := r.db.Create(&newQuestion).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}

	err = r.db.Create(&newQuestionChoices).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}

	question.ID = newQuestion.ID

	return nil
}

func (r *questionRepo) GetByID(ctx context.Context, id uuid.UUID) (*question.Question, error) {
	var question entities.Question

	err := r.db.WithContext(ctx).Model(&entities.Question{}).Where("id = ?", id).First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var questionChoices []entities.QuestionChoices
	err = r.db.WithContext(ctx).Model(&entities.QuestionChoices{}).Where("question_id = ?", id).Find(&questionChoices).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	questionDomain := mappers.QuestionEntityToDomain(question, questionChoices)
	return &questionDomain, nil
}
