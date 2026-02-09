package healthcheck

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// healthCheckData represents the health check information.
type healthCheckData struct {
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"instance-123"`
}

type pingErrorResponse struct {
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"instance-123"`
	Error       string `json:"error"`
}

// @Summary Health check
// @Description Health check
// @Tags health_check
// @Produce json
// @Success 200 {object} dto.SuccessResponse[healthCheckData]
// @Failure 503 {object} pingErrorResponse "Service Unavailable - dependency unhealthy"
// @Router /health-check [get]
func (h *healthCheckHandler) Ping(c *gin.Context) {
	message, serviceName, instanceID, err := h.svc.Check(c)
	if err != nil {
		log.Error().
			Str("service_name", serviceName).
			Str("instance_id", instanceID).
			Err(err).
			Msg("Health check failed")
		c.JSON(http.StatusServiceUnavailable,
			pingErrorResponse{
				Message:     message,
				ServiceName: serviceName,
				InstanceID:  instanceID,
				Error:       "Internal Server Error",
			})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse[*healthCheckData]{
		Message: message,
		Data: &healthCheckData{
			ServiceName: serviceName,
			InstanceID:  instanceID,
		},
	})
}
