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
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestUserHandler_UpdateSelfInfo tests the UpdateSelfInfo handler method.
// It uses table-driven tests with mocked User service to verify:
//   - Successful profile update returns 200 status
//   - Missing/invalid JWT token returns 401 status
//   - No update data provided returns 400 status
//   - User not found returns 404 status
//   - Internal server error returns 500 status
//   - Invalid email format returns 400 status
func TestUserHandler_UpdateSelfInfo(t *testing.T) {
	t.Parallel()

	// Disable Gin debug mode for cleaner test output
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		// JWT claims to set in context (nil means no claims)
		jwtClaims jwt.MapClaims

		// Request body
		requestBody map[string]string

		// Mock setup - returns the mocked service
		setupMockSvc func(t *testing.T, ctx context.Context) *mocks.User

		// Expected response
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success - update both fields",
			jwtClaims: jwt.MapClaims{
				"sub": "test-user-id",
			},
			requestBody: map[string]string{
				"display_name": "New Display Name",
				"email":        "new@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUser",
					ctx,
					"test-user-id", "New Display Name", "new@example.com",
				).Return(nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Edit current user successfully!",
			},
		},
		{
			name: "success - update display_name only",
			jwtClaims: jwt.MapClaims{
				"sub": "test-user-id",
			},
			requestBody: map[string]string{
				"display_name": "New Display Name",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUser",
					ctx,
					"test-user-id", "New Display Name", "",
				).Return(nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Edit current user successfully!",
			},
		},
		{
			name: "success - update email only",
			jwtClaims: jwt.MapClaims{
				"sub": "test-user-id",
			},
			requestBody: map[string]string{
				"email": "new@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUser",
					ctx,
					"test-user-id", "", "new@example.com",
				).Return(nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Edit current user successfully!",
			},
		},
		{
			name:      "error - missing JWT claims",
			jwtClaims: nil, // No claims set
			requestBody: map[string]string{
				"display_name": "New Display Name",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when JWT is invalid
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"message": "Invalid token",
			},
		},
		{
			name: "error - no update data provided",
			jwtClaims: jwt.MapClaims{
				"sub": "test-user-id",
			},
			requestBody: map[string]string{}, // Empty body
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUser",
					ctx,
					"test-user-id", "", "",
				).Return(service.ErrClientNoUpdate)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": "No data provided for update. Please provide at least one field to update.",
			},
		},
		{
			name: "error - user not found",
			jwtClaims: jwt.MapClaims{
				"sub": "non-existent-id",
			},
			requestBody: map[string]string{
				"display_name": "New Display Name",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUser",
					ctx,
					"non-existent-id", "New Display Name", "",
				).Return(dbutils.ErrNotFoundType)
				return svcMock
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]any{
				"message": "User does not exist",
			},
		},
		{
			name: "error - internal server error",
			jwtClaims: jwt.MapClaims{
				"sub": "test-user-id",
			},
			requestBody: map[string]string{
				"display_name": "New Display Name",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUser",
					ctx,
					"test-user-id", "New Display Name", "",
				).Return(assert.AnError) // generic error
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"message": "Processing error",
			},
		},
		{
			name: "error - invalid email format",
			jwtClaims: jwt.MapClaims{
				"sub": "test-user-id",
			},
			requestBody: map[string]string{
				"email": "invalid-email",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when validation fails
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Just check status code for validation errors
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test context with JSON body and JWT claims using helper
			testCtx := handlertest.NewTestContext(http.MethodPut, "/v1/self/info").
				WithJSONBody(tc.requestBody).
				WithJWTClaims(tc.jwtClaims)

			// Setup mock service
			svcMock := tc.setupMockSvc(t, testCtx.Ctx)

			// Create handler with mock service
			handler := NewUserHandler(svcMock)

			// Call the handler
			handler.UpdateSelfInfo(testCtx.Ctx)

			// Assert response using helper
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
