package entities

import (
	"github.com/google/uuid"
	"time"
)

type Question struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Index       uint      `gorm:"type:integer;not null"`
	SurveyID    uuid.UUID `gorm:"type:uuid;not null"`
	Survey      Survey    `gorm:"foreignKey:SurveyID;references:ID"`
	Text        string    `gorm:"type:text;not null"`
	Type        uint      `gorm:"type:smallint;not null"`
	IsMandatory bool      `gorm:"type:boolean;not null"`
	MediaPath   string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null"`
}

type QuestionChoices struct {
	ID         uint      `gorm:"type:int;primaryKey"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"`
	Question   Question  `gorm:"foreignKey:QuestionID;references:ID"`
	Value      string    `gorm:"type:text;not null"`
	IsAnswer   bool      `gorm:"type:boolean;not null"`
}
