// Package service provides business logic implementations for the application.
// This file contains the health check service which verifies service health
// including external dependency connectivity.
package service

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
)

const (
	// HealthCheckOK is the message returned when the service is healthy.
	HealthCheckOK = "OK"
	// HealthCheckUnhealthy is the message returned when a dependency is unhealthy.
	HealthCheckUnhealthy = "UNHEALTHY"
)

// healthCheckService holds the configuration values and dependencies needed for health checks.
type healthCheckService struct {
	serviceName   string
	instanceID    string
	healthChecker repository.HealthChecker
}

// HealthCheck defines the interface for health check operations.
//
//go:generate mockery --name HealthCheck --filename health_check_service.go
type HealthCheck interface {
	// Check verifies the health of the service and its dependencies.
	// Returns status message, service name, instance ID, and any error.
	Check(ctx context.Context) (string, string, string, error)
}

// NewHealthCheck creates a new HealthCheck service with the provided config values
// and optional health checker for verifying external dependencies.
//
// Parameters:
//   - serviceName: The name of this service.
//   - instanceID: The unique identifier for this instance.
//   - healthChecker: Optional HealthChecker for verifying dependencies (can be nil).
//
// Returns:
//   - HealthCheck: The health check service implementation.
func NewHealthCheck(serviceName, instanceID string, healthChecker repository.HealthChecker) HealthCheck {
	return &healthCheckService{
		serviceName:   serviceName,
		instanceID:    instanceID,
		healthChecker: healthChecker,
	}
}

// Check verifies the health of the service and its dependencies.
// If a HealthChecker is configured, it pings the dependency (e.g., Redis).
// Returns an error if any dependency is unhealthy.
func (s *healthCheckService) Check(ctx context.Context) (string, string, string, error) {
	// If a health checker is configured, verify the dependency
	if s.healthChecker != nil {
		if err := s.healthChecker.Ping(ctx); err != nil {
			return HealthCheckUnhealthy, s.serviceName, s.instanceID, err
		}
	}
	return HealthCheckOK, s.serviceName, s.instanceID, nil
}
