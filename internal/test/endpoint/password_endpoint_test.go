// Package endpoint contains integration tests for the API endpoints.
// Unlike unit tests that test components in isolation with mocks,
// integration tests verify the complete request-response cycle through
// the entire application stack (routing → handlers → services → response).
//
// These tests use httptest.NewRequest and httptest.NewRecorder to simulate
// HTTP requests without starting an actual HTTP server, making them fast
// and suitable for CI/CD pipelines while still testing real integrations.
//
// Key differences from unit tests:
//   - Uses real service implementations (no mocks)
//   - Tests the full routing configuration
//   - Validates middleware behavior
//   - Ensures all layers work together correctly
package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestPasswordEndpoint is an integration test for the password generation endpoint.
// It tests the complete flow from HTTP request to response, including:
//   - Route registration ("/gen-pass")
//   - Handler execution
//   - Service layer (real password generation)
//   - Response formatting
//
// This test uses a table-driven approach with the following structure:
//   - name: Human-readable test case identifier
//   - setupTestHTTP: Function that creates the request and captures the response
//   - expectedStatus: Expected HTTP status code
//   - expectedRespLen: Expected length of the response body (password length)
//
// The test runs in parallel for improved performance.
func TestPasswordEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus  int
		expectedRespLen int
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/gen-pass", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus:  http.StatusOK,
			expectedRespLen: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := tc.setupTestHTTP(api.New(gin.New(), &api.Config{}, redisPkg.InitMockRedis(t), nil, nil))

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedRespLen, len(rec.Body.String()))
		})
	}
}
