package handler

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

// In Go, when you want to use json tags (or any other reflection-based serialization),
// struct fields must be exported (start with an uppercase letter).
// Because Go's reflection cannot access unexported fields from other packages

// The json tag only controls the JSON key name in the output
type healthCheckResponse struct {
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"instance-123"`
}

type healthCheckHandler struct {
	healthCheckSvc service.HealthCheck
}

type HealthCheck interface {
	Check(c *gin.Context)
}

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
// @Failure 500 {string} Internal Server Error
// @Router /health-check [get]
func (h *healthCheckHandler) Check(c *gin.Context) {
	message, serviceName, instanceID := h.healthCheckSvc.Check()
	c.JSON(http.StatusOK, healthCheckResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
