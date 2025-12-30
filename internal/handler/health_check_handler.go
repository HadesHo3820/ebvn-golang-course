// Package handler provides HTTP handlers for the bookmark service API.
//
// This package implements the presentation layer in a Hexagonal Architecture,
// handling incoming HTTP requests and delegating business logic to the service layer.
// Handlers are responsible for request parsing, response formatting, and HTTP status codes.
package handler

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

// healthCheckResponse represents the JSON response structure for health check endpoints.
//
// Note: In Go, struct fields must be exported (start with an uppercase letter) to be
// accessible via reflection, which is required for JSON serialization. The json tag
// controls the key name in the JSON output, while the example tag provides sample
// values for Swagger documentation.
type healthCheckResponse struct {
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"instance-123"`
}

// healthCheckHandler is the concrete implementation of the HealthCheck interface.
// It holds a reference to the health check service for performing health status checks.
type healthCheckHandler struct {
	healthCheckSvc service.HealthCheck
}

// HealthCheck defines the interface for health check HTTP handlers.
// This interface allows for dependency injection and easier testing.
type HealthCheck interface {
	// Check handles health check HTTP requests and returns the service status.
	Check(c *gin.Context)
}

// NewHealthCheck creates a new HealthCheck handler with the provided service dependency.
// It follows the constructor injection pattern for dependency management.
//
// Parameters:
//   - svc: The health check service that provides the actual health status logic.
//
// Returns:
//   - HealthCheck: An interface implementation that handles health check HTTP requests.
func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheckHandler{
		healthCheckSvc: svc,
	}
}

// @Summary Health check
// @Description Health check
// @Tags health_check
// @Produce json
// @Success 200 {object} healthCheckResponse
// @Failure 503 {object} map[string]string "Service Unavailable - dependency unhealthy"
// @Router /v1/health-check [get]
func (h *healthCheckHandler) Check(c *gin.Context) {
	message, serviceName, instanceID, err := h.healthCheckSvc.Check(c)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message":      message,
			"service_name": serviceName,
			"instance_id":  instanceID,
			"error":        "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, healthCheckResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
