package bookmark

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetBookmarks returns a paginated list of bookmarks.
// @Summary      List bookmarks
// @Description  Get a paginated list of bookmarks for the authenticated user
// @Tags         Bookmark
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Page number (default 1)"
// @Param        limit  query     int  false  "Items per page (default 10)"
// @Success      200    {object}  dto.SuccessResponse[[]model.Bookmark]
// @Failure      401    {object}  response.Message "Unauthorized"
// @Failure      500    {object}  response.Message "Internal server error"
// @Router       /v1/bookmarks [get]
func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	input, uid, err := utils.BindInputFromRequestWithAuth[dto.Request](c)
	if err != nil {
		return
	}

	res, err := h.svc.GetBookmarks(c, uid, input)
	if err != nil {
		log.Error().Err(err).Str("uid", uid).Msg("Failed to list bookmarks")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[[]*model.Bookmark]{
		Data:     res.Data,
		Metadata: &res.Metadata,
	})
}
