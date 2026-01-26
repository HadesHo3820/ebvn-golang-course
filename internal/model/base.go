// Package model defines the domain entities and data structures
// used throughout the application. These models represent database
// tables and are used for data transfer between layers.
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common fields shared across all domain models.
// It provides standard auditing fields (ID, timestamps) and soft delete support.
// All models should embed this struct to ensure consistency across the database schema.
//
// Fields:
//   - ID: Unique identifier (UUID) auto-generated before record creation
//   - CreatedAt: Timestamp when the record was created (auto-managed by GORM)
//   - UpdatedAt: Timestamp when the record was last updated (auto-managed by GORM)
//   - DeletedAt: Timestamp for soft deletion (NULL means active, non-NULL means deleted)
type Base struct {
	ID        string    `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
}

// BeforeCreate is a GORM hook that runs automatically before inserting a new record.
// It generates a UUID for the ID field if one is not already set, ensuring every
// record has a unique identifier without requiring manual ID assignment.
//
// This hook is triggered for any model that embeds the Base struct.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}
