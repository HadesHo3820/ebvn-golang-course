package bookmark

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
}

type bookmarkHandler struct {
	svc bookmark.Service
}

func NewHandler(svc bookmark.Service) Handler {
	return &bookmarkHandler{svc: svc}
}
