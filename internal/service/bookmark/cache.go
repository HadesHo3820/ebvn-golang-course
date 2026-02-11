// Package bookmark provides bookmark management services with optional caching support.
package bookmark

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/repository/cache"
	"github.com/rs/zerolog/log"
)

const (
	getBookmarksCacheGroupFormat = "get_bookmarks_%s" // get_bookmarks_<userID>
	getBookmarksCacheKeyFormat   = "%d_%d"            // <page>_<limit>
	getBookmarksCacheDuration    = 24 * time.Hour
)

// ErrNilRequest is returned when GetBookmarks receives a nil request.
var ErrNilRequest = errors.New("request cannot be nil")

type bookmarkServiceWithCache struct {
	s     Service
	cache cache.DB
}

// NewServiceWithCache wraps a Service with Redis caching capabilities.
// It returns a Service that implements Cache-Aside for reads and
// Write-Invalidate for mutations.
func NewServiceWithCache(s Service, cacheDB cache.DB) Service {
	return &bookmarkServiceWithCache{
		s:     s,
		cache: cacheDB,
	}
}

// CreateBookmark creates a new bookmark and automatically invalidates the user's bookmark cache.
// This method uses the "Write-Invalidate" strategy:
// 1. Primary Write: The bookmark is first created in the persistent database.
// 2. Cache Invalidation: The user's entire bookmark cache group is deleted.
//
// Error Handling:
//   - Database errors block the operation and are returned to the caller.
//   - Cache invalidation errors are logged but ignored (Best Effort), ensuring that
//     transient cache issues do not prevent users from creating bookmarks.
func (s *bookmarkServiceWithCache) CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	// 1. Create in DB first (Source of Truth)
	newBookmark, err := s.s.CreateBookmark(ctx, description, url, userID)
	if err != nil {
		return nil, err
	}

	// 2. Invalidate Cache (Best Effort)
	// We delete the cache AFTER successful DB write to ensure consistency.
	// If cache deletion fails, we log the error but don't fail the request,
	// preserving availability (the user still gets their bookmark created).
	if err := s.cache.DeleteCacheData(ctx, fmt.Sprintf(getBookmarksCacheGroupFormat, userID)); err != nil {
		log.Error().Err(err).
			Str("userID", userID).
			Msg("Failed to invalidate cache after bookmark creation")
	}

	return newBookmark, nil
}

// GetBookmarks retrieves a list of bookmarks with caching support.
// It follows the Cache-Aside pattern:
// 1. Check Cache: Try to get data from Redis.
// 2. DB Fallback: If cache miss or error, fetch from the database.
// 3. Generic Return: Return data to user.
// 4. Async Cache: Populate cache with new data for subsequent requests.
func (s *bookmarkServiceWithCache) GetBookmarks(ctx context.Context, userID string, req *dto.Request) (*dto.Response[*model.Bookmark], error) {
	// 1. Validate the request to prevent potential panics.
	if req == nil {
		return nil, ErrNilRequest
	}

	// 2. Sanitize pagination parameters and generate cache keys.
	// Sanitization ensures consistent cache keys (e.g., "1_10" instead of "0_0")
	// even when query params are missing.
	req.Sanitize()

	// - groupKey: groups all pages for this user (allows bulk invalidation).
	// - key: identifies specific page & limit combination.
	groupKey := fmt.Sprintf(getBookmarksCacheGroupFormat, userID)
	key := fmt.Sprintf(getBookmarksCacheKeyFormat, req.Page, req.Limit)

	// 3. Try to retrieve data from the cache.
	cachedData, err := s.cache.GetCacheData(ctx, groupKey, key)
	if err == nil && len(cachedData) > 0 {
		response := &dto.Response[*model.Bookmark]{}

		// 4. Cache Hit: Attempt to deserialize.
		// If successful, return immediately. If failed, log warning and fall through to DB.
		if err := json.Unmarshal(cachedData, response); err != nil {
			log.Warn().Err(err).
				Str("groupKey", groupKey).
				Str("key", key).
				Msg("Failed to unmarshal cached data")
		} else {
			return response, nil
		}
	}

	// 5. Cache Miss: Fetch fresh data from the underlying service (Database).
	result, err := s.s.GetBookmarks(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// 6. Cache Population: Store the new result.
	// Marshal the result to JSON.
	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Warn().Err(err).
			Msg("Failed to marshal result for caching")
	} else {
		// Cache operations inherit the request context. If the request has a short deadline and Redis is slow,
		// you might fail to cache even though the DB succeeded. This is usually fine (cache is best-effort),
		// but for long-running batch operations, consider using a detached context for cache writes
		cacheCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Save to cache with expiration. Log errors but don't fail the request.
		if cachedErr := s.cache.SetCacheData(cacheCtx, groupKey, key, resultBytes, getBookmarksCacheDuration); cachedErr != nil {
			log.Error().Err(cachedErr).
				Str("groupKey", groupKey).
				Str("key", key).
				Msg("Cannot cache this data")
		}
	}

	return result, nil
}

// UpdateBookmark updates an existing bookmark and invalidates the user's cache.
// Strategy: Write-Invalidate
// 1. Update DB (Source of Truth)
// 2. Invalidate Cache (Best Effort)
func (s *bookmarkServiceWithCache) UpdateBookmark(ctx context.Context, bookmarkID, userID, description, url string) error {
	// 1. DB Update
	if err := s.s.UpdateBookmark(ctx, bookmarkID, userID, description, url); err != nil {
		return err
	}

	// 2. Cache Invalidation
	if err := s.cache.DeleteCacheData(ctx, fmt.Sprintf(getBookmarksCacheGroupFormat, userID)); err != nil {
		log.Error().Err(err).
			Str("userID", userID).
			Msg("Failed to invalidate cache after bookmark update")
	}

	return nil
}

// DeleteBookmark removes a bookmark and invalidates the user's cache.
// Strategy: Write-Invalidate
// 1. Delete from DB (Source of Truth)
// 2. Invalidate Cache (Best Effort)
func (s *bookmarkServiceWithCache) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	// 1. DB Delete
	if err := s.s.DeleteBookmark(ctx, bookmarkID, userID); err != nil {
		return err
	}

	// 2. Cache Invalidation
	if err := s.cache.DeleteCacheData(ctx, fmt.Sprintf(getBookmarksCacheGroupFormat, userID)); err != nil {
		log.Error().Err(err).
			Str("userID", userID).
			Msg("Failed to invalidate cache after bookmark deletion")
	}

	return nil
}
