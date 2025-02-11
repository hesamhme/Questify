package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"` // UUID as primary key
	Email        string    `gorm:"size:100;unique;not null"`                        // UNIQUE, NOT NULL
	Password     string    `gorm:"size:100;not null"`                               // Password (hashed)
	NationalCode string    `gorm:"type:char(10);unique;not null"`                   // CHAR(10), UNIQUE, NOT NULL
	TfaCode      string    `gorm:"size:6"`                                          // Temporary TFA code
	TfaExpiresAt time.Time `gorm:""`                                                // Expiration time for TFA
	Role         string    `gorm:"default:user"`
	IsVerified   bool      `gorm:"type:boolean;default:false"` // New column
	CreatedAt    time.Time `gorm:"autoCreateTime"`             // Auto set creation time
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`             // Auto set update time
}
