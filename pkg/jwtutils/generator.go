// Package jwtutils provides utilities for generating and validating JSON Web Tokens (JWT)
// using RSA asymmetric key pairs. It supports RS256 signing algorithm for secure
// token-based authentication.
package jwtutils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator defines the interface for generating JWT tokens.
// Implementations of this interface are responsible for creating
// signed tokens using RSA private keys.
//
//go:generate mockery --name JWTGenerator --filename jwt_generator.go
type JWTGenerator interface {
	// GenerateToken creates a new JWT token with the provided claims.
	// It returns the signed token string or an error if signing fails.
	GenerateToken(jwtContent jwt.MapClaims) (string, error)
}

// jwtGenerator is the concrete implementation of JWTGenerator.
// It holds the RSA private key used for signing tokens.
type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

// NewJWTGenerator creates a new JWTGenerator instance by loading an RSA private key
// from the specified file path. The private key file must be in PEM format.
//
// Parameters:
//   - privateKeyPath: The file path to the RSA private key in PEM format.
//
// Returns:
//   - JWTGenerator: A new generator instance ready to create signed tokens.
//   - error: An error if the file cannot be read or the key cannot be parsed.
func NewJWTGenerator(privateKeyPath string) (JWTGenerator, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{privateKey: privateKey}, nil
}

// GenerateToken creates a new JWT token signed with RS256 algorithm.
// The token contains the provided claims and is signed using the generator's
// RSA private key.
//
// Parameters:
//   - jwtContent: A map of claims to include in the token payload.
//
// Returns:
//   - string: The signed JWT token string.
//   - error: An error if the token signing process fails.
func (g *jwtGenerator) GenerateToken(jwtContent jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
