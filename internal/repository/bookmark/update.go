package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
)

// UpdateBookmark updates an existing bookmark's description and URL.
// It performs an ownership check to ensure only the bookmark owner can update it.
//
// Parameters:
//   - ctx: Context for the operation
//   - bookmarkID: The ID of the bookmark to update
//   - userID: The ID of the user attempting the update (for ownership validation)
//   - description: The new description for the bookmark
//   - url: The new URL for the bookmark
//
// Returns:
//   - error: nil on success, ErrNotFoundType if bookmark doesn't exist or user doesn't own it
func (r *bookmarkRepo) UpdateBookmark(ctx context.Context, bookmarkID, userID, description, url string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("id = ? AND user_id = ?", bookmarkID, userID).
		Updates(map[string]any{
			"description": description,
			"url":         url,
		})

	if result.Error != nil {
		return dbutils.CatchDBErr(result.Error)
	}

	// Check if any row was actually updated
	if result.RowsAffected == 0 {
		return dbutils.ErrNotFoundType
	}

	return nil
}
