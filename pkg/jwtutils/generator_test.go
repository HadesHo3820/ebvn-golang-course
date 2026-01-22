package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestNewJWTGenerator tests the NewJWTGenerator constructor function.
// It uses table-driven tests to cover:
//   - Valid private key path
//   - File not found error
//   - Invalid PEM content
func TestNewJWTGenerator(t *testing.T) {
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
			name:      "success - valid private key path",
			keyPath:   filepath.FromSlash("./private.test.pem"),
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

// TestJWTGenerator_GenerateToken tests the GenerateToken method.
// It uses table-driven tests to cover various claim scenarios.
func TestJWTGenerator_GenerateToken(t *testing.T) {
	t.Parallel()

	// Setup generator with valid key
	generator, err := NewJWTGenerator(filepath.FromSlash("./private.test.pem"))
	assert.NoError(t, err)
	assert.NotNil(t, generator)

	// Setup validator for token verification
	validator, err := NewJWTValidator(filepath.FromSlash("./public.test.pem"))
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		claims      jwt.MapClaims
		expectErr   bool
		verifyClaim func(t *testing.T, outputClaims jwt.MapClaims)
	}{
		{
			name: "success - basic claims",
			claims: jwt.MapClaims{
				"id":   "1234",
				"name": "John",
			},
			expectErr: false,
			verifyClaim: func(t *testing.T, outputClaims jwt.MapClaims) {
				assert.Equal(t, "1234", outputClaims["id"])
				assert.Equal(t, "John", outputClaims["name"])
			},
		},
		{
			name: "success - with sub claim",
			claims: jwt.MapClaims{
				"sub":   "user-uuid-123",
				"email": "test@example.com",
			},
			expectErr: false,
			verifyClaim: func(t *testing.T, outputClaims jwt.MapClaims) {
				assert.Equal(t, "user-uuid-123", outputClaims["sub"])
				assert.Equal(t, "test@example.com", outputClaims["email"])
			},
		},
		{
			name: "success - with numeric claims",
			claims: jwt.MapClaims{
				"user_id": float64(42),
				"exp":     float64(9999999999),
			},
			expectErr: false,
			verifyClaim: func(t *testing.T, outputClaims jwt.MapClaims) {
				assert.Equal(t, float64(42), outputClaims["user_id"])
				assert.Equal(t, float64(9999999999), outputClaims["exp"])
			},
		},
		{
			name:      "success - empty claims",
			claims:    jwt.MapClaims{},
			expectErr: false,
			verifyClaim: func(t *testing.T, outputClaims jwt.MapClaims) {
				assert.NotNil(t, outputClaims)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Generate token
			token, err := generator.GenerateToken(tc.claims)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Validate the generated token
			outputClaims, err := validator.ValidateToken(token)
			assert.NoError(t, err)
			assert.NotNil(t, outputClaims)

			// Verify claims if callback provided
			if tc.verifyClaim != nil {
				tc.verifyClaim(t, outputClaims)
			}
		})
	}
}
