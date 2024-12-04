package storage

import (
	"Questify/internal/question"
	"Questify/pkg/adapters/storage/entities"
	"Questify/pkg/adapters/storage/mappers"
	"context"
	"errors"
	"fmt"

	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type questionRepo struct {
	db *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) question.Repo {
	return &questionRepo{
		db: db,
	}
}

// Create creates a new question along with its choices
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

// GetByID retrieves a question by its ID along with its choices
func (r *questionRepo) GetByID(ctx context.Context, id uuid.UUID) (*question.Question, error) {
	var questionEntity entities.Question

	err := r.db.WithContext(ctx).Model(&entities.Question{}).Where("id = ?", id).First(&questionEntity).Error
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

	questionDomain := mappers.QuestionEntityToDomain(questionEntity, questionChoices)
	return &questionDomain, nil
}

// CreateAnswer adds a new answer to the database
func (r *questionRepo) CreateAnswer(ctx context.Context, answer *question.Answer) error {
	newAnswer := mappers.AnswerDomainToEntity(*answer)
	err := r.db.WithContext(ctx).Create(&newAnswer).Error
	if err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}
	answer.ID = newAnswer.ID
	return nil
}

func (r *questionRepo) GetAnswersByQuestion(ctx context.Context, questionID uuid.UUID, limit, offset int) ([]question.Answer, error) {
	var answerEntities []entities.Answer
	err := r.db.WithContext(ctx).
		Where("question_id = ?", questionID).
		Limit(limit).Offset(offset). // Add pagination
		Find(&answerEntities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get answers by question: %w", err)
	}
	return mappers.BatchAnswerEntityToDomain(answerEntities), nil
}

func (r *questionRepo) GetAnswersByUser(ctx context.Context, userID, surveyID uuid.UUID, limit, offset int) ([]question.Answer, error) {
	var answerEntities []entities.Answer
	err := r.db.WithContext(ctx).
		Joins("JOIN questions ON answers.question_id = questions.id").
		Where("answers.user_id = ? AND questions.survey_id = ?", userID, surveyID). // Filter by user and survey
		Limit(limit).Offset(offset).                                                // Add pagination
		Find(&answerEntities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get answers by user: %w", err)
	}
	return mappers.BatchAnswerEntityToDomain(answerEntities), nil
}
