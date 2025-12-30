// Package handler provides HTTP handlers for the bookmark service API.
// This file contains the URL shortening handler which generates short codes for URLs.
package handler

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
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

// urlShortenResponse represents the JSON response for a successful URL shortening.
type urlShortenResponse struct {
	// Message indicates the status of the operation.
	Message string `json:"message" example:"Shorten URL generated successfully!"`
	// Code is the generated short code that maps to the original URL.
	Code string `json:"code" example:"string"`
}

// urlShortenHandler is the concrete implementation of the UrlShorten interface.
// It holds a reference to the URL shortening service.
type urlShortenHandler struct {
	urlService service.ShortenUrl
}

// UrlShorten defines the interface for URL shortening HTTP handlers.
// This interface allows for dependency injection and easier testing.
type UrlShorten interface {
	// ShortenUrl handles HTTP requests to shorten a URL.
	ShortenUrl(c *gin.Context)
}

// NewUrlShorten creates a new UrlShorten handler with the provided service dependency.
// It follows the constructor injection pattern for dependency management.
//
// Parameters:
//   - svc: The URL shortening service that provides the business logic.
//
// Returns:
//   - UrlShorten: An interface implementation that handles URL shortening HTTP requests.
func NewUrlShorten(svc service.ShortenUrl) UrlShorten {
	return &urlShortenHandler{
		urlService: svc,
	}
}

// @Summary Shorten URL
// @Description Generate a short code for the provided URL
// @Tags URL
// @Accept json
// @Produce json
// @Param request body urlShortenRequest true "URL shorten request"
// @Success 200 {object} urlShortenResponse
// @Failure 400 {object} map[string]string "Bad Request - invalid URL or validation error"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /v1/links/shorten [post]
func (h *urlShortenHandler) ShortenUrl(c *gin.Context) {
	var req urlShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong input"})
		return
	}

	code, err := h.urlService.ShortenUrl(c, req.Url, req.Exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, urlShortenResponse{
		Message: "Shorten URL generated successfully!",
		Code:    code,
	})
}
