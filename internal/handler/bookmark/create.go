package bookmark

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	_ "github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type createBookmarkInput struct {
	// Description of the bookmark
	Description string `json:"description" example:"Your description here" validate:"lte=255,gte=1"`
	// URL to be shortened
	URL string `json:"url" example:"https://example.com" validate:"required,url,lte=2048"`
}

// CreateBookmark creates a new bookmark for the authenticated user.
//
// @Summary      Create a new bookmark
// @Description  Create a new bookmark with a description and target URL. Returns the created bookmark with its short code.
// @Tags         Bookmark
// @Accept       json
// @Produce      json
// @Security 	 BearerAuth
// @Param        request        body      createBookmarkInput  true  "Bookmark details"
// @Success      200            {object}  model.Bookmark
// @Failure      400            {object}  response.Message     "Invalid input"
// @Failure      401            {object}  response.Message     "Unauthorized"
// @Failure      500            {object}  response.Message     "Internal server error"
// @Router       /v1/bookmarks [post]
func (h *bookmarkHandler) CreateBookmark(c *gin.Context) {
	// Get user id from JWT token
	uid, err := utils.GetUIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid token",
		})
		return
	}

	// Getting input from request and validate
	input, err := utils.BindInputFromRequest[createBookmarkInput](c)
	if err != nil {
		return
	}

	res, err := h.svc.CreateBookmark(c, input.Description, input.URL, uid)
	if err != nil {
		log.Error().Err(err).Str("uid", uid).Msg("Failed to create bookmark")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
