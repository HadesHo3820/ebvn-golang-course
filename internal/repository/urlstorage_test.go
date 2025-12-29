// Package repository contains unit tests for the URL storage functionality.
// These tests use miniredis to simulate a Redis server in-memory,
// allowing for fast and isolated tests without requiring a real Redis instance.
package repository

import (
	"context"
	"testing"

	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestUrlStorage_StoreURL validates the StoreUrl method of the UrlStorage interface.
// It uses table-driven tests to cover various scenarios for storing URLs in Redis.
//
// Test Structure:
//   - name: A descriptive name for the test case, shown in test output.
//   - setupMock: A function that initializes and returns an in-memory Redis client.
//   - expectedErr: The expected error from StoreUrl (nil for success cases).
//   - verifyFunc: A callback to verify the stored data after a successful operation.
//
// The test runs in parallel (t.Parallel()) for improved performance.
// Each sub-test also runs in parallel to maximize test throughput.
func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	// testCases defines the table of test scenarios.
	// Each case sets up its own isolated Redis mock to prevent test interference.
	testCases := []struct {
		name        string                                     // Test case name for identification
		setupMock   func() *redis.Client                       // Factory function for mock Redis client
		expectedErr error                                      // Expected error result from StoreUrl
		verifyFunc  func(ctx context.Context, r *redis.Client) // Verification callback for success cases
	}{
		{
			// "normal case" tests the happy path where a URL is successfully stored.
			name: "normal case",
			setupMock: func() *redis.Client {
				// InitMockRedis creates an in-memory Redis instance managed by the test lifecycle.
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				// Verify the URL was stored correctly by fetching it directly from Redis.
				url, err := r.Get(ctx, "test").Result()
				assert.Nil(t, err)
				assert.Equal(t, url, "https://google.com")
			},
		},
	}

	// Iterate through each test case and execute as a sub-test.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run sub-tests in parallel for performance
			ctx := t.Context()

			// Setup: Create the mock Redis client and repository instance.
			redisMock := tc.setupMock()
			urlRepo := NewUrlStorage(redisMock)

			// Execute: Call StoreUrl with test data.
			err := urlRepo.StoreUrl(ctx, "test", "https://google.com")

			// Assert: Check that the error matches expectations.
			assert.Equal(t, tc.expectedErr, err)

			// Verify: If no error occurred, run the verification function
			// to confirm the data was stored correctly.
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}

// TestUrlStorage_Exists validates the Exists method of the UrlStorage interface.
// It tests whether the method correctly identifies existing and non-existing codes.
func TestUrlStorage_Exists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string               // Test case name
		setupMock      func() *redis.Client // Factory function for mock Redis client
		code           string               // Code to check for existence
		expectedExists bool                 // Expected result from Exists
		expectedErr    error                // Expected error (nil for success)
	}{
		{
			name: "code exists",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				// Pre-populate the mock Redis with a code
				mock.Set(context.Background(), "existing_code", "https://example.com", 0)
				return mock
			},
			code:           "existing_code",
			expectedExists: true,
			expectedErr:    nil,
		},
		{
			name: "code does not exist",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			code:           "non_existing_code",
			expectedExists: false,
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Setup
			redisMock := tc.setupMock()
			urlRepo := NewUrlStorage(redisMock)

			// Execute
			exists, err := urlRepo.Exists(ctx, tc.code)

			// Assert
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedExists, exists)
		})
	}
}

// TestUrlStorage_StoreUrlIfNotExists validates the atomic SETNX-based storage.
// It tests both successful storage (new code) and collision detection (code exists).
func TestUrlStorage_StoreUrlIfNotExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string               // Test case name
		setupMock      func() *redis.Client // Factory function for mock Redis client
		code           string               // Code to store
		url            string               // URL to store
		expectedStored bool                 // Expected result (true if stored, false if collision)
		expectedErr    error                // Expected error
	}{
		{
			name: "successful storage - new code",
			setupMock: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			code:           "newcode",
			url:            "https://example.com",
			expectedStored: true,
			expectedErr:    nil,
		},
		{
			name: "collision - code already exists",
			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				// Pre-populate the code
				mock.Set(context.Background(), "existing", "https://old-url.com", 0)
				return mock
			},
			code:           "existing",
			url:            "https://new-url.com",
			expectedStored: false,
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Setup
			redisMock := tc.setupMock()
			urlRepo := NewUrlStorage(redisMock)

			// Execute
			stored, err := urlRepo.StoreUrlIfNotExists(ctx, tc.code, tc.url, 0)

			// Assert
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedStored, stored)

			// Verify: if stored, check the value
			if stored {
				val, _ := redisMock.Get(ctx, tc.code).Result()
				assert.Equal(t, tc.url, val)
			}
		})
	}
}
