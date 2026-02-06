package bookmark

import (
	"errors"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type deleteBookmarkInput struct {
	// ID is the bookmark identifier from the URL path
	ID string `uri:"id" validate:"required,uuid"`
}

// DeleteBookmark deletes an existing bookmark for the authenticated user.
//
// @Summary      Delete a bookmark
// @Description  Delete an existing bookmark. Only the bookmark owner can delete it.
// @Tags         Bookmark
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string           true  "Bookmark ID (UUID)"
// @Success      200  {object}  response.Message "Success"
// @Failure      400  {object}  response.Message "Invalid input"
// @Failure      401  {object}  response.Message "Unauthorized"
// @Failure      404  {object}  response.Message "Bookmark not found"
// @Failure      500  {object}  response.Message "Internal server error"
// @Router       /v1/bookmarks/{id} [delete]
func (h *bookmarkHandler) DeleteBookmark(c *gin.Context) {
	// Get user id from JWT token
	uid, err := utils.GetUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid token",
		})
		return
	}

	// Getting input from request and validate
	input, err := utils.BindInputFromRequest[deleteBookmarkInput](c)
	if err != nil {
		return
	}

	err = h.svc.DeleteBookmark(c, input.ID, uid)
	if err != nil {
		if errors.Is(err, dbutils.ErrNotFoundType) {
			c.JSON(http.StatusNotFound, &response.Message{
				Message: "Bookmark not found",
			})
			return
		}

		log.Error().Err(err).Str("uid", uid).Str("bookmark_id", input.ID).Msg("Failed to delete bookmark")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, &response.Message{
		Message: "Success",
	})
}
