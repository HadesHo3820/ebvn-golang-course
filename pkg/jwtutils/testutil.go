package jwtutils

import (
	"os"
	"testing"
)

// CreateInvalidPEMFile creates a temporary file with invalid PEM content for testing.
// The file is automatically cleaned up when the test ends via t.Cleanup.
//
// Parameters:
//   - t: Testing instance for cleanup registration
//
// Returns the path to the temporary file.
func CreateInvalidPEMFile(t *testing.T) string {
	t.Helper()

	invalidPEMFile, err := os.CreateTemp("", "invalid_key_*.pem")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, _ = invalidPEMFile.WriteString("this is not a valid PEM file")
	invalidPEMFile.Close()

	// Register cleanup to remove the file when test ends
	t.Cleanup(func() {
		os.Remove(invalidPEMFile.Name())
	})

	return invalidPEMFile.Name()
}
