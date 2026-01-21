package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	repoMocks "github.com/HadesHo3820/ebvn-golang-course/internal/repository/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	jwtMocks "github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUser_CreateUser tests the CreateUser method of the User service.
// It uses table-driven tests with mocked repository to verify:
//   - Successful user creation with hashed password
//   - Error handling when repository returns an error
func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		inputUsername  string
		inputPassword  string
		inputDisplay   string
		inputEmail     string
		setupMock      func(mockRepo *repoMocks.User, ctx context.Context)
		expectedErr    error
		expectedOutput *model.User
	}{
		{
			name:          "success - create user",
			inputUsername: "testuser",
			inputPassword: "password123",
			inputDisplay:  "Test User",
			inputEmail:    "test@example.com",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				// mock.MatchedBy is a custom argument matcher from testify/mock.
				// It accepts a function that returns true if the argument matches expectations.
				// This is useful when:
				//   1. The argument is a complex struct that can't be compared with simple equality
				//   2. Some fields (like Password) are transformed before being passed to the mock
				//   3. You need to verify specific properties of an argument
				//
				// Here we verify:
				//   - Username, DisplayName, Email match the expected values
				//   - Password was correctly hashed (using VerifyPassword to check
				//     that the plain-text password matches the hash)

				// mock.MatchedBy ensures the service correctly transforms the input
				// (especially password hashing) before calling the repository
				mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(u *model.User) bool {
					return u.Username == "testuser" &&
						u.DisplayName == "Test User" &&
						u.Email == "test@example.com" &&
						utils.VerifyPassword("password123", u.Password)
				})).Return(&model.User{
					ID:          "test-uuid",
					Username:    "testuser",
					DisplayName: "Test User",
					Email:       "test@example.com",
				}, nil)
			},
			expectedOutput: &model.User{
				ID:          "test-uuid",
				Username:    "testuser",
				DisplayName: "Test User",
				Email:       "test@example.com",
			},
		},
		{
			name:          "error - duplicate user",
			inputUsername: "existinguser",
			inputPassword: "password123",
			inputDisplay:  "Existing User",
			inputEmail:    "existing@example.com",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil, dbutils.ErrDuplicationType)
			},
			expectedErr: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup mocks
			mockRepo := repoMocks.NewUser(t)
			mockJWT := jwtMocks.NewJWTGenerator(t)
			tc.setupMock(mockRepo, ctx)

			// Create service
			svc := NewUser(mockRepo, mockJWT)

			// Execute
			output, err := svc.CreateUser(ctx, tc.inputUsername, tc.inputPassword, tc.inputDisplay, tc.inputEmail)

			// Assert
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, output)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput.ID, output.ID)
			assert.Equal(t, tc.expectedOutput.Username, output.Username)
			assert.Equal(t, tc.expectedOutput.DisplayName, output.DisplayName)
			assert.Equal(t, tc.expectedOutput.Email, output.Email)
		})
	}
}

// TestUser_Login tests the Login method of the User service.
// It uses table-driven tests with mocked repository and JWT generator to verify:
//   - Successful login returns a valid token
//   - Error when user is not found
//   - Error when password is invalid
//   - Error when JWT generation fails
func TestUser_Login(t *testing.T) {
	t.Parallel()

	// Pre-hash a password for testing
	hashedPassword, _ := utils.HashPassword("correctpassword")

	testCases := []struct {
		name          string
		inputUsername string
		inputPassword string
		setupMock     func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator)
		expectedErr   error
		expectedToken string
	}{
		{
			name:          "success - valid login",
			inputUsername: "testuser",
			inputPassword: "correctpassword",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator) {
				mockRepo.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					ID:       "user-uuid",
					Username: "testuser",
					Password: hashedPassword,
				}, nil)
				mockJWT.On("GenerateToken", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					return claims["sub"] == "user-uuid"
				})).Return("valid.jwt.token", nil)
			},
			expectedToken: "valid.jwt.token",
		},
		{
			name:          "error - user not found",
			inputUsername: "nonexistent",
			inputPassword: "password",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator) {
				mockRepo.On("GetUserByUsername", ctx, "nonexistent").Return(nil, dbutils.ErrNotFoundType)
			},
			expectedErr: dbutils.ErrNotFoundType,
		},
		{
			name:          "error - invalid password",
			inputUsername: "testuser",
			inputPassword: "wrongpassword",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator) {
				mockRepo.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					ID:       "user-uuid",
					Username: "testuser",
					Password: hashedPassword,
				}, nil)
			},
			expectedErr: ErrClientErr,
		},
		{
			name:          "error - JWT generation fails",
			inputUsername: "testuser",
			inputPassword: "correctpassword",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator) {
				mockRepo.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					ID:       "user-uuid",
					Username: "testuser",
					Password: hashedPassword,
				}, nil)
				mockJWT.On("GenerateToken", mock.Anything).Return("", errors.New("jwt error"))
			},
			expectedErr: errors.New("jwt error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup mocks
			mockRepo := repoMocks.NewUser(t)
			mockJWT := jwtMocks.NewJWTGenerator(t)
			tc.setupMock(ctx, mockRepo, mockJWT)

			// Create service
			svc := NewUser(mockRepo, mockJWT)

			// Execute
			token, err := svc.Login(ctx, tc.inputUsername, tc.inputPassword)

			// Assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedToken, token)
		})
	}
}

