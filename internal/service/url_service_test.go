// Package service contains unit tests for the URL shortening service.
// These tests use mocks to isolate the service logic from the repository layer.
package service

import (
	"errors"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestShortenUrl_ShortenUrl validates the ShortenUrl method of the ShortenUrl service.
// It uses table-driven tests to cover various scenarios including success,
// collision handling, and error cases.
func TestShortenUrl_ShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func() *mocks.UrlStorage
		url         string
		exp         int    // expiration time in seconds
		expectCode  bool   // true if we expect a valid code to be returned
		expectedErr string // expected error message (empty for no error)
	}{
		{
			name: "successful storage on first attempt",
			setupMock: func() *mocks.UrlStorage {
				mockRepo := new(mocks.UrlStorage)
				// First call succeeds (stored = true)
				mockRepo.On("StoreUrlIfNotExists", mock.Anything, mock.AnythingOfType("string"), "https://example.com", 0).
					Return(true, nil).Once()
				return mockRepo
			},
			url:         "https://example.com",
			exp:         0,
			expectCode:  true,
			expectedErr: "",
		},
		{
			name: "successful storage after one collision",
			setupMock: func() *mocks.UrlStorage {
				mockRepo := new(mocks.UrlStorage)
				// First call: collision (stored = false)
				mockRepo.On("StoreUrlIfNotExists", mock.Anything, mock.AnythingOfType("string"), "https://example.com", 3600).
					Return(false, nil).Once()
				// Second call: success (stored = true)
				mockRepo.On("StoreUrlIfNotExists", mock.Anything, mock.AnythingOfType("string"), "https://example.com", 3600).
					Return(true, nil).Once()
				return mockRepo
			},
			url:         "https://example.com",
			exp:         3600,
			expectCode:  true,
			expectedErr: "",
		},
		{
			name: "max retries exceeded - all collisions",
			setupMock: func() *mocks.UrlStorage {
				mockRepo := new(mocks.UrlStorage)
				// All 5 attempts result in collision
				mockRepo.On("StoreUrlIfNotExists", mock.Anything, mock.AnythingOfType("string"), "https://example.com", 0).
					Return(false, nil).Times(5)
				return mockRepo
			},
			url:         "https://example.com",
			exp:         0,
			expectCode:  false,
			expectedErr: "failed to generate unique code after 5 attempts",
		},
		{
			name: "repository error",
			setupMock: func() *mocks.UrlStorage {
				mockRepo := new(mocks.UrlStorage)
				mockRepo.On("StoreUrlIfNotExists", mock.Anything, mock.AnythingOfType("string"), "https://example.com", 0).
					Return(false, errors.New("redis connection failed")).Once()
				return mockRepo
			},
			url:         "https://example.com",
			exp:         0,
			expectCode:  false,
			expectedErr: "redis connection failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Setup
			mockRepo := tc.setupMock()
			service := NewShortenUrl(mockRepo)

			// Execute
			code, err := service.ShortenUrl(ctx, tc.url, tc.exp)

			// Assert
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				if tc.expectCode {
					assert.NotEmpty(t, code)
					assert.Len(t, code, 7) // urlCodeLength
				}
			}

			// Verify mock expectations: ensures all methods set up with .On()
			// were called the expected number of times with the expected arguments.
			// This catches cases where:
			// - A method was never called when it should have been
			// - A method was called more/fewer times than expected (.Once(), .Times(n))
			// - A method was called with unexpected arguments
			mockRepo.AssertExpectations(t)
		})
	}
}
