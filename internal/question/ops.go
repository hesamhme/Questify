package question

import (
	"context"
	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo: repo}
}

func (o *Ops) Create(ctx context.Context, question *Question) error {

	if question.SurveyId == uuid.Nil {
		return ErrSurveyIdIsRequired
	}

	//TODO: Validate survey is exist!
	// if data, err := o.repo.GetBySurveyByID; err != nil {
	//	return ErrSurveyNotFound
	//}

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

	err := o.repo.Create(ctx, question)

	if err != nil {
		return err
	}
	return nil
}

func (o *Ops) GetByID(ctx context.Context, id uuid.UUID) (*Question, error) {
	return o.repo.GetByID(ctx, id)
}
