package fixture

import (
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"gorm.io/gorm"
)

// FixtureTimestamp is a common timestamp used across all fixture data
// to ensure consistent and deterministic test assertions.
var FixtureTimestamp = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

// UserCommonTestDB is a test fixture that provides a pre-configured database
// with User model schema and sample user data. It implements the Fixture interface
// and is commonly used across multiple User-related test cases to ensure
// consistent test data and reduce boilerplate setup code.
type UserCommonTestDB struct {
	base
}

// Migrate creates the User table schema in the test database.
// It uses GORM's AutoMigrate to create the table based on the User model definition.
// Returns an error if the migration fails.
func (f *UserCommonTestDB) Migrate() error {
	return f.db.AutoMigrate(&model.User{})
}

// GenerateData seeds the test database with sample user records.
// It creates a new database session and inserts predefined test users
// using batch insertion for efficiency. The sample data includes:
//   - Johnny Ho (ID: f47ac10b-58cc-4372-a567-0e02b2c3d479)
//   - Huy Ho (ID: 322ac10b-58cc-4372-a567-0e02b2c3d479)
//
// Returns an error if the data insertion fails.
func (f *UserCommonTestDB) GenerateData() error {
	db := f.db.Session(&gorm.Session{})

	users := []*model.User{
		{
			ID:          "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			DisplayName: "Johnny Ho",
			Username:    "johnny.ho",
			Email:       "johnny.ho@example.com",
			Password:    "$2a$$2a$10$wfpS7JvQgcHvvHLk86eFs.jhKCIucgr9fhPkyBLVQntSH0nB05106$wfpS23sf",
			CreatedAt:   FixtureTimestamp,
			UpdatedAt:   FixtureTimestamp,
		},
		{
			ID:          "322ac10b-58cc-4372-a567-0e02b2c3d479",
			DisplayName: "Huy Ho",
			Username:    "huy.ho",
			Email:       "huy.ho@example.com",
			Password:    "$2a$$2a$10$wfpS7JvQgcHvvHLk86eFs.jhKCIucgr9fhPkyBLVQntSH0nB05106$wfpS23sf",
			CreatedAt:   FixtureTimestamp,
			UpdatedAt:   FixtureTimestamp,
		},
	}

	return db.CreateInBatches(users, 10).Error
}
