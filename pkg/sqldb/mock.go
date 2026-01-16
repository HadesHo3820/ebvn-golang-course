// Package sqldb provides database utilities for the application.
// This file contains mock database utilities for testing purposes.
package sqldb

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMockDB creates and returns an in-memory SQLite database for testing.
// It uses a unique UUID-based connection string with shared cache mode,
// ensuring each test gets an isolated database instance.
//
// The database is configured with silent logging to reduce test output noise.
// If the database connection fails, the test is immediately marked as failed
// using t.Fatalf.
//
// Parameters:
//   - t: The testing.T instance used to report failures.
//
// Returns:
//   - *gorm.DB: A configured GORM database instance ready for testing.
func InitMockDB(t *testing.T) *gorm.DB {
	// Build SQLite connection string with the following URI format:
	// - file:<name>  : SQLite URI prefix for file-based databases
	// - uuid         : Unique database name to isolate each test instance
	// - mode=memory  : Store the database in RAM instead of disk. This means the DB
	//                  only exists in memory - all data is lost when the program exits.
	//                  This makes it ideal for testing since each test run starts fresh.
	// - cache=shared : Allow multiple connections to share the same in-memory database
	//                  (required for GORM connection pooling to work with in-memory DBs)
	cxn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String())
	db, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to create db: %v", err)
	}

	return db
}
