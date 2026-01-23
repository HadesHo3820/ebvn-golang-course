package url

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

// UrlHandler represents the HTTP handler for URL related requests.
type UrlHandler interface {
	// ShortenUrl handles the request to shorten a URL.
	ShortenUrl(c *gin.Context)
	// GetUrl handles the request to retrieve the original URL from a short code.
	GetUrl(c *gin.Context)
}

// urlHandler implements the UrlHandler interface.
type urlHandler struct {
	urlService service.ShortenUrl
}

// NewUrlHandler creates a new instance of UrlHandler with the given service.
func NewUrlHandler(urlService service.ShortenUrl) UrlHandler {
	return &urlHandler{urlService: urlService}
}
