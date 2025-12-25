// Package service provides unit tests for the service layer components.
//
// This file contains tests for the HealthCheckService, ensuring that the
// health check functionality works correctly across different configurations.
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHealthCheckService_Check validates the Check method of the HealthCheckService.
//
// This test uses a table-driven approach to verify that the health check service
// correctly returns the expected health status message, service name, and instance ID.
//
// Test coverage includes:
//   - Verifying that the service returns HealthCheckOK message
//   - Confirming that the service name is correctly propagated
//   - Ensuring the instance ID matches the configured value
//
// The test runs in parallel for improved performance.
func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                string
		inputServiceName    string
		inputInstanceID     string
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
	}{
		{
			name:                "normal case",
			inputServiceName:    "bookmark_service",
			inputInstanceID:     "instance-123",
			expectedMessage:     HealthCheckOK,
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Inject test values directly - no real config needed!
			testSvc := NewHealthCheck(tc.inputServiceName, tc.inputInstanceID)

			// Call the method
			message, serviceName, instanceID := testSvc.Check()

			assert.Equal(t, tc.expectedMessage, message)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)
		})
	}
}
