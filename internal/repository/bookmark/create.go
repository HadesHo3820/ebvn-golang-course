package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
)

// CreateBookmark inserts a new bookmark record into the database.
// It returns the created bookmark model with populated fields (like ID and timestamps) or an error.
// Any database errors are translated into application-specific errors using dbutils.CatchDBErr.
func (r *bookmarkRepo) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	err := r.db.WithContext(ctx).Create(&bookmark).Error
	if err != nil {
		return nil, dbutils.CatchDBErr(err)
	}
	return bookmark, nil
}
