package service

import (
	"Questify/internal/question"
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type SurveyService struct {
	questionOps  *question.Ops
	userProgress map[string]uint
	mu           sync.RWMutex
}

func NewSurveyService(questionOps *question.Ops) *SurveyService {
	return &SurveyService{
		questionOps:  questionOps,
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
