package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
)

// GetBookmarks retrieves a paginated list of bookmarks for a specific user.
// It returns the slice of bookmarks, the total count of records matching the criteria, and any error encountered.
//
// The pagination is implemented using a two-step approach:
// 1. Count the total number of records matching the user ID.
// 2. If records exist, retrieve the specific page of data using limit and offset.
func (r *bookmarkRepo) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error) {
	bookmarks := make([]*model.Bookmark, 0)
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return bookmarks, 0, nil
	}

	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&bookmarks).Error
	if err != nil {
		return nil, 0, err
	}
	return bookmarks, total, nil
}
