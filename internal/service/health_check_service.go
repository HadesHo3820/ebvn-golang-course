package service

const (
	HealthCheckOK = "OK"
)

// healthCheckService holds the configuration values needed for health checks.
type healthCheckService struct {
	serviceName string
	instanceID  string
}

// HealthCheck defines the interface for health check operations.
//go:generate mockery --name HealthCheck --filename health_check_service.go
type HealthCheck interface {
	Check() (string, string, string)
}

// NewHealthCheck creates a new HealthCheck service with the provided config values.
// This follows the Dependency Injection pattern - config is passed in from outside.
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// Check returns the service name, instance ID, and app port.
func (s *healthCheckService) Check() (string, string, string) {
	return HealthCheckOK, s.serviceName, s.instanceID
}
