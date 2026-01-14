package repository

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"gorm.io/gorm"
)

// User defines the interface for user-related database operations.
// This interface follows the repository pattern, abstracting data access
// and enabling easier testing through mock implementations.
type User interface {
	// CreateUser persists a new user record to the database.
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
}

// user is the concrete implementation of the User interface.
// It uses GORM as the underlying ORM for database operations.
type user struct {
	db *gorm.DB
}

// NewUser creates and returns a new User repository instance.
// It requires a GORM database connection which will be used for all
// subsequent database operations.
//
// Parameters:
//   - db: A GORM database connection instance
//
// Returns:
//   - User: An implementation of the User interface
func NewUser(db *gorm.DB) User {
	return &user{db: db}
}

// CreateUser inserts a new user record into the database.
// The operation is executed within the provided context, allowing for
// cancellation and timeout control.
//
// Parameters:
//   - ctx: Context for request cancellation and deadline control
//   - newUser: Pointer to the User model containing the data to be inserted
//
// Returns:
//   - *model.User: The created user with any database-generated fields populated
//   - error: An error if the database operation fails, nil otherwise
func (u *user) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := u.db.WithContext(ctx).Create(&newUser).Error
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
