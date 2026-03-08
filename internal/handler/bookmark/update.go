package bookmark

import (
	"errors"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type updateBookmarkInput struct {
	// ID is the bookmark identifier from the URL path
	ID string `uri:"id" validate:"required,uuid"`
	// Description of the bookmark
	Description string `json:"description" example:"Google" validate:"lte=255"`
	// URL to be shortened
	URL string `json:"url" example:"https://www.google.com" validate:"required,url,lte=2048"`
}

// UpdateBookmark updates an existing bookmark for the authenticated user.
//
// @Summary      Update a bookmark
// @Description  Update an existing bookmark's description and URL. Only the bookmark owner can update it.
// @Tags         Bookmark
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string               true  "Bookmark ID (UUID)"
// @Param        request  body      updateBookmarkInput  true  "Updated bookmark details"
// @Success      200      {object}  dto.SuccessResponse[any]  "Success"
// @Failure      400      {object}  response.Message     "Invalid input"
// @Failure      401      {object}  response.Message     "Unauthorized"
// @Failure      404      {object}  response.Message     "Bookmark not found"
// @Failure      500      {object}  response.Message     "Internal server error"
// @Router       /v1/bookmarks/{id} [put]
func (h *bookmarkHandler) UpdateBookmark(c *gin.Context) {
	// Getting input from request and validate
	input, uid, err := utils.BindInputFromRequestWithAuth[updateBookmarkInput](c)
	if err != nil {
		return
	}

	err = h.svc.UpdateBookmark(c, input.ID, uid, input.Description, input.URL)
	if err != nil {
		if errors.Is(err, dbutils.ErrNotFoundType) {
			c.JSON(http.StatusNotFound, &response.Message{
				Message: "Bookmark not found",
			})
			return
		}

		log.Error().Err(err).Str("uid", uid).Str("bookmark_id", input.ID).Msg("Failed to update bookmark")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[any]{
		Message: "Bookmark updated successfully",
	})
}
