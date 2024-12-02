package entities

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Survey struct {
	ID                 uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title              string        `gorm:"type:text;not null"`
	OwnerID            uuid.UUID     `gorm:"type:uuid;not null"`
	StartTime          time.Time     `gorm:"type:timestamp;not null"`
	EndTime            time.Time     `gorm:"type:timestamp;not null"`
	IsRandom           bool          `gorm:"type:boolean;default:false;not null"`
	IsCanceled         bool          `gorm:"type:boolean;default:false;not null"`
	AllowBack          bool          `gorm:"type:boolean;default:true;not null"`
	ParticipationLimit uint          `gorm:"type:integer;not null"`
	ResponseTimeLimit  time.Duration `gorm:"type:integer"`
	CreatedAt          time.Time     `gorm:"type:timestamp;not null"`
}

type SurveyRequirements struct {
	ID         uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SurveyID   uuid.UUID     `gorm:"type:uuid;not null"`
	Survey     Survey        `gorm:"foreignKey:SurveyID;references:ID"`
	MinAge     sql.NullInt32 `gorm:"type:integer;"`
	MaxAge     sql.NullInt32 `gorm:"type:integer;"`
	CityId     sql.NullInt32 `gorm:"type:int;not null"`
	City       City          `gorm:"foreignKey:CityId;references:ID"`
	Gender     string        `gorm:"type:text;"`
	ReviewTime time.Time     `gorm:"type:timestamp;"`
}
