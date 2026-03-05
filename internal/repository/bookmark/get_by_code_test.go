package bookmark

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_GetBookmarkByCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupDB          func(t *testing.T) *gorm.DB
		inputCode        int64
		expectedBookmark *model.Bookmark
		expectNotFound   bool
		expectAnyErr     bool
	}{
		{
			name: "success - found by code",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputCode: fixture.FixtureBookmarkOneSequenceCodeID,
			expectedBookmark: &model.Bookmark{
				URL:    fixture.FixtureBookmarkURL,
				UserID: fixture.FixtureUserOneID,
			},
		},
		{
			name: "not found - non-existent code",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputCode:      99999,
			expectNotFound: true,
		},
		{
			name: "error - database error (disconnected)",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
				sqlDB, _ := db.DB()
				sqlDB.Close()
				return db
			},
			inputCode:    fixture.FixtureBookmarkOneSequenceCodeID,
			expectAnyErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			bookmark, err := repo.GetBookmarkByCode(ctx, tc.inputCode)

			if tc.expectNotFound {
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				assert.Nil(t, bookmark)
				return
			}

			if tc.expectAnyErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, bookmark)
			assert.Equal(t, tc.expectedBookmark.URL, bookmark.URL)
			assert.Equal(t, tc.expectedBookmark.UserID, bookmark.UserID)
			assert.Equal(t, tc.inputCode, bookmark.SequenceCodeID)
		})
	}
}
