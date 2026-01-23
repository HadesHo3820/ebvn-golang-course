package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	// Cannot run parallel as it modifies global state (env vars and zerolog global level)

	t.Run("default level", func(t *testing.T) {
		os.Clearenv()
		SetLogLevel()
		assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
	})

	t.Run("valid level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "debug")
		defer os.Clearenv()
		SetLogLevel()
		assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
	})

	t.Run("invalid level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "invalid")
		defer os.Clearenv()
		SetLogLevel()
		assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
	})

	t.Run("empty level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "")
		defer os.Clearenv()
		SetLogLevel()
		assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
	})
}
