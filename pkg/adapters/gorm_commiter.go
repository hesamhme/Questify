package adapters

import (
	"gorm.io/gorm"
)

// GormCommitter provides transaction handling for GORM.
type GormCommitter struct {
	DB *gorm.DB
}

// WithTransaction executes the given function within a transaction.
func (gc *GormCommitter) WithTransaction(fn func(tx *gorm.DB) error) error {
	return gc.DB.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
