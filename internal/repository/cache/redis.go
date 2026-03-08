package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache implements the DB interface using Redis as the cache backend.
// It uses Redis Hash data structures (HSET/HGET) for storing cached data,
// which allows grouping related cache entries under a common key while
// maintaining individual field-level access.
type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis-backed cache instance.
// It accepts a configured Redis client and returns a DB interface implementation.
func NewRedisCache(redis *redis.Client) DB {
	return &redisCache{client: redis}
}

// SetCacheData stores a value in the cache using Redis Hash (HSET).
// The cacheGroupKey acts as the Redis key, and cacheKey is the hash field.
// This hierarchical structure allows grouping related cache entries together.
// After setting the value, it applies an expiration time to the entire hash group.
//
// Note: The expiration is set on the cacheGroupKey, meaning all entries in the
// hash group will expire together. Each call resets the expiration timer.
func (db *redisCache) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value any, exp time.Duration) error {
	err := db.client.HSet(ctx, cacheGroupKey, cacheKey, value).Err()
	if err != nil {
		return err
	}

	return db.client.Expire(ctx, cacheGroupKey, exp).Err()
}

// GetCacheData retrieves a cached value from Redis Hash (HGET).
// It returns the value as []byte, which allows the caller to deserialize
// the data into their specific type. Using []byte provides:
//   - Flexibility: Caller decides how to unmarshal (JSON, Protobuf, etc.)
//   - Efficiency: No intermediate type conversions or allocations
//   - Type Safety: Caller handles type assertion at deserialization time
//
// Returns redis.Nil error if the key or field does not exist.
func (db *redisCache) GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error) {
	val, err := db.client.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
	if err != nil {
		return nil, err
	}

	return val, nil
}

// DeleteCacheData removes an entire cache group from Redis.
// It deletes the hash key and all its associated fields in a single operation.
// This is useful for cache invalidation when the underlying data changes.
func (db *redisCache) DeleteCacheData(ctx context.Context, cacheGroupKey string) error {
	return db.client.Del(ctx, cacheGroupKey).Err()
}
