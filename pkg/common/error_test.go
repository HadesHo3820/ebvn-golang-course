package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHandleError validates the HandleError utility function.
// It verifies:
//   - nil errors do not cause a panic
//   - non-nil errors cause a panic with the error message
func TestHandleError(t *testing.T) {
	t.Parallel()

	t.Run("nil error - no panic", func(t *testing.T) {
		t.Parallel()
		// Should not panic
		assert.NotPanics(t, func() {
			HandleError(nil)
		})
	})

	t.Run("non-nil error - should panic", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		assert.PanicsWithError(t, "test error", func() {
			HandleError(testErr)
		})
	})
}
