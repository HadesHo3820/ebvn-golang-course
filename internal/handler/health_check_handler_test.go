// Package handler provides unit tests for the HTTP handler layer.
//
// This file contains tests for the healthCheckHandler, using mocks to isolate
// the handler from its service dependencies and verify HTTP response behavior.
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheckHandler_Check validates the Check method of the healthCheckHandler.
//
// This test uses a table-driven approach with the following testing patterns:
//   - Mock injection: Uses mockery-generated mocks to isolate the handler from the service layer
//   - HTTP simulation: Uses httptest.NewRecorder and gin.CreateTestContext for request/response testing
//   - Parallel execution: Runs test cases concurrently for improved performance
//
// Test structure for each case:
//   - setupRequest: Configures the incoming HTTP request (method, path, headers, body)
//   - setupMockSvc: Creates and configures the mock service with expected behavior
//   - expectedStatus: The expected HTTP status code
//   - expectedBody: The expected JSON response body
//
// Test coverage includes:
//   - Verifying correct HTTP status code is returned
//   - Validating JSON response body structure and content
//   - Ensuring proper delegation to the health check service
func TestHealthCheckHandler_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T) *mocks.HealthCheck
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "normal case",
			setupRequest: func(ctx *gin.Context) {
				// Create virtual request
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check").Return("OK", "bookmark_service", "instance-123")
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"OK","service_name":"bookmark_service","instance_id":"instance-123"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a httptest.NewRecorder to capture status code and response for the incomming request
			rec := httptest.NewRecorder()

			// Create a Gin test context to simulate a request
			gctx, _ := gin.CreateTestContext(rec)

			// Setup the request
			tc.setupRequest(gctx)

			// Setup the mock service
			svcMock := tc.setupMockSvc(t)

			// Create the handler with the mock service
			handler := NewHealthCheck(svcMock)

			// Call the handler
			handler.Check(gctx)

			// Check the response and status code
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, rec.Body.String())
		})
	}
}
