// Package model defines the domain entities and data structures
// used throughout the application. These models represent database
// tables and are used for data transfer between layers.
package model

// User represents a user account in the system.
// This struct maps to the "users" table in the database and contains
// all user-related information including authentication credentials.
//
// Fields:
//   - Base: Embedded struct providing ID, CreatedAt, UpdatedAt, and DeletedAt fields
//   - Username: Unique login name for the user
//   - Password: Hashed password (excluded from JSON serialization for security)
//   - DisplayName: User's display name shown in the UI
//   - Email: Unique email address for the user
type User struct {
	Base
	Username    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"-"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`
	Email       string `gorm:"unique;column:email" json:"email"`
}
