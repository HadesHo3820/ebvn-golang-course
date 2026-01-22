package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	handlertest "github.com/HadesHo3820/ebvn-golang-course/internal/test/handler"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUserHandler_Login tests the Login handler method.
// It uses table-driven tests with mocked User service to verify:
//   - Successful login returns a JWT token with 200 status
//   - Invalid credentials return 400 status
//   - User not found returns 404 status
//   - Internal server error returns 500 status
//   - Invalid request body returns 400 status (validation error)
func TestUserHandler_Login(t *testing.T) {
	t.Parallel()

	// Disable Gin debug mode for cleaner test output
	gin.SetMode(gin.TestMode)

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
			name: "success - valid login",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				// mock.Anything for gin.Context since it implements context.Context
				// The handler calls: svc.Login(c, username, password)
				// where c is *gin.Context
				svcMock.On("Login",
					ctx, // gin.Context (implements context.Context)
					"testuser", "password123",
				).Return("valid.jwt.token", nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Logged in successfully!",
				"data":    "valid.jwt.token",
			},
		},
		{
			name: "error - invalid credentials",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "wrongpassword123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login",
					ctx, "testuser", "wrongpassword123",
				).Return("", service.ErrClientErr)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": "invalid username or password",
			},
		},
		{
			name: "error - user not found",
			requestBody: map[string]string{
				"username": "nonexistent",
				"password": "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login",
					ctx, "nonexistent", "password123",
				).Return("", dbutils.ErrNotFoundType)
				return svcMock
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]any{
				"error": "invalid username or password",
			},
		},
		{
			name: "error - internal server error",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login",
					ctx, "testuser", "password123",
				).Return("", assert.AnError) // generic error
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"message": "Processing error",
			},
		},
		{
			name: "error - missing username",
			requestBody: map[string]string{
				"password": "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			// Validation error response - checking just status code
			expectedBody: nil,
		},
		{
			name: "error - password too short",
			requestBody: map[string]string{
				"username": "testuser",
				"password": "short", // less than 8 characters
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			// Validation error response - checking just status code
			expectedBody: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test context with JSON body using helper
			testCtx := handlertest.NewTestContext(http.MethodPost, "/v1/users/login").
				WithJSONBody(tc.requestBody)

			// Setup mock service
			svcMock := tc.setupMockSvc(t, testCtx.Ctx)

			// Create handler with mock service
			handler := NewUserHandler(svcMock)

			// Call the handler
			handler.Login(testCtx.Ctx)

			// Assert response using helper
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
