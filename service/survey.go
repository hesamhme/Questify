package service

import (
	"Questify/internal/question"
	"Questify/internal/survey"
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SurveyService struct {
	questionOps    *question.Ops
	surveyOps      *survey.Ops
	mu             sync.Mutex
	userProgress   map[string]uint
	shownQuestions map[string][]uuid.UUID
}

func NewSurveyService(questionOps *question.Ops, surveyOps *survey.Ops) *SurveyService {
	return &SurveyService{
		questionOps:    questionOps,
		surveyOps:      surveyOps,
		userProgress:   make(map[string]uint),
		shownQuestions: make(map[string][]uuid.UUID),
	}
}

func generateProgressKey(userID string, surveyID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", userID, surveyID.String())
}

func (s *SurveyService) CreateQuestion(ctx context.Context, question *question.Question) error {
	return s.questionOps.Create(ctx, question)
}

func (s *SurveyService) GetNextQuestion(ctx context.Context, surveyID uuid.UUID, userID string) (*question.Question, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	survey, err := s.surveyOps.GetByID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	questions, err := s.questionOps.GetQuestionsBySurveyID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	progressKey := generateProgressKey(userID, surveyID)
	shownQuestionsKey := progressKey

	if len(questions) == 0 {
		return nil, fmt.Errorf("no questions available for this survey")
	}

	var nextQuestion *question.Question

	if survey.IsRandom {
		notShownQuestions := getNotShownQuestions(questions, s.shownQuestions[shownQuestionsKey])

		if len(notShownQuestions) == 0 {
			return nil, fmt.Errorf("all questions have been shown")
		}

		nextQuestion = notShownQuestions[rand.Intn(len(notShownQuestions))]
	} else {
		currentIndex := s.userProgress[progressKey]
		if int(currentIndex) >= len(questions) {
			return nil, fmt.Errorf("no more questions")
		}
		nextQuestion = questions[currentIndex]
		s.userProgress[progressKey]++
	}

	s.shownQuestions[shownQuestionsKey] = append(s.shownQuestions[shownQuestionsKey], nextQuestion.ID)

	return nextQuestion, nil
}

func (s *SurveyService) GetPreviousQuestion(ctx context.Context, surveyID uuid.UUID, userID string) (*question.Question, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	survey, err := s.surveyOps.GetByID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	if !survey.AllowBack {
		return nil, fmt.Errorf("going back is not allowed for this survey")
	}

	progressKey := generateProgressKey(userID, surveyID)
	shownQuestionsKey := progressKey

	shownQuestions := s.shownQuestions[shownQuestionsKey]
	if len(shownQuestions) <= 1 {
		return nil, fmt.Errorf("no previous questions available")
	}

	var previousQuestionID uuid.UUID

	if survey.IsRandom {
		previousQuestionID = shownQuestions[len(shownQuestions)-2]
		s.shownQuestions[shownQuestionsKey] = shownQuestions[:len(shownQuestions)-1]
	} else {
		currentIndex := s.userProgress[progressKey]
		if currentIndex <= 1 {
			return nil, fmt.Errorf("no previous questions available")
		}
		s.userProgress[progressKey]--
		previousQuestionID = shownQuestions[currentIndex-2]
	}

	previousQuestion, err := s.questionOps.GetByID(ctx, previousQuestionID)
	if err != nil {
		return nil, err
	}

	return previousQuestion, nil
}

func (s *SurveyService) ResetUserProgress(userID string, surveyID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	progressKey := generateProgressKey(userID, surveyID)
	delete(s.userProgress, progressKey)
	delete(s.shownQuestions, progressKey)
}

func getNotShownQuestions(allQuestions []*question.Question, shownQuestionIDs []uuid.UUID) []*question.Question {
	notShown := make([]*question.Question, 0)
	shownMap := make(map[uuid.UUID]bool)

	for _, id := range shownQuestionIDs {
		shownMap[id] = true
	}

	for _, q := range allQuestions {
		if !shownMap[q.ID] {
			notShown = append(notShown, q)
		}
	}

	return notShown
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

func (s *SurveyService) GetUserSurvey(ctx context.Context, userId uuid.UUID, page, pageSize int) ([]survey.Survey, int64, error) {
	return s.surveyOps.GetUserSurvey(ctx, userId, page, pageSize)
}
