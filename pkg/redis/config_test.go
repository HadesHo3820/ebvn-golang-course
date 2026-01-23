package redis

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
		assert.Equal(t, "localhost:6379", cfg.Address)
		assert.Equal(t, "", cfg.Password)
		assert.Equal(t, 0, cfg.DB)
	})

	t.Run("environment variables", func(t *testing.T) {
		os.Setenv("REDIS_ADDR", "redis:6379")
		os.Setenv("REDIS_PASSWORD", "secret")
		os.Setenv("REDIS_DB", "1")
		defer os.Clearenv()

		cfg, err := newConfig("")
		assert.NoError(t, err)
		assert.Equal(t, "redis:6379", cfg.Address)
		assert.Equal(t, "secret", cfg.Password)
		assert.Equal(t, 1, cfg.DB)
	})

	t.Run("environment variables with prefix", func(t *testing.T) {
		os.Setenv("TEST_REDIS_ADDR", "prefixed:6379")
		defer os.Clearenv()

		cfg, err := newConfig("TEST")
		assert.NoError(t, err)
		assert.Equal(t, "prefixed:6379", cfg.Address)
	})
}
