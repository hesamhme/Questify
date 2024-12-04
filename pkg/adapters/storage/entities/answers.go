package entities

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null"` // Correctly reference the Question's ID
	UserID     uuid.UUID `gorm:"type:uuid;not null"` // Correctly reference the User's ID
	Response   string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
