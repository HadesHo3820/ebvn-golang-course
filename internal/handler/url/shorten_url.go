package url

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// urlShortenRequest represents the JSON request body for URL shortening.
// It contains the original URL to shorten and an optional expiration time.
type urlShortenRequest struct {
	// Url is the original URL to be shortened.
	// binding:"required" makes sure this field is present
	// binding:"url" validates that the string is a valid URL format
	Url string `json:"url" binding:"required,url" example:"https://google.com"`
	// Exp is the optional expiration time in seconds for the shortened URL.
	// binding:"gte=0" ensures the expiration time is greater than or equal to 0
	Exp int `json:"exp" binding:"required,gte=0,lte=604800" example:"86400"`
}

// @Summary Shorten URL
// @Description Generate a short code for the provided URL
// @Tags URL
// @Accept json
// @Produce json
// @Param request body urlShortenRequest true "URL shorten request"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /v1/links/shorten [post]
func (h *urlHandler) ShortenUrl(c *gin.Context) {
	req, err := utils.BindInputFromRequest[urlShortenRequest](c)
	if err != nil {
		return
	}

	code, err := h.urlService.ShortenUrl(c, req.Url, req.Exp)
	if err != nil {
		// Log the error using Zerolog's structured logging:
		// - .Str("url", ...): key-value pair for context
		// - .Err(err): standard error field
		// - .Msg(...): final message and write action
		log.Error().Str("url", req.Url).Err(err).Msg("Failed to shorten URL")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[string]{
		Message: "Shorten URL generated successfully!",
		Data:    code,
	})
}
