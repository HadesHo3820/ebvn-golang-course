package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/pagination"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
)

//go:generate mockery --name Service --filename service.go
type Service interface {
	CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, req *pagination.Request) (*pagination.Response[*model.Bookmark], error)
}

type BookmarkSvc struct {
	repo    bookmark.Repository
	codeGen stringutils.KeyGenerator
}

func NewBookmarkSvc(repo bookmark.Repository, codeGen stringutils.KeyGenerator) Service {
	return &BookmarkSvc{repo: repo, codeGen: codeGen}
}
