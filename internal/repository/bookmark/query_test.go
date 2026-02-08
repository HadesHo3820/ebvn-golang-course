package bookmark

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserID   string
		inputLimit    int
		inputOffset   int
		expectedLen   int
		expectedError error
		expectAnyErr  bool // true to check for any error, not specific type
	}{
		{
			name: "success - get existing bookmarks",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputUserID: fixture.FixtureUserOneID,
			inputLimit:  10,
			inputOffset: 0,
			expectedLen: 1, // Fixture creates 1 bookmark for UserOne
		},
		{
			name: "success - get empty for non-existent user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputUserID: "non-existent-uuid",
			inputLimit:  10,
			inputOffset: 0,
			expectedLen: 0,
		},
		{
			name: "success - pagination limit",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
				// Add 2 more bookmarks for User 1 to have total 3
				extraBookmarks := []*model.Bookmark{
					{
						Base:        model.Base{ID: "extra-1"},
						Code:        "extra1",
						URL:         "https://example.com/1",
						UserID:      fixture.FixtureUserOneID,
						Description: "Extra 1",
					},
					{
						Base:        model.Base{ID: "extra-2"},
						Code:        "extra2",
						URL:         "https://example.com/2",
						UserID:      fixture.FixtureUserOneID,
						Description: "Extra 2",
					},
				}
				err := db.Create(&extraBookmarks).Error
				assert.NoError(t, err)
				return db
			},
			inputUserID: fixture.FixtureUserOneID,
			inputLimit:  2,
			inputOffset: 0,
			expectedLen: 2,
		},
		{
			name: "error - database error (disconnected)",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
				// Close connection to simulate DB error
				sqlDB, _ := db.DB()
				sqlDB.Close()
				return db
			},
			inputUserID:  fixture.FixtureUserOneID,
			inputLimit:   10,
			inputOffset:  0,
			expectAnyErr: true, // Any database error is expected
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			bookmarks, err := repo.GetBookmarks(ctx, tc.inputUserID, tc.inputLimit, tc.inputOffset)

			if tc.expectAnyErr {
				assert.Error(t, err)
				return
			}

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, bookmarks, tc.expectedLen)
		})
	}
}

func TestBookmarkRepo_GetBookmarksCount(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserID   string
		expectedCount int64
		expectAnyErr  bool
	}{
		{
			name: "success - count existing bookmarks",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputUserID:   fixture.FixtureUserOneID,
			expectedCount: 1, // Fixture creates 1 bookmark for UserOne
		},
		{
			name: "success - count zero for non-existent user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputUserID:   "non-existent-uuid",
			expectedCount: 0,
		},
		{
			name: "success - count multiple bookmarks",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
				// Add 2 more bookmarks for User 1 to have total 3
				extraBookmarks := []*model.Bookmark{
					{
						Base:        model.Base{ID: "extra-1"},
						Code:        "extra1",
						URL:         "https://example.com/1",
						UserID:      fixture.FixtureUserOneID,
						Description: "Extra 1",
					},
					{
						Base:        model.Base{ID: "extra-2"},
						Code:        "extra2",
						URL:         "https://example.com/2",
						UserID:      fixture.FixtureUserOneID,
						Description: "Extra 2",
					},
				}
				err := db.Create(&extraBookmarks).Error
				assert.NoError(t, err)
				return db
			},
			inputUserID:   fixture.FixtureUserOneID,
			expectedCount: 3, // 1 (fixture) + 2 (extra)
		},
		{
			name: "error - database error (disconnected)",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
				sqlDB, _ := db.DB()
				sqlDB.Close()
				return db
			},
			inputUserID:  fixture.FixtureUserOneID,
			expectAnyErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			count, err := repo.GetBookmarksCount(ctx, tc.inputUserID)

			if tc.expectAnyErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCount, count)
		})
	}
}
