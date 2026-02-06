package bookmark

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_DeleteBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		setupDB         func(t *testing.T) *gorm.DB
		inputBookmarkID string
		inputUserID     string
		expectedErr     error
		expectAnyErr    bool // Set to true to check for any error (not specific type)
		verifyFunc      func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "success - delete existing bookmark",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmarkID: fixture.FixtureBookmarkOneID,
			inputUserID:     fixture.FixtureUserOneID,
			expectedErr:     nil,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				// Verify bookmark was soft-deleted
				var count int64
				db.Table("bookmarks").
					Where("id = ? AND deleted_at IS NULL", fixture.FixtureBookmarkOneID).
					Count(&count)
				assert.Equal(t, int64(0), count, "Bookmark should be soft-deleted")
			},
		},
		{
			name: "error - bookmark not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmarkID: "00000000-0000-0000-0000-000000000000",
			inputUserID:     fixture.FixtureUserOneID,
			expectedErr:     dbutils.ErrNotFoundType,
		},
		{
			name: "error - bookmark belongs to different user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmarkID: fixture.FixtureBookmarkOneID, // Belongs to User One
			inputUserID:     fixture.FixtureUserTwoID,     // Trying with User Two
			expectedErr:     dbutils.ErrNotFoundType,      // Should return not found for security
		},
		{
			name: "error - database error (disconnected)",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
				// Get the underlying SQL DB and close it to simulate connection error
				sqlDB, _ := db.DB()
				sqlDB.Close()
				return db
			},
			inputBookmarkID: fixture.FixtureBookmarkOneID,
			inputUserID:     fixture.FixtureUserOneID,
			expectAnyErr:    true, // Any database error is expected
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			err := repo.DeleteBookmark(ctx, tc.inputBookmarkID, tc.inputUserID)

			if tc.expectAnyErr {
				assert.Error(t, err)
				return
			}

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
