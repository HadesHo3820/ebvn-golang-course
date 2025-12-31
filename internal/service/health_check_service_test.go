// Package service provides unit tests for the service layer components.
//
// This file contains tests for the HealthCheckService, ensuring that the
// health check functionality works correctly across different configurations.
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testConnectionErr = errors.New("connection refused")

// TestHealthCheckService_Check validates the Check method of the HealthCheckService.
//
// This test uses a table-driven approach to verify that the health check service
// correctly returns the expected health status message, service name, instance ID,
// and handles dependency health checks appropriately.
//
// Test coverage includes:
//   - Verifying healthy status when Redis is available
//   - Verifying unhealthy status when Redis is unavailable
//   - Confirming behavior when no health checker is configured
func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                string
		inputServiceName    string
		inputInstanceID     string
		setupMock           func(t *testing.T) repository.HealthChecker
		expectedMessage     string
		expectedServiceName string
		expectedInstanceID  string
		expectedErr         error
	}{
		{
			name:             "healthy - Redis available",
			inputServiceName: "bookmark_service",
			inputInstanceID:  "instance-123",
			setupMock: func(t *testing.T) repository.HealthChecker {
				m := mocks.NewHealthChecker(t)
				m.On("Ping", mock.Anything).Return(nil).Once()
				return m
			},
			expectedMessage:     HealthCheckOK,
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-123",
			expectedErr:         nil,
		},
		{
			name:             "unhealthy - Redis unavailable",
			inputServiceName: "bookmark_service",
			inputInstanceID:  "instance-456",
			setupMock: func(t *testing.T) repository.HealthChecker {
				m := mocks.NewHealthChecker(t)
				m.On("Ping", mock.Anything).Return(testConnectionErr).Once()
				return m
			},
			expectedMessage:     HealthCheckUnhealthy,
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-456",
			expectedErr:         testConnectionErr,
		},
		{
			name:             "healthy - no health checker configured",
			inputServiceName: "bookmark_service",
			inputInstanceID:  "instance-789",
			setupMock: func(t *testing.T) repository.HealthChecker {
				return nil
			},
			expectedMessage:     HealthCheckOK,
			expectedServiceName: "bookmark_service",
			expectedInstanceID:  "instance-789",
			expectedErr:         nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup mock
			healthChecker := tc.setupMock(t)

			// Create service with mock
			testSvc := NewHealthCheck(tc.inputServiceName, tc.inputInstanceID, healthChecker)

			// Call the method
			message, serviceName, instanceID, err := testSvc.Check(ctx)

			// Assert results
			assert.Equal(t, tc.expectedMessage, message)
			assert.Equal(t, tc.expectedServiceName, serviceName)
			assert.Equal(t, tc.expectedInstanceID, instanceID)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
