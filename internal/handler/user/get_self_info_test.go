package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestUserHandler_GetSelfInfo tests the GetSelfInfo handler method.
// It uses table-driven tests with mocked User service to verify:
//   - Successful profile retrieval returns user data with 200 status
//   - Missing/invalid JWT token returns 401 status
//   - User not found returns 404 status
//   - Internal server error returns 500 status
func TestUserHandler_GetSelfInfo(t *testing.T) {
	t.Parallel()

	// Disable Gin debug mode for cleaner test output
	gin.SetMode(gin.TestMode)

	// Fixed timestamp for consistent test assertions
	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name string

		// JWT claims to set in context (nil means no claims)
		jwtClaims jwt.MapClaims

		// Mock setup - returns the mocked service
		setupMockSvc func(t *testing.T, ctx context.Context) *mocks.User

		// Expected response
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success - get user profile",
			jwtClaims: jwt.MapClaims{
				"sub": "test-user-id",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("GetUserByID",
					ctx,
					"test-user-id",
				).Return(&model.User{
					ID:          "test-user-id",
					Username:    "testuser",
					DisplayName: "Test User",
					Email:       "test@example.com",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				}, nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"data": map[string]any{
					"id":           "test-user-id",
					"username":     "testuser",
					"display_name": "Test User",
					"email":        "test@example.com",
					"created_at":   fixedTime.Format(time.RFC3339Nano),
					"updated_at":   fixedTime.Format(time.RFC3339Nano),
				},
			},
		},
		{
			name:      "error - missing JWT claims",
			jwtClaims: nil, // No claims set
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
			name: "error - empty uid in claims",
			jwtClaims: jwt.MapClaims{
				"sub": "", // Empty uid
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				// Service should not be called when uid is empty
				return mocks.NewUser(t)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"message": "Invalid token",
			},
		},
		{
			name: "error - user not found",
			jwtClaims: jwt.MapClaims{
				"sub": "non-existent-id",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("GetUserByID",
					ctx,
					"non-existent-id",
				).Return(nil, dbutils.ErrNotFoundType)
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
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("GetUserByID",
					ctx,
					"test-user-id",
				).Return(nil, assert.AnError) // generic error
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"message": "Processing error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a response recorder to capture the response
			rec := httptest.NewRecorder()

			// Create a Gin test context
			gctx, _ := gin.CreateTestContext(rec)

			// Create request
			gctx.Request = httptest.NewRequest(
				http.MethodGet,
				"/v1/self/info",
				nil,
			)

			// Set JWT claims in context (simulates JWT middleware)
			if tc.jwtClaims != nil {
				gctx.Set("claims", tc.jwtClaims)
			}

			// Setup mock service
			svcMock := tc.setupMockSvc(t, gctx)

			// Create handler with mock service
			handler := NewUserHandler(svcMock)

			// Call the handler
			handler.GetSelfInfo(gctx)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body
			if tc.expectedBody != nil {
				var actualBody map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, actualBody)
			}
		})
	}
}
