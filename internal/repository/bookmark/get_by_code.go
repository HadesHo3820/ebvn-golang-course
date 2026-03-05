package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
)

// GetBookmarkByCode retrieves a single bookmark record by its unique sequence code ID.
// It uses GORM's First method which returns gorm.ErrRecordNotFound if no matching record exists.
//
// Parameters:
//   - ctx: Context for the operation
//   - code: The unique auto-incremented sequence code ID
//
// Returns:
//   - *model.Bookmark: The found bookmark
//   - error: gorm.ErrRecordNotFound if not found, or other database errors
func (r *bookmarkRepo) GetBookmarkByCode(ctx context.Context, code int64) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	err := r.db.WithContext(ctx).Where("sequence_code_id = ?", code).First(&bookmark).Error
	if err != nil {
		return nil, err
	}

	return &bookmark, nil
}
