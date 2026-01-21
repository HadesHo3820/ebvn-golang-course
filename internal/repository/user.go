package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"gorm.io/gorm"
)

// User defines the interface for user-related database operations.
// This interface follows the repository pattern, abstracting data access
// and enabling easier testing through mock implementations.
//go:generate mockery --name User --filename user.go
type User interface {
	// CreateUser persists a new user record to the database.
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserById(ctx context.Context, userID string) (*model.User, error)
	// UpdateUser updates the display_name and email fields by user ID.
	UpdateUser(ctx context.Context, userID, displayName, email string) error
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
		return nil, dbutils.CatchDBErr(err)
	}

	return newUser, nil
}

// GetUserByUsername retrieves a user by their username from the database.
// The operation is executed within the provided context, allowing for
// cancellation and timeout control.
//
// Parameters:
//   - ctx: Context for request cancellation and deadline control
//   - username: The username of the user to retrieve
//
// Returns:
//   - *model.User: The retrieved user with any database-generated fields populated
//   - error: An error if the database operation fails, nil otherwise
func (u *user) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.GetUserByField(ctx, "username", username)
}

func (u *user) GetUserById(ctx context.Context, userID string) (*model.User, error) {
	return u.GetUserByField(ctx, "id", userID)
}

func (u *user) GetUserByField(ctx context.Context, field string, value string) (*model.User, error) {
	chosenUser := &model.User{}
	err := u.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(chosenUser).Error
	if err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return chosenUser, nil
}

var (
	ErrNoUpdate = errors.New("no update")
)

// UpdateUser updates the display_name and email fields for a user by their ID.
// Only non-empty fields will be updated in the database.
//
// Parameters:
//   - ctx: Context for request cancellation and deadline control
//   - userID: The unique identifier of the user to update
//   - displayName: The new display name (empty string means no update)
//   - email: The new email address (empty string means no update)
//
// Returns:
//   - error: An error if the database operation fails, nil otherwise
func (u *user) UpdateUser(ctx context.Context, userID, displayName, email string) error {
	// Build updates map
	// GORM's Updates() only accepts map[string]any
	updates := make(map[string]any)

	if displayName != "" {
		updates["display_name"] = displayName
	}
	if email != "" {
		updates["email"] = email
	}

	// Perform the update
	err := u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return dbutils.CatchDBErr(err)
	}

	return nil
}
