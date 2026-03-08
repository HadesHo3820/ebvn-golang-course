package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/base62"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"gorm.io/gorm"
)

// CreateBookmark inserts a new bookmark record into the database within a transaction.
//
// The transaction ensures atomicity of the following steps:
//  1. Insert the bookmark — the database auto-assigns a SequenceCodeID via the PostgreSQL sequence.
//  2. Encode the SequenceCodeID to a Base62 string.
//  3. Update the bookmark's EncodedBookmarkCode with the encoded value.
//
// If any step fails, the entire transaction is rolled back to maintain data integrity.
// Any database errors are translated into application-specific errors using dbutils.CatchDBErr.
func (r *bookmarkRepo) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Insert the bookmark — SequenceCodeID is auto-assigned by DB sequence
		if err := tx.Create(&bookmark).Error; err != nil {
			return err
		}

		// Step 2: Encode the auto-generated SequenceCodeID to Base62
		encodedCode, err := base62.Encode(bookmark.SequenceCodeID)
		if err != nil {
			return err
		}

		// Step 3: Update the bookmark with the encoded value
		bookmark.EncodedBookmarkCode = &encodedCode
		if err := tx.Model(&bookmark).Update("encoded_bookmark_code", encodedCode).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return bookmark, nil
}
