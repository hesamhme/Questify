package service

import (
	"Questify/internal/question"
	"Questify/internal/survey"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type SurveyService struct {
	questionOps *question.Ops
	surveyOps   *survey.Ops
}

func NewSurveyService(questionOps *question.Ops, surveyOps *survey.Ops) *SurveyService {
	return &SurveyService{
		questionOps: questionOps,
		surveyOps:   surveyOps,
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

func (s *SurveyService) CreateSurvey(ctx context.Context, survey *survey.Survey) error {
	// TODO: Validate survey fields (e.g., time ranges, title length)
	return s.surveyOps.Create(ctx, survey)
}

func (s *SurveyService) GetSurvey(ctx context.Context, id string) (*survey.Survey, error) {
	surveyID, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return nil, err
	}

	fetchedSurvey, err := s.surveyOps.GetByID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	return fetchedSurvey, nil
}
