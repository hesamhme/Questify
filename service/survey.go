package service

import (
	"Questify/internal/question"
	"Questify/internal/survey"
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type SurveyService struct {
	questionOps  *question.Ops
	userProgress map[string]uint
	mu           sync.RWMutex
	surveyOps   *survey.Ops
}

func NewSurveyService(questionOps *question.Ops, surveyOps *survey.Ops) *SurveyService {
	return &SurveyService{
		questionOps: questionOps,
		surveyOps:   surveyOps,
		userProgress: make(map[string]uint),
	}
}

func (s *SurveyService) CreateQuestion(ctx context.Context, question *question.Question) error {
	return s.questionOps.Create(ctx, question)
}

func (s *SurveyService) GetNextQuestion(ctx context.Context, surveyID uuid.UUID, userID string) (*question.Question, error) {
	s.mu.Lock()
	currentIndex := s.userProgress[userID]
	s.userProgress[userID]++
	s.mu.Unlock()

	questions, err := s.questionOps.GetQuestionsBySurveyID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	if int(currentIndex) >= len(questions) {
		return nil, fmt.Errorf("no more questions")
	}

	return questions[currentIndex], nil
}

func (s *SurveyService) ResetUserProgress(userID string) {
	s.mu.Lock()
	delete(s.userProgress, userID)
	s.mu.Unlock()
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