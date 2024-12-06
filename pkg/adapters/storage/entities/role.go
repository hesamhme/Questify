package entities

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role entity in the database.
type Role struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string     `gorm:"type:varchar(100);unique;not null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}

// Permission represents a permission entity in the database.
type Permission struct {
	ID          int    `gorm:"primaryKey;autoIncrement:false"`
	Description string `gorm:"type:varchar(255);not null"`
}

// UserRole represents the relationship between a user and a role.
type UserRole struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	RoleID    uuid.UUID  `gorm:"type:uuid;not null"`
	AssignedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt *time.Time `gorm:"type:timestamp"`
	Role      Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}
