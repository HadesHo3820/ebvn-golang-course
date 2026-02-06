package bookmark

import (
	"testing"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputBookmark *model.Bookmark
		expectedErr   error
		verifyFunc    func(t *testing.T, db *gorm.DB, bookmark *model.Bookmark)
	}{
		{
			name: "success - create valid bookmark",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmark: &model.Bookmark{
				UserID:      fixture.FixtureUserOneID,
				URL:         "https://example.com/unique-url",
				Code:        "uniq123",
				Description: "My Unique Bookmark",
			},
			verifyFunc: func(t *testing.T, db *gorm.DB, expected *model.Bookmark) {
				var actual model.Bookmark
				// Verify persistence:
				// 1. Query the DB for the bookmark using its unique code.
				// 2. Preload("User") fetches the associated User to ensure the foreign key relationship is valid.
				err := db.Preload("User").Where("code = ?", expected.Code).First(&actual).Error
				assert.NoError(t, err)

				assert.NotEmpty(t, actual.ID)
				assert.Equal(t, expected.URL, actual.URL)
				assert.Equal(t, expected.Code, actual.Code)
				assert.Equal(t, expected.Description, actual.Description)
				assert.Equal(t, expected.UserID, actual.UserID)

				// Verify relationship preload
				assert.Equal(t, fixture.FixtureUserOneID, actual.User.ID)
				assert.Equal(t, fixture.FixtureUserOneUsername, actual.User.Username)
			},
		},
		{
			name: "error - duplicate code",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
				// This fixture already seeds FixtureBookmarkOneCode ("abc12345")
				return db
			},
			inputBookmark: &model.Bookmark{
				UserID:      fixture.FixtureUserOneID,
				URL:         "https://example.com/duplicate",
				Code:        fixture.FixtureBookmarkOneCode, // Reusing existing code "abc12345"
				Description: "Duplicate Code Bookmark",
			},
			expectedErr: dbutils.ErrDuplicationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			created, err := repo.CreateBookmark(ctx, tc.inputBookmark)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, created)

			// Verify ID generation (BeforeCreate hook)
			assert.NotEmpty(t, created.ID)
			assert.WithinDuration(t, time.Now(), created.CreatedAt, 2*time.Second)

			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, created)
			}
		})
	}
}
