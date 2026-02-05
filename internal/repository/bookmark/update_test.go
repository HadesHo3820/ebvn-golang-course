package bookmark

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupDB          func(t *testing.T) *gorm.DB
		inputBookmarkID  string
		inputUserID      string
		inputDescription string
		inputURL         string
		expectedErr      error
		verifyFunc       func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "success - update existing bookmark",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmarkID:  fixture.FixtureBookmarkOneID,
			inputUserID:      fixture.FixtureUserOneID,
			inputDescription: "Updated Description",
			inputURL:         "https://updated-example.com",
			expectedErr:      nil,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var bookmark struct {
					Description string
					URL         string
				}
				err := db.Table("bookmarks").
					Where("id = ?", fixture.FixtureBookmarkOneID).
					First(&bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, "Updated Description", bookmark.Description)
				assert.Equal(t, "https://updated-example.com", bookmark.URL)
			},
		},
		{
			name: "error - bookmark not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmarkID:  "non-existent-id",
			inputUserID:      fixture.FixtureUserOneID,
			inputDescription: "Description",
			inputURL:         "https://example.com",
			expectedErr:      dbutils.ErrNotFoundType,
		},
		{
			name: "error - bookmark belongs to different user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmarkID:  fixture.FixtureBookmarkOneID, // Belongs to User One
			inputUserID:      fixture.FixtureUserTwoID,     // Trying with User Two
			inputDescription: "Updated Description",
			inputURL:         "https://updated-example.com",
			expectedErr:      dbutils.ErrNotFoundType, // Should return not found for security
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			err := repo.UpdateBookmark(ctx, tc.inputBookmarkID, tc.inputUserID, tc.inputDescription, tc.inputURL)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db)
			}
		})
	}
}
