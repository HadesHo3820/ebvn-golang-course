// Package sqldb provides utilities for connecting to PostgreSQL databases
// using GORM. It handles configuration loading from environment variables
// and establishes database connections with sensible defaults.
package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// config holds the PostgreSQL database connection parameters.
// All fields can be configured via environment variables with the specified
// envconfig tags, or will fall back to their default values.
//
// Environment variables (with optional prefix):
//   - DB_HOST:     Database host address (default: "localhost")
//   - DB_USER:     Database user name (default: "admin")
//   - DB_PASSWORD: Database user password (default: "admin")
//   - DB_NAME:     Database name (default: "bookmark_db")
//   - DB_PORT:     Database port (default: "5432")
type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	User     string `default:"admin" envconfig:"DB_USER"`
	Password string `default:"admin" envconfig:"DB_PASSWORD"`
	Name     string `default:"bookmark_db" envconfig:"DB_NAME"`
	Port     string `default:"5432" envconfig:"DB_PORT"`
}

// newConfig creates a new config instance by loading values from environment
// variables. The envPrefix parameter allows namespacing environment variables
// (e.g., prefix "APP" would look for "APP_DB_HOST" instead of "DB_HOST").
// Pass an empty string for no prefix.
//
// Returns an error if environment variable processing fails.
func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetDSN returns a PostgreSQL Data Source Name (DSN) connection string
// formatted for use with GORM's postgres driver. The connection uses
// sslmode=disable for local development convenience.
//
// Example output:
//
//	"host=localhost user=admin password=admin dbname=bookmark_db port=5432 sslmode=disable"
func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
}
