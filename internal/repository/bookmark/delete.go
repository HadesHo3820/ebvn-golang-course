package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
)

// DeleteBookmark soft-deletes an existing bookmark.
// It performs an ownership check to ensure only the bookmark owner can delete it.
//
// Parameters:
//   - ctx: Context for the operation
//   - bookmarkID: The ID of the bookmark to delete
//   - userID: The ID of the user attempting the deletion (for ownership validation)
//
// Returns:
//   - error: nil on success, ErrNotFoundType if bookmark doesn't exist or user doesn't own it
func (r *bookmarkRepo) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", bookmarkID, userID).
		Delete(&model.Bookmark{})

	if result.Error != nil {
		return dbutils.CatchDBErr(result.Error)
	}

	// Check if any row was actually deleted
	if result.RowsAffected == 0 {
		return dbutils.ErrNotFoundType
	}

	return nil
}
