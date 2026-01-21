package user

import (
	"errors"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type updateSelfInfoReqBody struct {
	DisplayName string `json:"display_name"`
	// `validate:"omitempty,email"` allow empty email (only validates format when a value is provided)
	Email       string `json:"email" validate:"omitempty,email"`
}

// UpdateSelfInfo handles user profile update requests.
// It extracts the user ID from the JWT token and updates the user's display name and/or email.
//
// @Summary Update user profile
// @Description Update the authenticated user's profile information (display name and/or email)
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body updateSelfInfoReqBody true "Update profile request"
// @Success 200 {object} response.Message "Profile updated successfully"
// @Failure 400 {object} response.Message "No data provided for update"
// @Failure 401 {object} response.Message "Invalid or missing token"
// @Failure 404 {object} response.Message "User does not exist"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/self/info [put]
func (u *userHandler) UpdateSelfInfo(c *gin.Context) {
	// Get user id from JWT token
	uid, err := utils.GetUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid token",
		})
		return
	}

	// bind request body
	reqBody, err := utils.BindInputFromRequest[updateSelfInfoReqBody](c)
	if err != nil {
		return
	}

	// call service
	err = u.svc.UpdateUser(c, uid, reqBody.DisplayName, reqBody.Email)
	switch {
	case errors.Is(err, dbutils.ErrNotFoundType):
		log.Error().Err(err).Str("userID", uid).Msg("UpdateSelfInfo err - User Does Not Exist")
		c.JSON(http.StatusNotFound, &response.Message{
			Message: "User does not exist"},
		)
		return
	case errors.Is(err, service.ErrClientNoUpdate):
		log.Error().Err(err).Str("userID", uid).Msg("UpdateSelfInfo err - No Update")
		c.JSON(http.StatusBadRequest, &response.Message{
			Message: "No data provided for update. Please provide at least one field to update."},
		)
		return
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("userID", uid).Msg("UpdateSelfInfo err - Internal Server Error")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	// return success
	c.JSON(http.StatusOK, &response.Message{
		Message: "Edit current user successfully!",
	})
}
