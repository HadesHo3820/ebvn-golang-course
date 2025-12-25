package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
