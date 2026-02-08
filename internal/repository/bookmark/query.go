package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
)

// GetBookmarksCount returns the total number of bookmarks associated with a specific user.
// It is used for pagination to determine the total number of pages available.
func (r *bookmarkRepo) GetBookmarksCount(ctx context.Context, userID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetBookmarks retrieves a paginated list of bookmarks for a specific user.
// It returns a slice of bookmarks and any error encountered.
// The results are ordered by creation date in descending order.
func (r *bookmarkRepo) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error) {
	// initializes a slice in Go with a specific type, length, and capacity.
	// Since you know you are querying for a maximum of limit records
	// Setting the capacity upfront prevents Go from having to resize and reallocate the underlying array multiple times as GORM appends the results
	// It is more efficient than make([]*model.Bookmark, 0).
	bookmarks := make([]*model.Bookmark, 0, limit)

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}
