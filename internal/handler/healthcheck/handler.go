package healthcheck

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

type HealthCheckHandler interface {
	Ping(c *gin.Context)
}

type healthCheckHandler struct {
	svc service.HealthCheck
}

func NewHealthCheckHandler(svc service.HealthCheck) HealthCheckHandler {
	return &healthCheckHandler{svc: svc}
}
