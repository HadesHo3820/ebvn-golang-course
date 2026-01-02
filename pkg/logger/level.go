package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// SetLogLevel sets the global zerolog level based on the "LOG_LEVEL" environment variable.
// It defaults to zerolog.InfoLevel if the environment variable is not set or invalid.
func SetLogLevel() {
	level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
