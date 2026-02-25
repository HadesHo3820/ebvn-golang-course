package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark"
)

//go:generate mockery --name Service --filename service.go
type Service interface {
	CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, req *dto.Request) (*dto.Response[*model.Bookmark], error)
	UpdateBookmark(ctx context.Context, bookmarkID, userID, description, url string) error
	DeleteBookmark(ctx context.Context, bookmarkID, userID string) error
	// GetUrlByCode retrieves the original URL for a given bookmark code.
	// Returns the URL if found, or an error if the code does not exist.
	GetUrlByCode(ctx context.Context, code int64) (string, error)
}

type BookmarkSvc struct {
	repo bookmark.Repository
}

func NewBookmarkSvc(repo bookmark.Repository) Service {
	return &BookmarkSvc{repo: repo}
}
