// Package sqldb provides utilities for SQL database migrations using golang-migrate.
// This package extends the base sqldb package with migration capabilities,
// allowing version-controlled database schema changes through SQL migration files.
package sqldb

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/gorm"

	// Blank import registers the file source driver for golang-migrate.
	// This allows migration files to be loaded from the filesystem using
	// the "file://" URL scheme in the migration path.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateSQLDB applies database migrations from the specified path with flexible control
// over migration direction and granularity. It uses golang-migrate to execute versioned
// SQL migration files against the PostgreSQL database connected through the provided GORM instance.
//
// The migration process:
//  1. Extracts the underlying *sql.DB from GORM
//  2. Creates a postgres driver instance for golang-migrate
//  3. Initializes a migrate instance with the migration files
//  4. Executes migrations based on the specified mode
//
// Migration files should be named following the pattern:
//
//	{version}_{description}.up.sql   (for applying changes)
//	{version}_{description}.down.sql (for reverting changes)
//
// Example:
//
//	000001_init_db.up.sql
//	000001_init_db.down.sql
//
// Parameters:
//   - db: GORM database instance to migrate
//   - migrationPath: Path to migration files directory, using URL format.
//     Example: "file://./migrations" for a local "migrations" directory
//   - mode: Migration mode, either "up" or "steps"
//   - "up": Apply all pending migrations
//   - "steps": Apply a specific number of migrations (controlled by steps parameter)
//   - steps: Number of migrations to apply (only used when mode="steps")
//   - Positive value: migrate forward (apply migrations)
//   - Negative value: migrate backward (revert migrations)
//   - Zero: returns an error
//
// Returns an error if:
//   - Unable to extract *sql.DB from GORM
//   - Postgres driver initialization fails
//   - Migration instance creation fails
//   - Invalid mode is specified
//   - steps is 0 when mode="steps"
//   - Any migration execution fails (excluding ErrNoChange)
func MigrateSQLDB(db *gorm.DB, migrationPath string, mode string, steps int) error {
	// Extract the underlying *sql.DB from GORM for use with golang-migrate
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Create postgres-specific driver instance for golang-migrate
	postgresDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	// Initialize migrate instance with migration files and database connection
	m, err := migrate.NewWithDatabaseInstance(migrationPath, db.Name(), postgresDriver)
	if err != nil {
		return err
	}

	return migrateSchema(m, mode, steps)
}

// migrateSchema executes database migrations based on the specified mode and steps.
// It handles two migration modes:
//   - "up": Applies all pending migrations in order
//   - "steps": Applies a specific number of migrations (forward or backward)
//
// The function gracefully handles the ErrNoChange error, treating it as a successful
// no-op when no migrations are pending. All other errors are wrapped with a descriptive
// prefix for easier debugging.
//
// Parameters:
//   - m: Migrate instance configured with database connection and migration files
//   - mode: Migration execution mode ("up" or "steps")
//   - steps: Number of migrations to apply (only used when mode="steps")
//   - Positive: migrate forward (e.g., 2 applies next 2 up migrations)
//   - Negative: migrate backward (e.g., -1 reverts last migration)
//   - Zero: returns an error
//
// Returns an error if:
//   - mode is not "up" or "steps"
//   - steps is 0 when mode="steps"
//   - Migration execution fails (excluding ErrNoChange)
func migrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var migrationErr error

	switch mode {
	case "up":
		migrationErr = m.Up()
	case "steps":
		if steps == 0 {
			return errors.New("[Database migration] Steps must not be 0. Please use a positive number to migrate up, a negative number to migrate down.")
		}
		migrationErr = m.Steps(steps)
	default:
		return errors.New("[Database migration] Invalid mode. Please use 'up' or 'steps'.")
	}

	// ErrNoChange is excluded because it's not actually an error condition 
	// it's an informational status that means "the database is already at the desired state"
	if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		return fmt.Errorf("[Database migration] Error: %s", migrationErr.Error())
	}

	return nil

}
