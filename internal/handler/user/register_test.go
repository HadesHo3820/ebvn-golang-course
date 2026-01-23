package user

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	handlertest "github.com/HadesHo3820/ebvn-golang-course/internal/test/handler"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUserHandler_Register tests the Register handler method.
// It uses table-driven tests with mocked User service to verify:
//   - Successful registration returns user data with 200 status
//   - Duplicate username/email returns 400 status
//   - Internal server error returns 500 status
//   - Invalid request body returns 400 status (validation error)
func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	// Disable Gin debug mode for cleaner test output
	gin.SetMode(gin.TestMode)

	// Fixed timestamp for consistent test assertions
	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Default values
	defaultUsername := "testuser"
	defaultPassword := "Password1!"
	defaultDisplayName := "Test User"
	defaultEmail := "test@example.com"

	testCases := []struct {
		name string

		// Request setup
		requestBody map[string]string

		// Mock setup - returns the mocked service
		setupMockSvc func(t *testing.T, ctx context.Context) *mocks.User

		// Expected response
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:        "success - register user",
			requestBody: fixture.DefaultRegisterBody(),
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				// mock.Anything for gin.Context since it implements context.Context
				svcMock.On("CreateUser",
					ctx,
					defaultUsername, defaultPassword, defaultDisplayName, defaultEmail,
				).Return(&model.User{
					ID:          "test-uuid",
					Username:    defaultUsername,
					DisplayName: defaultDisplayName,
					Email:       defaultEmail,
					UpdatedAt:   fixedTime,
				}, nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Register an user successfully!",
				"data": map[string]any{
					"id":           "test-uuid",
					"username":     defaultUsername,
					"display_name": defaultDisplayName,
					"email":        defaultEmail,
					"updated_at":   fixedTime.String(),
				},
			},
		},
		{
			name: "error - duplicate username or email",
			requestBody: fixture.DefaultRegisterBody(
				fixture.WithField("username", "existinguser"),
				fixture.WithField("display_name", "Existing User"),
				fixture.WithField("email", "existing@example.com"),
			),
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser",
					ctx,
					"existinguser", "Password1!", "Existing User", "existing@example.com",
				).Return(nil, dbutils.ErrDuplicationType)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": "username or email is already taken",
			},
		},
		{
			name:        "error - internal server error",
			requestBody: fixture.DefaultRegisterBody(),
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser",
					ctx,
					defaultUsername, defaultPassword, defaultDisplayName, defaultEmail,
				).Return(nil, assert.AnError) // generic error
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"message": response.InternalErrMessage,
			},
		},
		{
			name:        "error - missing username",
			requestBody: fixture.DefaultRegisterBody(fixture.WithField("username", "")),
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Just check status code for validation errors
		},
		{
			name:        "error - invalid email format",
			requestBody: fixture.DefaultRegisterBody(fixture.WithField("email", "invalid-email")),
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Just check status code for validation errors
		},
		{
			name:        "error - password too short",
			requestBody: fixture.DefaultRegisterBody(fixture.WithField("password", "Pass1!")),
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Just check status code for validation errors
		},
		{
			name:        "error - username too short",
			requestBody: fixture.DefaultRegisterBody(fixture.WithField("username", "a")),
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Just check status code for validation errors
		},
		{
			name:        "error - password missing special character",
			requestBody: fixture.DefaultRegisterBody(fixture.WithField("password", "Password123")), // No special character
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when password validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Just check status code for validation errors
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test context with JSON body using helper
			testCtx := handlertest.NewTestContext(http.MethodPost, "/v1/users/register").
				WithJSONBody(tc.requestBody)

			// Setup mock service
			svcMock := tc.setupMockSvc(t, testCtx.Ctx)

			// Create handler with mock service
			handler := NewUserHandler(svcMock)

			// Call the handler
			handler.Register(testCtx.Ctx)

			// Assert response using helper
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
