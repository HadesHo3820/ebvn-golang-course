package sqldb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Don't run parallel because we modify environment variables

	t.Run("default values", func(t *testing.T) {
		os.Clearenv()
		cfg, err := newConfig("")
		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.Host)
		assert.Equal(t, "admin", cfg.User)
		assert.Equal(t, "admin", cfg.Password)
		assert.Equal(t, "bookmark_db", cfg.Name)
		assert.Equal(t, "5432", cfg.Port)
	})

	t.Run("environment variables", func(t *testing.T) {
		os.Setenv("DB_HOST", "db-host")
		os.Setenv("DB_USER", "db-user")
		os.Setenv("DB_PASSWORD", "db-pass")
		os.Setenv("DB_NAME", "db-name")
		os.Setenv("DB_PORT", "5000")
		defer os.Clearenv()

		cfg, err := newConfig("")
		assert.NoError(t, err)
		assert.Equal(t, "db-host", cfg.Host)
		assert.Equal(t, "db-user", cfg.User)
		assert.Equal(t, "db-pass", cfg.Password)
		assert.Equal(t, "db-name", cfg.Name)
		assert.Equal(t, "5000", cfg.Port)
	})
}

func TestConfig_GetDSN(t *testing.T) {
	t.Parallel()

	cfg := &config{
		Host:     "host",
		User:     "user",
		Password: "password",
		Name:     "dbname",
		Port:     "5432",
	}

	expected := "host=host user=user password=password dbname=dbname port=5432 sslmode=disable"
	assert.Equal(t, expected, cfg.GetDSN())
}
