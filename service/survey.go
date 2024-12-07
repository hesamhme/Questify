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
	surveyOps    *survey.Ops
	mu           sync.Mutex
	userProgress map[string]uint
}

func generateProgressKey(userID string, surveyID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", userID, surveyID.String())
}

func NewSurveyService(questionOps *question.Ops, surveyOps *survey.Ops) *SurveyService {
	return &SurveyService{
		questionOps:  questionOps,
		surveyOps:    surveyOps,
		userProgress: make(map[string]uint),
	}
}

func (s *SurveyService) CreateQuestion(ctx context.Context, question *question.Question) error {
	return s.questionOps.Create(ctx, question)
}

func (s *SurveyService) GetNextQuestion(ctx context.Context, surveyID uuid.UUID, userID string) (*question.Question, error) {
    s.mu.Lock()
    progressKey := generateProgressKey(userID, surveyID)
    currentIndex := s.userProgress[progressKey]
    s.userProgress[progressKey]++
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

func (s *SurveyService) GetPreviousQuestion(ctx context.Context, surveyID uuid.UUID, userID string) (*question.Question, error) {
    resultSurvey, err := s.surveyOps.GetByID(ctx, surveyID)
    if err != nil {
        return nil, err
    }

    if !resultSurvey.AllowBack {
        return nil, fmt.Errorf("going back is not allowed for this survey")
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    progressKey := generateProgressKey(userID, surveyID)
    currentIndex, exists := s.userProgress[progressKey]
    if !exists || currentIndex == 0 {
        return nil, fmt.Errorf("no previous questions available")
    }

    questions, err := s.questionOps.GetQuestionsBySurveyID(ctx, surveyID)
    if err != nil {
        return nil, err
    }

    if int(currentIndex) > len(questions) {
        return nil, fmt.Errorf("you have reached the end of the questions")
    }

    s.userProgress[progressKey]--

    currentIndex = s.userProgress[progressKey]
    return questions[currentIndex], nil
}

func (s *SurveyService) GetQuestion(ctx context.Context, questionID uuid.UUID) (*question.Question, error) {
	fetchedQuestion, err := s.questionOps.GetByID(ctx, questionID)
	if err != nil {
		return nil, err
	}

	return fetchedQuestion, nil
}

func (s *SurveyService) UpdateQuestion(ctx context.Context, question *question.Question, questionId uuid.UUID) error {
	return s.questionOps.Update(ctx, question, questionId)
}

func (s *SurveyService) CreateSurvey(ctx context.Context, survey *survey.Survey) error {
	return s.surveyOps.Create(ctx, survey)
}

func (s *SurveyService) GetSurvey(ctx context.Context, surveyID uuid.UUID) (*survey.Survey, error) {
	fetchedSurvey, err := s.surveyOps.GetByID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	return fetchedSurvey, nil
}
