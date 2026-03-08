package redis_test

import (
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func TestNewClientWithDB(t *testing.T) {
	testCases := []struct {
		name         string
		envPrefix    string
		db           int
		setupEnv     func(*testing.T)
		expectedAddr string
		expectedDB   int
		expectError  bool
	}{
		{
			name:         "success - default config",
			envPrefix:    "",
			db:           0,
			setupEnv:     func(t *testing.T) {},
			expectedAddr: "localhost:6379",
			expectedDB:   0,
			expectError:  false,
		},
		{
			name:         "success - explicit db overrides default",
			envPrefix:    "",
			db:           5,
			setupEnv:     func(t *testing.T) {},
			expectedAddr: "localhost:6379",
			expectedDB:   5,
			expectError:  false,
		},
		{
			name:      "success - env vars config",
			envPrefix: "",
			db:        2,
			setupEnv: func(t *testing.T) {
				t.Setenv("REDIS_ADDR", "127.0.0.1:6380")
				t.Setenv("REDIS_PASSWORD", "secret")
				t.Setenv("REDIS_DB", "10") // Should be ignored/overridden by parameter
			},
			expectedAddr: "127.0.0.1:6380",
			expectedDB:   2,
			expectError:  false,
		},
		{
			name:      "success - with prefix",
			envPrefix: "APP",
			db:        1,
			setupEnv: func(t *testing.T) {
				t.Setenv("APP_REDIS_ADDR", "redis-host:6379")
			},
			expectedAddr: "redis-host:6379",
			expectedDB:   1,
			expectError:  false,
		},
		{
			name:      "error - invalid config",
			envPrefix: "",
			db:        0,
			setupEnv: func(t *testing.T) {
				t.Setenv("REDIS_DB", "invalid-int") // Causes envconfig to fail
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Setup environment variables for this test case
			tc.setupEnv(t)

			client, err := redis.NewClientWithDB(tc.envPrefix, tc.db)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)

				// Verify client options
				opts := client.Options()
				assert.Equal(t, tc.expectedAddr, opts.Addr)
				assert.Equal(t, tc.expectedDB, opts.DB)

				// Cleanup
				client.Close()
			}
		})
	}
}

// Ensure NewClientWithDB sets reasonable defaults for timeouts, etc.
// (Checking standard go-redis defaults implicitly via successful creation)
