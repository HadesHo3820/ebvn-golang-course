package password

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

// PasswordHandler represents the HTTP handler for password related requests.
type PasswordHandler interface {
	// GenPass handles the password generation request.
	GenPass(c *gin.Context)
}

// passwordHandler implements the PasswordHandler interface.
type passwordHandler struct {
	svc service.Password
}

// NewPasswordHandler creates a new instance of PasswordHandler with the given service.
func NewPasswordHandler(svc service.Password) PasswordHandler {
	return &passwordHandler{svc: svc}
}
