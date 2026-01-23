package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTValidator defines the interface for validating JWT tokens.
// Implementations of this interface are responsible for verifying
// token signatures using RSA public keys and extracting claims.
//go:generate mockery --name JWTValidator --filename jwt_validator.go
type JWTValidator interface {
	// ValidateToken verifies the token signature and returns the claims if valid.
	// It returns an error if the token is invalid or verification fails.
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

// jwtValidator is the concrete implementation of JWTValidator.
// It holds the RSA public key used for verifying token signatures.
type jwtValidator struct {
	publicKey *rsa.PublicKey
}

// NewJWTValidator creates a new JWTValidator instance by loading an RSA public key
// from the specified file path. The public key file must be in PEM format.
//
// Parameters:
//   - publicKeyPath: The file path to the RSA public key in PEM format.
//
// Returns:
//   - JWTValidator: A new validator instance ready to verify tokens.
//   - error: An error if the file cannot be read or the key cannot be parsed.
func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return &jwtValidator{publicKey: publicKey}, nil
}

// errInvalidToken is returned when a token fails validation.
// This can occur due to an invalid signature, expiration, or malformed token.
var errInvalidToken = errors.New("invalid token")

// ValidateToken verifies the JWT token signature using the validator's RSA public key
// and extracts the claims if the token is valid.
//
// Parameters:
//   - tokenString: The JWT token string to validate.
//
// Returns:
//   - jwt.MapClaims: The claims extracted from the token if validation succeeds.
//   - error: errInvalidToken if the token is invalid, expired, or has an invalid signature.
func (v *jwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return v.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	return token.Claims.(jwt.MapClaims), nil
}
