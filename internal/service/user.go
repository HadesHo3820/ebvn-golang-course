// Package service contains the business logic layer of the application.
// It orchestrates data flow between handlers and repositories,
// implementing domain-specific rules and transformations.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenLast = 24 * time.Hour
)

// User defines the interface for user-related business operations.
// This abstraction allows for easy testing and swapping of implementations.
//
//go:generate mockery --name User --filename user.go
type User interface {
	// CreateUser registers a new user with the provided credentials and profile information.
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	UpdateUser(ctx context.Context, userID, displayName, email string) error
}

// user is the concrete implementation of the User interface.
// It coordinates between the repository layer and applies business rules.
type user struct {
	repo            repository.User
	jwtGen          jwtutils.JWTGenerator
	passwordHashing utils.PasswordHashing
}

// NewUser creates and returns a new User service instance.
// It requires a User repository for data persistence operations.
//
// Parameters:
//   - repo: A repository.User implementation for database operations
//
// Returns:
//   - User: An implementation of the User service interface
func NewUser(repo repository.User, jwtGen jwtutils.JWTGenerator, hash utils.PasswordHashing) User {
	return &user{repo: repo, jwtGen: jwtGen, passwordHashing: hash}
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
	hashPwd, err := u.passwordHashing.Hash(password)
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

var (
	ErrClientErr      = errors.New("invalid username or password")
	ErrClientNoUpdate = errors.New("no update")
)

func (u *user) Login(ctx context.Context, username, password string) (string, error) {
	// check if user exist
	chosenUser, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	// check if password is valid
	check := u.passwordHashing.CompareHashAndPassword(chosenUser.Password, password)
	if !check {
		return "", ErrClientErr
	}

	// create token
	jwtContent := jwt.MapClaims{
		"sub": chosenUser.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenLast).Unix(),
	}
	token, err := u.jwtGen.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	// return token
	return token, nil
}

func (u *user) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	return u.repo.GetUserById(ctx, userId)
}

func (u *user) UpdateUser(ctx context.Context, userID, displayName, email string) error {
	if displayName == "" && email == "" {
		return ErrClientNoUpdate
	}
	return u.repo.UpdateUser(ctx, userID, displayName, email)
}
