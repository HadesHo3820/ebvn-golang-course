package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// pingResponse represents the JSON response structure for health check endpoints.
//
// Note: In Go, struct fields must be exported (start with an uppercase letter) to be
// accessible via reflection, which is required for JSON serialization. The json tag
// controls the key name in the JSON output, while the example tag provides sample
// values for Swagger documentation.
type pingResponse struct {
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"bookmark_service"`
	InstanceID  string `json:"instance_id" example:"instance-123"`
}

type pingErrorResponse struct {
	pingResponse
	Error string `json:"error"`
}

// @Summary Health check
// @Description Health check
// @Tags health_check
// @Produce json
// @Success 200 {object} pingResponse
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
				pingResponse: pingResponse{
					Message:     message,
					ServiceName: serviceName,
					InstanceID:  instanceID,
				},
				Error: "Internal Server Error",
			})
		return
	}
	c.JSON(http.StatusOK, pingResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
