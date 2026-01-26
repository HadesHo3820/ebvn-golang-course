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
	"github.com/HadesHo3820/ebvn-golang-course/pkg/utils/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	// Test user constants
	testUserUsername    = "testuser"
	testUserDisplayName = "Test User"
	testUserEmail       = "test@example.com"
	testUserID          = "test-uuid"
	testHashedPassword  = "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e"

	// Existing user constants
	existingUserUsername    = "existinguser"
	existingUserDisplayName = "Existing User"
	existingUserEmail       = "existing@example.com"
	existingUserID          = "existing-user-id"

	// New user constants (for update)
	newUserDisplayName = "New Name"
	newUserEmail       = "new@example.com"
)

// TestUser_CreateUser tests the CreateUser method of the User service.
// It uses table-driven tests with mocked repository to verify:
//   - Successful user creation with hashed password
//   - Error handling when repository returns an error
func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                     string
		setupMockPasswordHashing func(t *testing.T) *mocks.PasswordHashing
		inputUsername            string
		inputPassword            string
		inputDisplay             string
		inputEmail               string
		setupMock                func(mockRepo *repoMocks.User, ctx context.Context)
		expectedErr              error
		expectedOutput           *model.User
	}{
		{
			name: "success - create user",
			setupMockPasswordHashing: func(t *testing.T) *mocks.PasswordHashing {
				hashingMock := mocks.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").Return(testHashedPassword, nil)
				return hashingMock
			},
			inputUsername: testUserUsername,
			inputPassword: "password123",
			inputDisplay:  testUserDisplayName,
			inputEmail:    testUserEmail,
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("CreateUser", ctx, &model.User{
					Username:    testUserUsername,
					Password:    testHashedPassword,
					DisplayName: testUserDisplayName,
					Email:       testUserEmail,
				}).Return(&model.User{
					Base: model.Base{
						ID: testUserID,
					},
					Username:    testUserUsername,
					DisplayName: testUserDisplayName,
					Email:       testUserEmail,
				}, nil)
			},
			expectedOutput: &model.User{
				Base: model.Base{
					ID: testUserID,
				},
				Username:    testUserUsername,
				DisplayName: testUserDisplayName,
				Email:       testUserEmail,
			},
		},
		{
			name: "error - duplicate user",
			setupMockPasswordHashing: func(t *testing.T) *mocks.PasswordHashing {
				hashingMock := mocks.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").Return(testHashedPassword, nil)
				return hashingMock
			},
			inputUsername: existingUserUsername,
			inputPassword: "password123",
			inputDisplay:  existingUserDisplayName,
			inputEmail:    existingUserEmail,
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil, dbutils.ErrDuplicationType)
			},
			expectedErr: dbutils.ErrDuplicationType,
		},
		{
			name: "error - password hashing fails",
			setupMockPasswordHashing: func(t *testing.T) *mocks.PasswordHashing {
				hashingMock := mocks.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").Return("", utils.ErrCannotGenerateHash)
				return hashingMock
			},
			inputUsername: testUserUsername,
			inputPassword: "password123",
			inputDisplay:  testUserDisplayName,
			inputEmail:    testUserEmail,
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				// No repo mock needed - should fail before reaching repo
			},
			expectedErr: utils.ErrCannotGenerateHash,
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
			mockPasswordHashing := tc.setupMockPasswordHashing(t)

			// Create service
			svc := NewUser(mockRepo, mockJWT, mockPasswordHashing)

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
	passwordHashing := utils.NewPasswordHashing()
	hashedPassword, _ := passwordHashing.Hash("correctpassword")

	testCases := []struct {
		name          string
		inputUsername string
		inputPassword string
		setupMock     func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator, mockPasswordHashing *mocks.PasswordHashing)
		expectedErr   error
		expectedToken string
	}{
		{
			name:          "success - valid login",
			inputUsername: testUserUsername,
			inputPassword: "correctpassword",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator, mockPasswordHashing *mocks.PasswordHashing) {
				mockRepo.On("GetUserByUsername", ctx, testUserUsername).Return(&model.User{
					Base: model.Base{
						ID: testUserID,
					},
					Username: "testuser",
					Password: hashedPassword,
				}, nil)
				mockPasswordHashing.On("CompareHashAndPassword", hashedPassword, "correctpassword").Return(true)
				mockJWT.On("GenerateToken", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					return claims["sub"] == testUserID
				})).Return("valid.jwt.token", nil)
			},
			expectedToken: "valid.jwt.token",
		},
		{
			name:          "error - user not found",
			inputUsername: "nonexistent",
			inputPassword: "password",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator, mockPasswordHashing *mocks.PasswordHashing) {
				mockRepo.On("GetUserByUsername", ctx, "nonexistent").Return(nil, dbutils.ErrNotFoundType)
			},
			expectedErr: dbutils.ErrNotFoundType,
		},
		{
			name:          "error - invalid password",
			inputUsername: testUserUsername,
			inputPassword: "wrongpassword",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator, mockPasswordHashing *mocks.PasswordHashing) {
				mockRepo.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					Base: model.Base{
						ID: testUserID,
					},
					Username: "testuser",
					Password: hashedPassword,
				}, nil)
				mockPasswordHashing.On("CompareHashAndPassword", hashedPassword, "wrongpassword").Return(false)
			},
			expectedErr: ErrClientErr,
		},
		{
			name:          "error - JWT generation fails",
			inputUsername: testUserUsername,
			inputPassword: "correctpassword",
			setupMock: func(ctx context.Context, mockRepo *repoMocks.User, mockJWT *jwtMocks.JWTGenerator, mockPasswordHashing *mocks.PasswordHashing) {
				mockRepo.On("GetUserByUsername", ctx, "testuser").Return(&model.User{
					Base: model.Base{
						ID: testUserID,
					},
					Username: "testuser",
					Password: hashedPassword,
				}, nil)
				mockPasswordHashing.On("CompareHashAndPassword", hashedPassword, "correctpassword").Return(true)
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
			mockPasswordHashing := mocks.NewPasswordHashing(t)
			tc.setupMock(ctx, mockRepo, mockJWT, mockPasswordHashing)

			// Create service
			svc := NewUser(mockRepo, mockJWT, mockPasswordHashing)

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
			inputUserID: existingUserID,
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("GetUserById", ctx, existingUserID).Return(&model.User{
					Base: model.Base{
						ID:        existingUserID,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					},
					Username:    "testuser",
					DisplayName: "Test User",
					Email:       "test@example.com",
				}, nil)
			},
			expectedOutput: &model.User{
				Base: model.Base{
					ID:        existingUserID,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				Username:    testUserUsername,
				DisplayName: testUserDisplayName,
				Email:       testUserEmail,
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

			//Setup mocks
			mockRepo := repoMocks.NewUser(t)
			mockJWT := jwtMocks.NewJWTGenerator(t)
			mockPasswordHashing := mocks.NewPasswordHashing(t)
			tc.setupMock(mockRepo, ctx)

			// Create service
			svc := NewUser(mockRepo, mockJWT, mockPasswordHashing)

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
			inputUserID:      testUserID,
			inputDisplayName: newUserDisplayName,
			inputEmail:       newUserEmail,
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, testUserID, newUserDisplayName, newUserEmail).Return(nil)
			},
		},
		{
			name:             "success - update display_name only",
			inputUserID:      testUserID,
			inputDisplayName: newUserDisplayName,
			inputEmail:       "",
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, testUserID, newUserDisplayName, "").Return(nil)
			},
		},
		{
			name:             "success - update email only",
			inputUserID:      testUserID,
			inputDisplayName: "",
			inputEmail:       newUserEmail,
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, testUserID, "", newUserEmail).Return(nil)
			},
		},
		{
			name:             "error - no fields provided",
			inputUserID:      testUserID,
			inputDisplayName: "",
			inputEmail:       "",
			// Do not pass fields requirements validation, so the call to UpdateUser should not be made
			setupMock:   func(mockRepo *repoMocks.User, ctx context.Context) {},
			expectedErr: ErrClientNoUpdate,
		},
		{
			name:             "error - duplicate email",
			inputUserID:      testUserID,
			inputDisplayName: "",
			inputEmail:       existingUserEmail,
			setupMock: func(mockRepo *repoMocks.User, ctx context.Context) {
				mockRepo.On("UpdateUser", ctx, testUserID, "", existingUserEmail).Return(dbutils.ErrDuplicationType)
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
			mockPasswordHashing := mocks.NewPasswordHashing(t)
			tc.setupMock(mockRepo, ctx)

			// Create service
			svc := NewUser(mockRepo, mockJWT, mockPasswordHashing)

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
