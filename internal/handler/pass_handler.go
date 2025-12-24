package handler

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

// passwordHandler is the HTTP adapter for password-related operations.
// It implements the Password interface and depends on service.Password
// following the Hexagonal Architecture pattern (Dependency Injection).
type passwordHandler struct {
	svc service.Password
}

// Password defines the interface for password HTTP handlers.
// This acts as a Port in Hexagonal Architecture, allowing the router
// to depend on an abstraction rather than a concrete implementation.
type Password interface {
	// GenPass handles HTTP requests to generate a new password.
	GenPass(c *gin.Context)
}

// NewPassword creates a new password handler with the given password service.
// It accepts service.Password interface to enable loose coupling and testability.
func NewPassword(svc service.Password) Password {
	return &passwordHandler{svc: svc}
}

// GenPass handles the password generation HTTP request.
// It delegates the password generation to the service layer and returns
// the generated password as a plain text response.
//
// Response:
//   - 200 OK: Returns the generated password as plain text.
//   - 500 Internal Server Error: If password generation fails.
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
		return
	}
	c.String(http.StatusOK, pass)
}
