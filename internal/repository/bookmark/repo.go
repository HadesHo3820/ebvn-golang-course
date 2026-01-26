package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"gorm.io/gorm"
)

// Repository defines the interface for bookmark-related database operations.
// It abstracts the underlying data access logic, allowing for easier testing and maintenance.
//go:generate mockery --name Repository --filename bookmark.go
type Repository interface {
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error)
}

// bookmarkRepo is the concrete implementation of the Repository interface using GORM.
type bookmarkRepo struct {
	db *gorm.DB
}

// NewRepository creates a new instance of bookmarkRepo with the provided GORM database connection.
func NewRepository(db *gorm.DB) Repository {
	return &bookmarkRepo{db: db}
}
