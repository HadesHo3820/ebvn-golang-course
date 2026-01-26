package fixture

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"gorm.io/gorm"
)

const (
	FixtureBookmarkOneID = "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	FixtureBookmarkTwoID = "322ac10b-58cc-4372-a567-0e02b2c3d479"
	FixtureBookmarkOneCode = "abc12345"
	FixtureBookmarkTwoCode = "def12345"
	FixtureBookmarkURL="https://example.com/long-url"
	FixtureBookmarkDescription="My First Bookmark"
)

type BookmarkCommonTestDB struct {
	base
}

func (f *BookmarkCommonTestDB) Migrate() error {
	return f.db.AutoMigrate(&model.Bookmark{}, &model.User{})
}

// GenerateData seeds the test database.
// Currently, it just reuses User seeding because we need users to exist
// before we can create bookmarks in our tests.
func (f *BookmarkCommonTestDB) GenerateData() error {
	// This will allow us to skip the BeforeCreate hook
	db := f.db.Session(&gorm.Session{SkipHooks: true})

	users := []*model.User{
		{
			Base: model.Base{
				ID:        FixtureUserOneID,
				CreatedAt: FixtureTimestamp,
				UpdatedAt: FixtureTimestamp,
			},
			DisplayName: FixtureUserOneDisplayName,
			Username:    FixtureUserOneUsername,
			Email:       FixtureUserOneEmail,
			Password:    FixtureUserPassword,
		},
		{
			Base: model.Base{
				ID:        FixtureUserTwoID,
				CreatedAt: FixtureTimestamp,
				UpdatedAt: FixtureTimestamp,
			},
			DisplayName: FixtureUserTwoDisplayName,
			Username:    FixtureUserTwoUsername,
			Email:       FixtureUserTwoEmail,
			Password:    FixtureUserPassword,
		},
	}

	err := db.CreateInBatches(users, 100).Error
	if err != nil {
		return err
	}

	bookmarks := []*model.Bookmark{
		{
			Base: model.Base{
				ID:        FixtureBookmarkOneID,
				CreatedAt: FixtureTimestamp,
				UpdatedAt: FixtureTimestamp,
			},
			URL:         FixtureBookmarkURL,
			Code:        FixtureBookmarkOneCode,
			Description: FixtureBookmarkDescription,
			UserID:      FixtureUserOneID,
			User: users[0],
		},
		{
			Base: model.Base{
				ID:        FixtureBookmarkTwoID,
				CreatedAt: FixtureTimestamp,
				UpdatedAt: FixtureTimestamp,
			},
			URL:         FixtureBookmarkURL,
			Code:        FixtureBookmarkTwoCode,
			Description: FixtureBookmarkDescription,
			UserID:      FixtureUserTwoID,
			User: users[1],
		},
	}
	

	return db.CreateInBatches(bookmarks, 100).Error
}
