package healthcheck

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

// HealthCheckHandler represents the HTTP handler for health check requests.
type HealthCheckHandler interface {
	// Ping handles the GET /ping request to check the service health.
	Ping(c *gin.Context)
}

// healthCheckHandler implements the HealthCheckHandler interface.
type healthCheckHandler struct {
	svc service.HealthCheck
}

// NewHealthCheckHandler creates a new instance of HealthCheckHandler with the given service.
func NewHealthCheckHandler(svc service.HealthCheck) HealthCheckHandler {
	return &healthCheckHandler{svc: svc}
}
