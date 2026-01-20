package password

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

type PasswordHandler interface {
	GenPass(c *gin.Context)
}

type passwordHandler struct {
	svc service.Password
}

func NewPasswordHandler(svc service.Password) PasswordHandler {
	return &passwordHandler{svc: svc}
}
