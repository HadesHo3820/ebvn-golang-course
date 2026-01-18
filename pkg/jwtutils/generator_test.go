package jwtutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestNewJWTGenerator tests the NewJWTGenerator constructor function.
func TestNewJWTGenerator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		keyPath   string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid private key path",
			keyPath:   filepath.FromSlash("./private.test.pem"),
			expectErr: false,
		},
		{
			name:      "invalid private key path - file not found",
			keyPath:   filepath.FromSlash("./non-existent.pem"),
			expectErr: true,
			errMsg:    "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generator, err := NewJWTGenerator(tc.keyPath)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, generator)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, generator)
		})
	}
}

// TestNewJWTGenerator_InvalidPEM tests NewJWTGenerator with invalid PEM content.
func TestNewJWTGenerator_InvalidPEM(t *testing.T) {
	t.Parallel()

	// Create a temporary file with invalid PEM content
	tempFile, err := os.CreateTemp("", "invalid_key_*.pem")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("this is not a valid PEM file")
	assert.NoError(t, err)
	tempFile.Close()

	generator, err := NewJWTGenerator(tempFile.Name())
	assert.Error(t, err)
	assert.Nil(t, generator)
}

// TestJWTGenerator_GenerateToken_Integration tests that generated tokens can be validated.
func TestJWTGenerator_GenerateToken_Integration(t *testing.T) {
	t.Parallel()

	// Create generator and validator with matching key pair.
	// filepath.FromSlash converts forward slashes to the OS-specific path separator,
	// ensuring cross-platform compatibility (e.g., "/" on Unix, "\" on Windows).
	generator, err := NewJWTGenerator(filepath.FromSlash("./private.test.pem"))
	assert.NoError(t, err)

	validator, err := NewJWTValidator(filepath.FromSlash("./public.test.pem"))
	assert.NoError(t, err)

	// Generate a token
	inputClaims := jwt.MapClaims{
		"id":   "1234",
		"name": "John",
	}
	token, err := generator.GenerateToken(inputClaims)
	assert.NoError(t, err)

	// Validate the token
	outputClaims, err := validator.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, outputClaims)

	// Verify the claims match
	assert.Equal(t, inputClaims["id"], outputClaims["id"])
	assert.Equal(t, inputClaims["name"], outputClaims["name"])
}
