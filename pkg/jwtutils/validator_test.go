package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestNewJWTValidator tests the NewJWTValidator constructor function.
// It uses table-driven tests to cover:
//   - Valid public key path
//   - File not found error
//   - Invalid PEM content
func TestNewJWTValidator(t *testing.T) {
	t.Parallel()

	// Create invalid PEM file using helper
	invalidPEMPath := CreateInvalidPEMFile(t)

	testCases := []struct {
		name      string
		keyPath   string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "success - valid public key path",
			keyPath:   filepath.FromSlash("./public.test.pem"),
			expectErr: false,
		},
		{
			name:      "error - file not found",
			keyPath:   filepath.FromSlash("./non-existent.pem"),
			expectErr: true,
			errMsg:    "no such file or directory",
		},
		{
			name:      "error - invalid PEM content",
			keyPath:   invalidPEMPath,
			expectErr: true,
			errMsg:    "invalid",
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

// TestJWTValidator_ValidateToken tests the ValidateToken method.
// It uses table-driven tests to cover various token validation scenarios.
func TestJWTValidator_ValidateToken(t *testing.T) {
	t.Parallel()

	// Setup: Generate a valid token using the generator
	generator, err := NewJWTGenerator(filepath.FromSlash("./private.test.pem"))
	assert.NoError(t, err)

	validClaims := jwt.MapClaims{
		"id":   "1234",
		"name": "John",
	}
	validToken, err := generator.GenerateToken(validClaims)
	assert.NoError(t, err)

	// Setup: Create validator with matching public key
	validator, err := NewJWTValidator(filepath.FromSlash("./public.test.pem"))
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		tokenString  string
		expectErr    bool
		verifyClaims func(t *testing.T, claims jwt.MapClaims)
	}{
		{
			name:        "success - valid token",
			tokenString: validToken,
			expectErr:   false,
			verifyClaims: func(t *testing.T, claims jwt.MapClaims) {
				assert.Equal(t, "1234", claims["id"])
				assert.Equal(t, "John", claims["name"])
			},
		},
		{
			name:        "error - empty token",
			tokenString: "",
			expectErr:   true,
		},
		{
			name:        "error - malformed token (wrong segments)",
			tokenString: "not.a.valid.jwt.token",
			expectErr:   true,
		},
		{
			name:        "error - invalid signature (signed with different key)",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
			expectErr:   true,
		},
		{
			name:        "error - truncated token",
			tokenString: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpX",
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			claims, err := validator.ValidateToken(tc.tokenString)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, errInvalidToken, err)
				assert.Nil(t, claims)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, claims)

			// Verify claims if callback provided
			if tc.verifyClaims != nil {
				tc.verifyClaims(t, claims)
			}
		})
	}
}
