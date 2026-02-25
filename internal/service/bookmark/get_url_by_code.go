package bookmark

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// ErrBookmarkNotFound is a sentinel error returned when a bookmark code
// does not exist in the database. Callers should use errors.Is()
// to check for this specific error condition.
var ErrBookmarkNotFound = errors.New("bookmark not found")

// GetUrlByCode retrieves the original URL for a given bookmark code.
// It queries the repository and translates gorm.ErrRecordNotFound to ErrBookmarkNotFound
// for a cleaner abstraction that doesn't leak repository implementation details.
//
// Returns:
//   - The original URL if the bookmark exists.
//   - ErrBookmarkNotFound if the code does not exist.
//   - Other errors for repository/connection failures.
func (s *BookmarkSvc) GetUrlByCode(ctx context.Context, code int64) (string, error) {
	bookmark, err := s.repo.GetBookmarkByCode(ctx, code)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrBookmarkNotFound
	}
	if err != nil {
		return "", err
	}

	return bookmark.URL, nil
}
