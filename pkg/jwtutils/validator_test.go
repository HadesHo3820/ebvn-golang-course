package jwtutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestNewJWTValidator tests the NewJWTValidator constructor function.
func TestNewJWTValidator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		keyPath   string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid public key path",
			keyPath:   filepath.FromSlash("./public.test.pem"),
			expectErr: false,
		},
		{
			name:      "invalid public key path - file not found",
			keyPath:   filepath.FromSlash("./non-existent.pem"),
			expectErr: true,
			errMsg:    "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.keyPath)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, validator)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, validator)
		})
	}
}

// TestNewJWTValidator_InvalidPEM tests NewJWTValidator with invalid PEM content.
func TestNewJWTValidator_InvalidPEM(t *testing.T) {
	t.Parallel()

	// Create a temporary file with invalid PEM content
	tempFile, err := os.CreateTemp("", "invalid_key_*.pem")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("this is not a valid PEM file")
	assert.NoError(t, err)
	tempFile.Close()

	validator, err := NewJWTValidator(tempFile.Name())
	assert.Error(t, err)
	assert.Nil(t, validator)
}

// TestJWTValidator_ValidateToken tests the ValidateToken method.
func TestJWTValidator_ValidateToken(t *testing.T) {
	t.Parallel()

	// Generate a valid token using the generator
	generator, err := NewJWTGenerator(filepath.FromSlash("./private.test.pem"))
	assert.NoError(t, err)

	validClaims := jwt.MapClaims{
		"id":   "1234",
		"name": "John",
	}
	validToken, err := generator.GenerateToken(validClaims)
	assert.NoError(t, err)

	// Create validator with matching public key
	validator, err := NewJWTValidator(filepath.FromSlash("./public.test.pem"))
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		tokenString    string
		expectErr      bool
		expectedClaims jwt.MapClaims
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			expectErr:   false,
			expectedClaims: jwt.MapClaims{
				"id":   "1234",
				"name": "John",
			},
		},
		{
			name:        "empty token",
			tokenString: "",
			expectErr:   true,
		},
		{
			name:        "malformed token",
			tokenString: "not.a.valid.jwt.token",
			expectErr:   true,
		},
		{
			name:        "invalid signature - token signed with different key",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := validator.ValidateToken(tc.tokenString)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, errInvalidToken, err)
				assert.Nil(t, claims)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, claims)
			assert.Equal(t, tc.expectedClaims["id"], claims["id"])
			assert.Equal(t, tc.expectedClaims["name"], claims["name"])
		})
	}
}
