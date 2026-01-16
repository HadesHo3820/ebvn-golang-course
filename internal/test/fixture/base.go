// Package fixture provides test fixtures for setting up and managing test databases.
// It defines a common interface for test database operations and utilities
// to create isolated test environments with pre-populated data.
package fixture

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/pkg/sqldb"
	"gorm.io/gorm"
)

// Fixture defines the interface for test database fixtures.
// Implementations of this interface provide the necessary setup logic
// for specific test scenarios, including database initialization,
// schema migration, and test data generation.
type Fixture interface {
	// SetupDB initializes the fixture with the provided database connection.
	// This method should store the database reference for later use.
	SetupDB(db *gorm.DB)

	// Migrate applies the necessary database schema migrations.
	// It should create all required tables and indexes for the test.
	// Returns an error if migration fails.
	Migrate() error

	// GenerateData populates the database with test data.
	// This method is called after Migrate() to seed the database
	// with the fixtures needed for testing.
	// Returns an error if data generation fails.
	GenerateData() error

	// DB returns the underlying GORM database connection.
	// This allows test code to perform additional database operations.
	DB() *gorm.DB
}

// NewFixture creates and initializes a test database using the provided fixture.
// It performs the following steps:
//  1. Creates a mock database connection for testing
//  2. Runs schema migrations via the fixture's Migrate method
//  3. Seeds the database with test data via GenerateData
//
// If any step fails, the test is immediately terminated with t.Fatalf.
// Returns the initialized GORM database connection for use in tests.
func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	// create test database
	fix.SetupDB(sqldb.InitMockDB(t))

	// migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatalf("failed to migrate db for testing: %v", err)
	}

	// create test data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("failed to generate test data for testing: %v", err)
	}

	return fix.DB()
}
