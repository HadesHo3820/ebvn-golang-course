// Package handler provides unit tests for the HTTP handler layer.
//
// This file contains tests for the healthCheckHandler, using mocks to isolate
// the handler from its service dependencies and verify HTTP response behavior.
package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestHealthCheckHandler_Check validates the Check method of the healthCheckHandler.
//
// This test uses a table-driven approach with the following testing patterns:
//   - Mock injection: Uses mockery-generated mocks to isolate the handler from the service layer
//   - HTTP simulation: Uses httptest.NewRecorder and gin.CreateTestContext for request/response testing
//   - Parallel execution: Runs test cases concurrently for improved performance
//
// Test coverage includes:
//   - Verifying correct HTTP 200 status when healthy
//   - Verifying correct HTTP 503 status when dependency is unhealthy
//   - Validating JSON response body structure and content
func TestHealthCheckHandler_Check(t *testing.T) {
	t.Parallel()

	// Set Gin to test mode to reduce noise in test output
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupRequest   func(ctx *gin.Context)
		setupMockSvc   func(t *testing.T) *mocks.HealthCheck
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "healthy - all dependencies available",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", mock.Anything).
					Return("OK", "bookmark_service", "instance-123", nil).Once()
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"OK","service_name":"bookmark_service","instance_id":"instance-123"}`,
		},
		{
			name: "unhealthy - Redis unavailable",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", mock.Anything).
					Return("UNHEALTHY", "bookmark_service", "instance-456", errors.New("connection refused")).Once()
				return mockSvc
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   `{"error":"connection refused","instance_id":"instance-456","message":"UNHEALTHY","service_name":"bookmark_service"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a httptest.NewRecorder to capture status code and response
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
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())

			// Verify mock expectations
			svcMock.AssertExpectations(t)
		})
	}
}
