package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	t.Run("success - hash password", func(t *testing.T) {
		t.Parallel()
		password := "mysecretpassword"
		hash, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Parallel()

	password := "mysecretpassword"
	hash, _ := HashPassword(password)

	t.Run("success - correct password", func(t *testing.T) {
		t.Parallel()
		valid := VerifyPassword(password, hash)
		assert.True(t, valid)
	})

	t.Run("success - incorrect password", func(t *testing.T) {
		t.Parallel()
		valid := VerifyPassword("wrongpassword", hash)
		assert.False(t, valid)
	})

	t.Run("success - invalid hash", func(t *testing.T) {
		t.Parallel()
		valid := VerifyPassword(password, "invalidhash")
		assert.False(t, valid)
	})
}
