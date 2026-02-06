package bookmark

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark"
	"github.com/gin-gonic/gin"
)

// Handler defines the interface for bookmark HTTP handlers.
type Handler interface {
	// CreateBookmark handles the creation of a new bookmark.
	CreateBookmark(c *gin.Context)
	// GetBookmarks retrieves a list of bookmarks.
	GetBookmarks(c *gin.Context)
	// UpdateBookmark handles updating an existing bookmark.
	UpdateBookmark(c *gin.Context)
	// DeleteBookmark handles the deletion of a bookmark.
	DeleteBookmark(c *gin.Context)
}

type bookmarkHandler struct {
	svc bookmark.Service
}

// NewHandler creates a new instance of the bookmark handler.
func NewHandler(svc bookmark.Service) Handler {
	return &bookmarkHandler{svc: svc}
}
