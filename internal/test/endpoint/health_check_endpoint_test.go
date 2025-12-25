// Package endpoint provides integration tests for API endpoints.
//
// Unlike unit tests that mock dependencies, endpoint tests validate the full HTTP stack
// including routing, middleware, handlers, and real service implementations. These tests
// ensure that all layers work together correctly when processing HTTP requests.
package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheckEndpoint validates the /health-check endpoint through the full HTTP stack.
//
// This is an integration test that exercises:
//   - HTTP routing configuration
//   - Request handling through the Gin engine
//   - Handler-to-service delegation with real dependencies
//   - JSON response serialization
//
// Test methodology:
//   - Uses httptest.NewRequest to simulate incoming HTTP requests
//   - Uses httptest.NewRecorder to capture the response
//   - Calls api.Engine.ServeHTTP directly without starting a real server
//
// Test coverage includes:
//   - Verifying the endpoint is correctly registered at /health-check
//   - Validating HTTP status codes
//   - Asserting the response message matches service.HealthCheckOK
//
// The test runs in parallel for improved performance.
func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		name            string
		setupTestHTTP   func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: service.HealthCheckOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := tc.setupTestHTTP(api.New(cfg))

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]any
			json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.Equal(t, tc.expectedMessage, resp["message"])
		})
	}
}
