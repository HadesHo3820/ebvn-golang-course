// Package cache provides an abstraction layer for cache storage operations.
// It defines a common interface that can be implemented by different cache backends
// (e.g., Redis, Memcached, in-memory), allowing the application to switch between
// cache implementations without modifying the business logic.
package cache

import (
	"context"
	"time"
)

// DB defines the interface for cache database operations.
// It provides methods for setting, getting, and deleting cached data using
// a hierarchical key structure (cacheGroupKey -> cacheKey -> value).
//
// The interface uses []byte as the return type for GetCacheData because:
//   - Type Agnosticism: []byte is a universal data representation that works with
//     any serialized format (JSON, Protobuf, MessagePack, etc.), allowing the caller
//     to deserialize into their specific type.
//   - Efficiency: Returning raw bytes avoids unnecessary intermediate conversions
//     and allocations. The caller can unmarshal directly into their target struct.
//   - Flexibility: Different consumers may need to deserialize the same cached data
//     into different types, which []byte enables without coupling to a specific type.
//   - Redis Compatibility: Redis stores values as byte strings internally, so []byte
//     provides a natural mapping to the underlying storage format.
//go:generate mockery --name DB --filename db.go --outpkg mock_cache
type DB interface {
	SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value any, exp time.Duration) error
	GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error)
	DeleteCacheData(ctx context.Context, cacheGroupKey string) error
}
