package bookmark

import (
	"context"
)

// DeleteBookmark implements the business logic for deleting an existing bookmark.
// It delegates the delete operation to the repository layer.
//
// Parameters:
//   - ctx: Context for the operation
//   - bookmarkID: The ID of the bookmark to delete
//   - userID: The ID of the user requesting the deletion (for ownership validation)
//
// Returns:
//   - error: nil on success, or an error from the repository layer
func (s *BookmarkSvc) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	return s.repo.DeleteBookmark(ctx, bookmarkID, userID)
}
