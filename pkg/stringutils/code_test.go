package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	t.Run("success - generate valid code", func(t *testing.T) {
		t.Parallel()
		length := 6
		code, err := GenerateCode(length)
		assert.NoError(t, err)
		assert.Len(t, code, length)
		for _, char := range code {
			assert.Contains(t, charset, string(char))
		}
	})

	t.Run("success - generate zero length code", func(t *testing.T) {
		t.Parallel()
		code, err := GenerateCode(0)
		assert.NoError(t, err)
		assert.Empty(t, code)
	})
}

func TestKeyGenerator_GenerateCode(t *testing.T) {
	t.Parallel()

	kg := NewKeyGenerator()
	assert.NotNil(t, kg)

	code, err := kg.GenerateCode(10)
	assert.NoError(t, err)
	assert.Len(t, code, 10)
}
