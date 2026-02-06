package bookmark

import (
	"context"
)

// UpdateBookmark implements the business logic for updating an existing bookmark.
// It delegates the update operation to the repository layer after any necessary
// business logic processing.
//
// Parameters:
//   - ctx: Context for the operation
//   - bookmarkID: The ID of the bookmark to update
//   - userID: The ID of the user requesting the update (for ownership validation)
//   - description: The new description for the bookmark
//   - url: The new URL for the bookmark
//
// Returns:
//   - error: nil on success, or an error from the repository layer
func (s *BookmarkSvc) UpdateBookmark(ctx context.Context, bookmarkID, userID, description, url string) error {
	return s.repo.UpdateBookmark(ctx, bookmarkID, userID, description, url)
}
