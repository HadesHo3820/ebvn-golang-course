package response

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// TestInputFieldError tests the InputFieldError function.
// It uses table-driven tests to cover:
//   - Generic error (non-validation error)
//   - validator.ValidationErrors with single field error
//   - validator.ValidationErrors with multiple field errors
func TestInputFieldError(t *testing.T) {
	t.Parallel()

	// Create a validator for generating real validation errors
	validate := validator.New()

	// Test struct for validation
	type testInput struct {
		Username string `validate:"required"`
		Email    string `validate:"required,email"`
		Age      int    `validate:"gte=18"`
	}

	testCases := []struct {
		name            string
		setupError      func() error
		expectedMessage string
		expectedDetails any
	}{
		{
			name: "generic error - returns InputErrResponse",
			setupError: func() error {
				return errors.New("some random error")
			},
			expectedMessage: InputErrMessage,
			expectedDetails: nil,
		},
		{
			name: "validation error - single field",
			setupError: func() error {
				input := testInput{Username: "", Email: "valid@example.com", Age: 20}
				return validate.Struct(input)
			},
			expectedMessage: InputErrMessage,
			expectedDetails: []string{"Username is invalid (required)"},
		},
		{
			name: "validation error - multiple fields",
			setupError: func() error {
				input := testInput{Username: "", Email: "invalid-email", Age: 10}
				return validate.Struct(input)
			},
			expectedMessage: InputErrMessage,
			expectedDetails: []string{
				"Username is invalid (required)",
				"Email is invalid (email)",
				"Age is invalid (gte)",
			},
		},
		{
			name: "validation error - email format",
			setupError: func() error {
				input := testInput{Username: "user", Email: "not-an-email", Age: 20}
				return validate.Struct(input)
			},
			expectedMessage: InputErrMessage,
			expectedDetails: []string{"Email is invalid (email)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.setupError()
			result := InputFieldError(err)

			assert.Equal(t, tc.expectedMessage, result.Message)
			assert.Equal(t, tc.expectedDetails, result.Details)
		})
	}
}

// TestMessage tests the Message struct initialization and field values.
func TestMessage(t *testing.T) {
	t.Parallel()

	t.Run("message with details", func(t *testing.T) {
		msg := Message{
			Message: "Test message",
			Details: []string{"detail1", "detail2"},
		}
		assert.Equal(t, "Test message", msg.Message)
		assert.Equal(t, []string{"detail1", "detail2"}, msg.Details)
	})

	t.Run("message without details", func(t *testing.T) {
		msg := Message{
			Message: "Test message",
		}
		assert.Equal(t, "Test message", msg.Message)
		assert.Nil(t, msg.Details)
	})
}

// TestPredefinedResponses tests the predefined response variables.
func TestPredefinedResponses(t *testing.T) {
	t.Parallel()

	t.Run("InternalErrResponse", func(t *testing.T) {
		assert.Equal(t, InternalErrMessage, InternalErrResponse.Message)
		assert.Nil(t, InternalErrResponse.Details)
	})

	t.Run("InputErrResponse", func(t *testing.T) {
		assert.Equal(t, InputErrMessage, InputErrResponse.Message)
		assert.Nil(t, InputErrResponse.Details)
	})
}
