package password

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GenPass handles the password generation HTTP request.
// It delegates the password generation to the service layer and returns
// the generated password as a plain text response.
//
// @Summary Generate a random password
// @Description Generates a cryptographically secure random password
// @Tags password
// @Produce plain
// @Success 200 {string} string "Generated password"
// @Failure 500 {object} response.Message
// @Router /v1/gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate password")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	c.String(http.StatusOK, pass)
}
