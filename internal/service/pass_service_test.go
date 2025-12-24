// Package service provides unit tests for the service layer.
// These tests validate the business logic implementations following
// Go testing best practices including table-driven tests and parallel execution.
package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// urlSafeRegex is a precompiled regular expression that validates URL-safe characters.
// It matches strings containing only alphanumeric characters (a-z, A-Z, 0-9).
// This is used to verify that generated passwords are safe for use in URLs
// without requiring encoding.
var urlSafeRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// TestPasswordService_GeneratePassword tests the GeneratePassword method of the Password service.
// It uses table-driven tests to validate:
//   - Password length matches the expected length (10 characters)
//   - No errors are returned during generation
//   - Generated password contains only URL-safe alphanumeric characters
//
// The test runs in parallel for improved performance.
func TestPasswordService_GeneratePassword(t *testing.T) {
	// call t.Parallel() to run the test in parallel
	t.Parallel()

	testCases := []struct {
		name        string
		expectedLen int
		expectErr   error
	}{
		{
			name:        "normal case",
			expectedLen: 10,
			expectErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// call t.Parallel() to indicate every test case can run in parallel
			t.Parallel()

			testSvc := NewPassword()
			pass, err := testSvc.GeneratePassword()

			assert.Equal(t, tc.expectedLen, len(pass))
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(pass), true)
		})
	}
}
