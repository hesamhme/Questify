package presenter

import (
	"Questify/internal/survey"
	"github.com/google/uuid"
	"log"
	"time"
)

type Survey struct {
	ID                 uuid.UUID `json:"id"`
	Title              string    `json:"title"`
	OwnerID            uuid.UUID `json:"owner_id"`
	StartTime          string    `json:"start_time"`
	EndTime            string    `json:"end_time"`
	IsRandom           bool      `json:"is_random"`
	IsCanceled         bool      `json:"is_canceled"`
	AllowBack          bool      `json:"allow_back"`
	ParticipationLimit uint      `json:"participation_limit"`
	ResponseTimeLimit  uint      `json:"response_time_limit"`
}

func MapPresenterToSurvey(presenterSurvey *Survey) *survey.Survey {
	return &survey.Survey{
		ID:                 presenterSurvey.ID,
		Title:              presenterSurvey.Title,
		OwnerID:            presenterSurvey.OwnerID,
		StartTime:          parseTime(presenterSurvey.StartTime),
		EndTime:            parseTime(presenterSurvey.EndTime),
		IsRandom:           presenterSurvey.IsRandom,
		IsCanceled:         presenterSurvey.IsCanceled,
		AllowBack:          presenterSurvey.AllowBack,
		ParticipationLimit: presenterSurvey.ParticipationLimit,
		ResponseTimeLimit:  toDuration(presenterSurvey.ResponseTimeLimit),
	}
}

func parseTime(timeStr string) time.Time {
	const isoLayout = "2006-01-02T15:04:05Z" // Example: "2023-12-01T14:00:00Z"
	parsedTime, err := time.Parse(isoLayout, timeStr)
	if err != nil {
		log.Fatalf("Error parsing time string: %s, error: %v", timeStr, err)
	}
	return parsedTime
}
func toDuration(seconds uint) time.Duration { return time.Duration(seconds) * time.Second }
