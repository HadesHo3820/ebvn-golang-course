package bookmark

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/pagination"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// listBookmarksResponse is a helper struct for Swagger documentation
type listBookmarksResponse struct {
	Data     []*model.Bookmark    `json:"data"`
	Metadata pagination.Metadata `json:"metadata"`
}

// GetBookmarks returns a paginated list of bookmarks.
// @Summary      List bookmarks
// @Description  Get a paginated list of bookmarks for the authenticated user
// @Tags         Bookmark
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Page number (default 1)"
// @Param        limit  query     int  false  "Items per page (default 10)"
// @Success      200    {object}  listBookmarksResponse
// @Failure      401    {object}  response.Message "Unauthorized"
// @Failure      500    {object}  response.Message "Internal server error"
// @Router       /v1/bookmarks [get]
func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	// Get user id from JWT token
	uid, err := utils.GetUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid token",
		})
		return
	}

	input, err := utils.BindInputFromRequest[pagination.Request](c)
	if err != nil {
		return
	}

	res, err := h.svc.GetBookmarks(c, uid, input)
	if err != nil {
		log.Error().Err(err).Str("uid", uid).Msg("Failed to list bookmarks")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, listBookmarksResponse{
		Data:     res.Data,
		Metadata: res.Metadata,
	})
}
