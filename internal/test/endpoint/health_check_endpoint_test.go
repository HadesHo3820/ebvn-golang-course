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
//   - Redis connectivity verification
//   - JSON response serialization
//
// Prerequisites:
//   - Redis must be running and accessible for the healthy test case
//
// Test coverage includes:
//   - Verifying the endpoint is correctly registered at /health-check
//   - Validating HTTP status codes (200 when Redis is available)
//   - Asserting the response message matches service.HealthCheckOK
//   - Validating service_name and instance_id are present in response
func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	cfg, err := api.NewConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	testCases := []struct {
		name           string
		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus int
		validateBody   func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "healthy - Redis available",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, service.HealthCheckOK, body["message"])
				assert.NotEmpty(t, body["service_name"])
				assert.NotEmpty(t, body["instance_id"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := tc.setupTestHTTP(api.New(cfg))

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tc.validateBody != nil {
				tc.validateBody(t, resp)
			}
		})
	}
}
