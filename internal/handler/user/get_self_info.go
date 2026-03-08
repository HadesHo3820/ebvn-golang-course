package user

import (
	"errors"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetSelfInfo handles user self info retrieval requests.
// It extracts the user ID from the JWT token context and returns the user's profile.
//
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse[model.User]
// @Failure 401 {object} response.Message "Invalid or missing token"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/self/info [get]
func (u *userHandler) GetSelfInfo(c *gin.Context) {
	// Get user id from JWT token
	uid, err := utils.GetUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid token",
		})
		return
	}

	// call service
	user, err := u.svc.GetUserByID(c, uid)
	switch {
	case errors.Is(err, dbutils.ErrNotFoundType):
		log.Error().Err(err).Str("userID", uid).Msg("GetSelfInfo err - User Does Not Exist")
		c.JSON(http.StatusNotFound, &response.Message{
			Message: "User does not exist"},
		)
		return
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("userID", uid).Msg("GetSelfInfo err - Internal Server Error")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	// return user
	c.JSON(http.StatusOK, dto.SuccessResponse[*model.User]{
		Data: user,
	})
}
