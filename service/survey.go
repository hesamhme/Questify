package service

import (
	"Questify/internal/question"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type SurveyService struct {
	questionOps *question.Ops
}

func NewSurveyService(questionOps *question.Ops) *SurveyService {
	return &SurveyService{
		questionOps: questionOps,
	}
}

func (s *SurveyService) CreateQuestion(ctx context.Context, question *question.Question) error {
	//TODO: Check if survey exist!
	return s.questionOps.Create(ctx, question)
}

func (s *SurveyService) GetQuestion(ctx context.Context, id string) (*question.Question, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	fetchedQuestion, err := s.questionOps.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return fetchedQuestion, nil
}
