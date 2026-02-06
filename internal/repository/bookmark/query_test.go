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
		expectedTotal int64
		expectedError error
		expectAnyErr  bool // true to check for any error, not specific type
	}{
		{
			name: "success - get existing bookmarks",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputUserID:   fixture.FixtureUserOneID,
			inputLimit:    10,
			inputOffset:   0,
			expectedLen:   1, // Fixture creates 1 bookmark for UserOne
			expectedTotal: 1,
		},
		{
			name: "success - get empty for non-existent user",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputUserID:   "non-existent-uuid",
			inputLimit:    10,
			inputOffset:   0,
			expectedLen:   0,
			expectedTotal: 0,
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
				// Use CreateInBatches or Save. Create might trigger hooks, which is fine if ID is set.
				err := db.Create(&extraBookmarks).Error
				assert.NoError(t, err)
				return db
			},
			inputUserID:   fixture.FixtureUserOneID,
			inputLimit:    2,
			inputOffset:   0,
			expectedLen:   2,
			expectedTotal: 3, // 1 (fixture) + 2 (extra)
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

			bookmarks, total, err := repo.GetBookmarks(ctx, tc.inputUserID, tc.inputLimit, tc.inputOffset)

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
			assert.Equal(t, tc.expectedTotal, total)
		})
	}
}
