package question

import (
	domainSurvey "Questify/internal/survey"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Ops struct {
	repo       Repo
	surveyRepo domainSurvey.Repo
}

func NewOps(repo Repo, surveyRepo domainSurvey.Repo) *Ops {
	return &Ops{
		repo:       repo,
		surveyRepo: surveyRepo,
	}
}

func (o *Ops) Create(ctx context.Context, question *Question) error {
    if question.SurveyId == uuid.Nil {
        return ErrSurveyIdIsRequired
    }

    if data, err := o.surveyRepo.GetByID(ctx, question.SurveyId); err != nil || data == nil {
        return ErrSurveyNotFound
    }

    if question.Type == DESCRIPTION {
        if question.QuestionChoices != nil && len(*question.QuestionChoices) > 0 {
            return ErrQuestionDescriptionShouldNotHaveMultipleChoiceList
        }
    } else if question.Type == MULTIPLE_CHOICE {
        if question.QuestionChoices == nil || len(*question.QuestionChoices) <= 0 {
            return ErrQuestionMultipleChoiceOptionsIsEmpty
        }

        if len(*question.QuestionChoices) <= 1 {
            return ErrQuestionMultipleChoiceItemsCountGreaterThanOne
        }

        seenValues := make(map[string]bool)
        for _, choice := range *question.QuestionChoices {
            if seenValues[choice.Value] {
                return ErrDuplicateValueForQuestionChoicesNotAllowed
            }
            seenValues[choice.Value] = true
        }
    }

    maxIndex, err := o.repo.GetMaxQuestionIndexBySurveyID(ctx, question.SurveyId)
    if err != nil {
        return err
    }
    question.Index = maxIndex + 1

    err = o.repo.Create(ctx, question)
    if err != nil {
        return err
    }
    return nil
}

func (o *Ops) Update(ctx context.Context, question *Question, questionId uuid.UUID) error {
	if question.SurveyId == uuid.Nil {
		return ErrSurveyIdIsRequired
	}

	existingQuestion, err := o.repo.GetByID(ctx, questionId)
	if err != nil {
		return fmt.Errorf("failed to get existing question: %w", err)
	}
	if existingQuestion == nil {
		return ErrQuestionNotFound
	}

	if existingQuestion.SurveyId != question.SurveyId {
		return ErrCannotChangeSurveyId
	}

	if question.Type == DESCRIPTION && len(*question.QuestionChoices) > 0 {
		return ErrQuestionDescriptionShouldNotHaveMultipleChoiceList
	}

	if question.Type == MULTIPLE_CHOICE {
		if len(*question.QuestionChoices) <= 0 {
			return ErrQuestionMultipleChoiceOptionsIsEmpty
		}

		if len(*question.QuestionChoices) <= 1 {
			return ErrQuestionMultipleChoiceItemsCountGreaterThanOne
		}

		seenValues := make(map[string]bool)
		for _, choice := range *question.QuestionChoices {
			if seenValues[choice.Value] {
				return ErrDuplicateValueForQuestionChoicesNotAllowed
			}
			seenValues[choice.Value] = true
		}
	}

	question.ID = questionId
	question.Index = existingQuestion.Index

	err = o.repo.Update(ctx, question)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}

	return nil
}

func (o *Ops) GetByID(ctx context.Context, id uuid.UUID) (*Question, error) {
	return o.repo.GetByID(ctx, id)
}

// CreateAnswer validates and adds a new answer to the database
func (o *Ops) CreateAnswer(ctx context.Context, answer *Answer) error {
	// TODO: Add validation logic if needed

	err := o.repo.CreateAnswer(ctx, answer)
	if err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}
	return nil
}

// GetAnswersByQuestion retrieves all answers for a specific question with pagination
func (o *Ops) GetAnswersByQuestion(ctx context.Context, questionID uuid.UUID, limit, offset int) ([]Answer, error) {
	answers, err := o.repo.GetAnswersByQuestion(ctx, questionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get answers by question: %w", err)
	}
	return answers, nil
}

// GetAnswersByUser retrieves all answers submitted by a specific user for a survey with pagination
func (o *Ops) GetAnswersByUser(ctx context.Context, userID, surveyID uuid.UUID, limit, offset int) ([]Answer, error) {
	answers, err := o.repo.GetAnswersByUser(ctx, userID, surveyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get answers by user: %w", err)
	}
	return answers, nil
}

func (o *Ops) GetQuestionsBySurveyID(ctx context.Context, surveyID uuid.UUID) ([]*Question, error) {
	return o.repo.GetBySurveyID(ctx, surveyID)

}
