package service

import "github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"

const (
	// passLength defines the default length of the generated password.
	passLength = 10
)

// passwordService implements the Password interface.
type passwordService struct{}

// Password defines the interface for password-related operations.
//
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
	return stringutils.GenerateCode(passLength)
}
