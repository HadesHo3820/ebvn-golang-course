package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
)

const codeLength = 9

// CreateBookmark implements the business logic for creating a new bookmark.
// It generates a unique short code for the URL using the configured KeyGenerator
// and persists the bookmark data to the repository.
//
// Parameters:
//   - ctx: Context for the operation
//   - description: User-provided description
//   - url: The target URL to shorten
//   - userID: The ID of the owner
//
// Returns:
//   - *model.Bookmark: The created bookmark with generated ID and code
//   - error: Any error during generation or persistence
func (s *BookmarkSvc) CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	// create code
	code, err := s.codeGen.GenerateCode(codeLength)
	if err != nil {
		return nil, err
	}

	// create bookmark
	bookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		Code:        code,
		UserID:      userID,
	}

	bookmarkModel, err := s.repo.CreateBookmark(ctx, bookmark)
	if err != nil {
		return nil, err
	}

	return bookmarkModel, nil
}
