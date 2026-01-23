package user

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for user related requests.
type UserHandler interface {
	// Register handles the user registration request.
	Register(c *gin.Context)
	// Login handles the user login request.
	Login(c *gin.Context)
	// GetSelfInfo handles the request to get the currently authenticated user's info.
	GetSelfInfo(c *gin.Context)
	// UpdateSelfInfo handles the request to update the currently authenticated user's info.
	UpdateSelfInfo(c *gin.Context)
}

// userHandler implements the UserHandler interface.
type userHandler struct {
	svc service.User
}

// NewUserHandler creates a new instance of UserHandler with the given service.
func NewUserHandler(svc service.User) UserHandler {
	return &userHandler{svc: svc}
}