// TestUser_GetUserByID tests the GetUserByID method of the User service.
// It uses table-driven tests with mocked repository to verify:
//   - Successful retrieval of user by ID
//   - Error handling when user is not found
func TestUser_GetUserByID(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		inputUserID    string
		setupMock      func(mockRepo *repoMocks.User, ctx context.Context)
		expectedErr    error
		expectedOutput *model.User
	}{
		{
			name:        "success - user found",
			inputUserID: "existing-user-id",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("GetUserById", ctx, "existing-user-id").Return(&model.User{
					ID:          "existing-user-id",
					Username:    "testuser",
					DisplayName: "Test User",
					Email:       "test@example.com",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				}, nil)
			},
			expectedOutput: &model.User{
				ID:          "existing-user-id",
				Username:    "testuser",
				DisplayName: "Test User",
				Email:       "test@example.com",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
		},
		{
			name:        "error - user not found",
			inputUserID: "non-existent-id",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("GetUserById", ctx, "non-existent-id").Return(nil, dbutils.ErrNotFoundType)
			},
			expectedErr: dbutils.ErrNotFoundType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup mocks
			mockRepo := repoMocks.NewUser(t)
			mockJWT := jwtMocks.NewJWTGenerator(t)
			tc.setupMock(mockRepo, ctx)

			// Create service
			svc := NewUser(mockRepo, mockJWT)

			// Execute
			output, err := svc.GetUserByID(ctx, tc.inputUserID)

			// Assert
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, output)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

// TestUser_UpdateUser tests the UpdateUser method of the User service.
// It uses table-driven tests with mocked repository to verify:
//   - Successful update of user profile
//   - Error when no fields provided for update
//   - Error handling when repository returns an error
func TestUser_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		inputUserID      string
		inputDisplayName string
		inputEmail       string
		setupMock        func(mockRepo *repoMocks.User, ctx context.Context)
		expectedErr      error
	}{
		{
			name:             "success - update both fields",
			inputUserID:      "user-uuid",
			inputDisplayName: "New Name",
			inputEmail:       "new@example.com",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, "user-uuid", "New Name", "new@example.com").Return(nil)
			},
		},
		{
			name:             "success - update display_name only",
			inputUserID:      "user-uuid",
			inputDisplayName: "New Name",
			inputEmail:       "",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, "user-uuid", "New Name", "").Return(nil)
			},
		},
		{
			name:             "success - update email only",
			inputUserID:      "user-uuid",
			inputDisplayName: "",
			inputEmail:       "new@example.com",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, "user-uuid", "", "new@example.com").Return(nil)
			},
		},
		{
			name:             "error - no fields provided",
			inputUserID:      "user-uuid",
			inputDisplayName: "",
			inputEmail:       "",
			setupMock:        func(mockRepo *repoMocks.User, ctx context.Context) {},
			expectedErr:      ErrClientNoUpdate,
		},
		{
			name:             "error - duplicate email",
			inputUserID:      "user-uuid",
			inputDisplayName: "",
			inputEmail:       "existing@example.com",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, "user-uuid", "", "existing@example.com").Return(dbutils.ErrDuplicationType)
			},
			expectedErr: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup mocks
			mockRepo := repoMocks.NewUser(t)
			mockJWT := jwtMocks.NewJWTGenerator(t)
			tc.setupMock(mockRepo, ctx)

			// Create service
			svc := NewUser(mockRepo, mockJWT)

			// Execute
			err := svc.UpdateUser(ctx, tc.inputUserID, tc.inputDisplayName, tc.inputEmail)

			// Assert
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
