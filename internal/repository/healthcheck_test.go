package repository

import (
	"testing"

	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestHealthChecker_Ping tests the Ping method of the Redis health checker.
// It verifies both successful connectivity and error handling when Redis is unavailable.
func TestHealthChecker_Ping(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(*testing.T) *redis.Client
		expectedErr error
	}{
		{
			name: "normal case",
			setupMock: func(t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedErr: nil,
		},
		{
			name: "redis connection closed",
			setupMock: func(t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			// Setup
			redisMock := tc.setupMock(t)
			healthChecker := NewRedisHealthChecker(redisMock)

			// Execute
			err := healthChecker.Ping(ctx)

			// Assert
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
