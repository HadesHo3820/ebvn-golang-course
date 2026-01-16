// Package utils provides common utility functions used across the application.
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash of the provided plaintext password.
// It uses bcrypt.DefaultCost (currently 10) for the work factor.
// Returns the hashed password as a string and an error if hashing fails.
func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err
}

// VerifyPassword compares a plaintext password against a bcrypt hash.
// Returns true if the password matches the hash, false otherwise.
// This is a constant-time comparison to prevent timing attacks.
func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
