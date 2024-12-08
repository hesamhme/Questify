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

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newQuestion).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return nil
			}
			return err
		}

		if newQuestionChoices != nil {
			for _, choice := range *newQuestionChoices {
				choice.QuestionID = newQuestion.ID
				if err := tx.Create(&choice).Error; err != nil {
					if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
						continue
					}
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
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

func (r *questionRepo) Update(ctx context.Context, question *question.Question) error {
	// Convert domain model to entity
	updatedQuestion, updatedChoices := mappers.QuestionDomainToEntity(question)

	// Start a transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update the question
		if err := tx.Model(&entities.Question{}).Where("id = ?", updatedQuestion.ID).Updates(updatedQuestion).Error; err != nil {
			return fmt.Errorf("failed to update question: %w", err)
		}

		// Delete existing choices
		if err := tx.Where("question_id = ?", updatedQuestion.ID).Delete(&entities.QuestionChoices{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing choices: %w", err)
		}

		// Insert new choices
		if len(*updatedChoices) > 0 {
			if err := tx.Create(&updatedChoices).Error; err != nil {
				return fmt.Errorf("failed to create new choices: %w", err)
			}
		}

		return nil
	})
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

func (r *questionRepo) GetMaxQuestionIndexBySurveyID(ctx context.Context, surveyId uuid.UUID) (uint, error) {
	var maxIndex uint
	query := `SELECT COALESCE(MAX(index), 0) FROM questions WHERE survey_id = $1`
	err := r.db.WithContext(ctx).Raw(query, surveyId).Scan(&maxIndex)
	if err != nil {
		return 0, err.Error
	}
	return maxIndex, nil
}

func (r *questionRepo) GetBySurveyID(ctx context.Context, surveyID uuid.UUID) ([]*question.Question, error) {
	var questionEntities []entities.Question

	err := r.db.WithContext(ctx).
		Where("survey_id = ?", surveyID).
		Order("index ASC").
		Find(&questionEntities).Error
	if err != nil {
		return nil, err
	}

	questions := make([]*question.Question, 0, len(questionEntities))

	for _, questionEntity := range questionEntities {
		var questionChoices []entities.QuestionChoices
		err = r.db.WithContext(ctx).
			Where("question_id = ?", questionEntity.ID).
			Find(&questionChoices).Error
		if err != nil {
			return nil, err
		}

		questionDomain := mappers.QuestionEntityToDomain(questionEntity, questionChoices)
		questions = append(questions, &questionDomain)
	}

	return questions, nil
}

func (r *questionRepo) GetAnswerByUserAndQuestion(ctx context.Context, userID, questionID uuid.UUID) (*question.Answer, error) {
	var answerEntity entities.Answer
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		First(&answerEntity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, question.ErrAnswerNotFound
		}
		return nil, fmt.Errorf("failed to get answer: %w", err)
	}

	answer := mappers.AnswerEntityToDomain(answerEntity)
	return &answer, nil
}
