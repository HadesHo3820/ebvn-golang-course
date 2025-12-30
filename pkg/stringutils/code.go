package stringutils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	// charset contains the alphanumeric characters used for password generation.
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// KeyGenerator defines the interface for generating random codes.
//go:generate mockery --name KeyGenerator --filename key_generator.go
type KeyGenerator interface {
	GenerateCode(length int) (string, error)
}

type keyGen struct{}

func NewKeyGenerator() KeyGenerator {
	return &keyGen{}
}

func (k *keyGen) GenerateCode(length int) (string, error) {
	return GenerateCode(length)
}

func GenerateCode(length int) (string, error) {
	var strBuilder bytes.Buffer

	// generate random password of length passLength
	for range length {
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
