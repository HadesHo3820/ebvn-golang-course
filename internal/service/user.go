// Package service contains the business logic layer of the application.
// It orchestrates data flow between handlers and repositories,
// implementing domain-specific rules and transformations.
package service

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/utils"
)

// User defines the interface for user-related business operations.
// This abstraction allows for easy testing and swapping of implementations.
type User interface {
	// CreateUser registers a new user with the provided credentials and profile information.
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

// user is the concrete implementation of the User interface.
// It coordinates between the repository layer and applies business rules.
type user struct {
	repo repository.User
}

// NewUser creates and returns a new User service instance.
// It requires a User repository for data persistence operations.
//
// Parameters:
//   - repo: A repository.User implementation for database operations
//
// Returns:
//   - User: An implementation of the User service interface
func NewUser(repo repository.User) User {
	return &user{repo: repo}
}

// CreateUser handles user registration by hashing the password and persisting the user.
// It applies security measures (password hashing) before delegating to the repository.
//
// Parameters:
//   - ctx: Context for request cancellation and deadline control
//   - username: Unique login identifier for the user
//   - password: Plain-text password (will be hashed before storage)
//   - displayName: User's display name shown in the UI
//   - email: User's email address
//
// Returns:
//   - *model.User: The created user with database-generated fields (e.g., ID)
//   - error: An error if password hashing or database operation fails
func (u *user) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	newUser := &model.User{
		Username:    username,
		Password:    hashPwd,
		DisplayName: displayName,
		Email:       email,
	}

	res, err := u.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return res, nil
}
