// Package model defines the domain entities and data structures
// used throughout the application. These models represent database
// tables and are used for data transfer between layers.
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user account in the system.
// This struct maps to the "users" table in the database and contains
// all user-related information including authentication credentials.
//
// Fields:
//   - ID: Unique identifier (UUID) for the user, serves as primary key
//   - Username: Unique login name for the user
//   - Password: Hashed password (excluded from JSON serialization for security)
//   - DisplayName: User's display name shown in the UI
//   - Email: Unique email address for the user
type User struct {
	ID          string    `gorm:"type:uuid;primarykey;column:id" json:"id"`
	Username    string    `gorm:"unique;column:username" json:"username"`
	Password    string    `gorm:"column:password" json:"-"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	Email       string    `gorm:"unique;column:email" json:"email"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before inserting a new User record.
// It automatically generates a UUID for the ID field if one is not already set,
// ensuring every user has a unique identifier.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
