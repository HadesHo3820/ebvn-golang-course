// Package redis provides utilities for creating and managing Redis client connections.
// It handles configuration loading from environment variables and client initialization.
package redis

import "github.com/redis/go-redis/v9"

// NewClientWithDB creates a new Redis client with an explicitly specified database number.
// This allows creating multiple clients connected to different logical databases within
// the same Redis instance. Redis supports databases 0-15 by default.
//
// Parameters:
//   - envPrefix: Prefix for environment variable names (empty string for defaults)
//   - db: The Redis database number to connect to (0-15)
//
// Returns:
//   - *redis.Client: A configured Redis client connected to the specified database.
//   - error: An error if configuration loading fails.
//
// Example:
//
//	cacheClient, _ := NewClientWithDB("", 0)  // DB 0 for cache
//	generalClient, _ := NewClientWithDB("", 1) // DB 1 for general purposes
func NewClientWithDB(envPrefix string, db int) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}
	// Override the DB from config with the explicitly provided db parameter
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       db,
	})

	return redisClient, nil
}
