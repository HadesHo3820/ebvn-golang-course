package sqldb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient creates and returns a new GORM database client configured for PostgreSQL.
// It loads database connection parameters from environment variables using the specified
// prefix for namespacing (e.g., prefix "APP" looks for "APP_DB_HOST", "APP_DB_USER", etc.).
// Pass an empty string for no prefix.
//
// The function performs the following steps:
//  1. Loads configuration from environment variables via newConfig
//  2. Constructs a PostgreSQL DSN connection string
//  3. Opens a connection to the database using GORM
//
// Returns a *gorm.DB instance on success, or an error if configuration loading
// or database connection fails.
//
// Example usage:
//
//	db, err := sqldb.NewClient("")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close()
func NewClient(envPrefix string) (*gorm.DB, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return db, nil
}
