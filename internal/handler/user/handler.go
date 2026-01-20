package user

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetSelfInfo(c *gin.Context)
}

type userHandler struct {
	svc service.User
}

func NewUserHandler(svc service.User) UserHandler {
	return &userHandler{svc: svc}
}
