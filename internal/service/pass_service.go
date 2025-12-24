package service

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	// charset contains the alphanumeric characters used for password generation.
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// passLength defines the default length of the generated password.
	passLength = 10
)

// passwordService implements the Password interface.
type passwordService struct{}

// Password defines the interface for password-related operations.
//go:generate mockery --name Password --filename pass_service.go
type Password interface {
	// GeneratePassword creates a new random password.
	GeneratePassword() (string, error)
}

// NewPassword creates a new instance of the password service.
func NewPassword() Password {
	return &passwordService{}
}

// GeneratePassword generates a cryptographically secure random password.
// It uses characters from the predefined charset and has a fixed length of 10.
// Returns the generated password string or an error if the random number generator fails.
func (s *passwordService) GeneratePassword() (string, error) {
	var strBuilder bytes.Buffer

	// generate random password of length passLength
	for range passLength {
		// Generate a random index using crypto/rand for cryptographic security.
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		// Select the character at the random index and append it to the result.
		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}
	return strBuilder.String(), nil
}
