package entities

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	QuestionID Question `gorm:"foreignKey:ID;references:ID"` 
	UserID     User `gorm:"foreignKey:UserID;references:ID"`
	Response   string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
