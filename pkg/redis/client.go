// Package redis provides utilities for creating and managing Redis client connections.
// It handles configuration loading from environment variables and client initialization.
package redis

import "github.com/redis/go-redis/v9"

// NewClient creates a new Redis client using configuration loaded from environment variables.
// The envPrefix parameter allows namespacing environment variables (e.g., "CACHE_" would look
// for CACHE_REDIS_ADDRESS instead of REDIS_ADDRESS). Pass an empty string for default variable names.
//
// Returns:
//   - *redis.Client: A configured Redis client ready for use.
//   - error: An error if configuration loading fails.
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return redisClient, nil
}
