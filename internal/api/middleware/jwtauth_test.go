package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtMocks "github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestJWTAuth tests the JWTAuth middleware handler.
// It uses table-driven tests to cover:
//   - Missing Authorization header
//   - Invalid header format (not "Bearer <token>")
//   - Invalid token (validation failure)
//   - Valid token (claims stored in context)
func TestJWTAuth(t *testing.T) {
	t.Parallel()

	// Disable Gin debug mode for cleaner test output
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		// Request setup
		authHeader string

		// Mock setup
		setupMock func(*jwtMocks.JWTValidator)

		// Expected response
		expectedStatus int
		expectedBody   map[string]any

		// Whether claims should be stored in context
		expectClaims bool
	}{
		{
			name:       "error - missing Authorization header",
			authHeader: "",
			// No mock setup needed - middleware rejects before token validation
			setupMock:      func(m *jwtMocks.JWTValidator) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Authorization header is required",
			},
			expectClaims: false,
		},
		{
			name:       "error - invalid header format (no Bearer prefix)",
			authHeader: "InvalidToken123",
			// No mock setup needed - middleware rejects before token validation
			setupMock:      func(m *jwtMocks.JWTValidator) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Invalid authorization header format",
			},
			expectClaims: false,
		},
		{
			name:       "error - invalid header format (wrong prefix)",
			authHeader: "Basic sometoken",
			// No mock setup needed - middleware rejects before token validation
			setupMock:      func(m *jwtMocks.JWTValidator) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Invalid authorization header format",
			},
			expectClaims: false,
		},
		{
			name:       "error - token validation fails",
			authHeader: "Bearer invalid.jwt.token",
			setupMock: func(m *jwtMocks.JWTValidator) {
				m.On("ValidateToken", "invalid.jwt.token").
					Return(nil, errors.New("token expired"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"error": "Invalid token",
			},
			expectClaims: false,
		},
		{
			name:       "success - valid token",
			authHeader: "Bearer valid.jwt.token",
			setupMock: func(m *jwtMocks.JWTValidator) {
				m.On("ValidateToken", "valid.jwt.token").
					Return(jwt.MapClaims{
						"sub": "user-123",
						"exp": float64(9999999999),
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil, // Handler will set its own response
			expectClaims:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create response recorder
			rec := httptest.NewRecorder()

			// Create Gin test context
			gctx, router := gin.CreateTestContext(rec)

			// Setup mock validator
			mockValidator := jwtMocks.NewJWTValidator(t)
			tc.setupMock(mockValidator)

			// Create middleware
			middleware := NewJWTAuth(mockValidator)

			// Variable to capture claims from context
			var capturedClaims any
			var claimsExist bool

			// Setup a route with the middleware
			router.GET("/test", middleware.JWTAuth(), func(c *gin.Context) {
				capturedClaims, claimsExist = c.Get("claims")
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create and execute request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			gctx.Request = req

			// Serve the request
			router.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert claims in context
			if tc.expectClaims {
				assert.True(t, claimsExist, "Expected claims to be stored in context")
				assert.NotNil(t, capturedClaims, "Expected claims to not be nil")
			}

			// Assert response body for error cases
			if tc.expectedBody != nil {
				assert.JSONEq(t,
					mustMarshal(tc.expectedBody),
					rec.Body.String(),
				)
			}
		})
	}
}

// TestNewJWTAuth tests the JWTAuth constructor.
func TestNewJWTAuth(t *testing.T) {
	t.Parallel()

	mockValidator := jwtMocks.NewJWTValidator(t)
	middleware := NewJWTAuth(mockValidator)

	assert.NotNil(t, middleware, "Expected middleware to be created")
	assert.IsType(t, &jwtAuth{}, middleware, "Expected *jwtAuth type")
}

// mustMarshal is a helper to marshal JSON for test assertions
func mustMarshal(v map[string]any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
