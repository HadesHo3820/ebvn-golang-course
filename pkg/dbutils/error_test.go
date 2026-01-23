package dbutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestCatchDBErr tests the CatchDBErr function.
// It uses table-driven tests to cover all error translation scenarios.
func TestCatchDBErr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		inputError  error
		expectedErr error
	}{
		{
			name:        "unique constraint violation - lowercase",
			inputError:  errors.New("unique constraint violation on column 'email'"),
			expectedErr: ErrDuplicationType,
		},
		{
			name:        "unique constraint violation - uppercase",
			inputError:  errors.New("UNIQUE CONSTRAINT violation"),
			expectedErr: ErrDuplicationType,
		},
		{
			name:        "unique constraint violation - mixed case",
			inputError:  errors.New("Unique Constraint failed for 'username'"),
			expectedErr: ErrDuplicationType,
		},
		{
			name:        "record not found - gorm error",
			inputError:  gorm.ErrRecordNotFound,
			expectedErr: ErrNotFoundType,
		},
		{
			name:        "generic error - no filter matches",
			inputError:  errors.New("some random database error"),
			expectedErr: errors.New("some random database error"),
		},
		{
			name:        "connection error - no filter matches",
			inputError:  errors.New("connection refused"),
			expectedErr: errors.New("connection refused"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := CatchDBErr(tc.inputError)

			// For sentinel errors, use errors.Is
			if errors.Is(tc.expectedErr, ErrDuplicationType) || errors.Is(tc.expectedErr, ErrNotFoundType) {
				assert.ErrorIs(t, result, tc.expectedErr)
			} else {
				// For non-matched errors, compare the error message
				assert.Equal(t, tc.expectedErr.Error(), result.Error())
			}
		})
	}
}

// TestFilterDuplicationType tests the filterDuplicationType function directly.
func TestFilterDuplicationType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		inputError error
		matches    bool
	}{
		{
			name:       "matches - lowercase",
			inputError: errors.New("unique constraint violation"),
			matches:    true,
		},
		{
			name:       "matches - uppercase",
			inputError: errors.New("UNIQUE CONSTRAINT error"),
			matches:    true,
		},
		{
			name:       "matches - partial",
			inputError: errors.New("ERROR: duplicate key violates unique constraint \"users_email_key\""),
			matches:    true,
		},
		{
			name:       "no match - different error",
			inputError: errors.New("foreign key violation"),
			matches:    false,
		},
		{
			name:       "no match - record not found",
			inputError: gorm.ErrRecordNotFound,
			matches:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matched, err := filterDuplicationType(tc.inputError)

			assert.Equal(t, tc.matches, matched)
			if matched {
				assert.Equal(t, ErrDuplicationType, err)
			}
		})
	}
}

// TestFilterRecordNotFound tests the filterRecordNotFound function directly.
func TestFilterRecordNotFound(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		inputError error
		matches    bool
	}{
		{
			name:       "matches - gorm.ErrRecordNotFound",
			inputError: gorm.ErrRecordNotFound,
			matches:    true,
		},
		{
			name:       "no match - different error",
			inputError: errors.New("some other error"),
			matches:    false,
		},
		{
			name:       "no match - similar message but not gorm error",
			inputError: errors.New("record not found"),
			matches:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matched, err := filterRecordNotFound(tc.inputError)

			assert.Equal(t, tc.matches, matched)
			if matched {
				assert.Equal(t, ErrNotFoundType, err)
			}
		})
	}
}

// TestSentinelErrors tests that sentinel errors are correctly defined.
func TestSentinelErrors(t *testing.T) {
	t.Parallel()

	t.Run("ErrDuplicationType", func(t *testing.T) {
		assert.Error(t, ErrDuplicationType)
		assert.Equal(t, "duplication type", ErrDuplicationType.Error())
	})

	t.Run("ErrNotFoundType", func(t *testing.T) {
		assert.Error(t, ErrNotFoundType)
		assert.Equal(t, "not found type", ErrNotFoundType.Error())
	})
}
