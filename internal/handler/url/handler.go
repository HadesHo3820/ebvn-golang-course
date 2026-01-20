package url

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/gin-gonic/gin"
)

type UrlHandler interface {
	ShortenUrl(c *gin.Context)
	GetUrl(c *gin.Context)
}

type urlHandler struct {
	urlService service.ShortenUrl
}

func NewUrlHandler(urlService service.ShortenUrl) UrlHandler {
	return &urlHandler{urlService: urlService}
}
