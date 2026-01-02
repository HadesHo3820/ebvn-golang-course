// Package repository provides the data access layer for the application.
package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// urlExpTime defines the expiration duration for URLs stored in the cache.
const (
	urlExpTime = 24 * time.Hour
)

// UrlStorage defines the interface for storing and retrieving URLs.
//
//go:generate mockery --name UrlStorage --filename url_storage.go
type UrlStorage interface {
	// StoreUrl associates a code with a URL.
	StoreUrl(ctx context.Context, code, url string) error
	// StoreUrlIfNotExists atomically stores the URL only if the code doesn't exist.
	// Returns true if stored successfully, false if the code already exists.
	StoreUrlIfNotExists(ctx context.Context, code, url string, exp int) (bool, error)
	// GetUrl retrieves the URL for a given code.
	GetUrl(ctx context.Context, code string) (string, error)
	// Exists checks if a code is already stored.
	Exists(ctx context.Context, code string) (bool, error)
}

// urlStorage is a Redis-backed implementation of UrlStorage.
type urlStorage struct {
	c *redis.Client
}

// NewUrlStorage creates a new instance of UrlStorage.
func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

// StoreUrl saves the code and URL pair in Redis with an expiration time.
func (s *urlStorage) StoreUrl(ctx context.Context, code, url string) error {
	return s.c.Set(ctx, code, url, urlExpTime).Err()
}

// StoreUrlIfNotExists atomically stores the URL using Redis SETNX.
// This operation is atomic: the key is only set if it doesn't already exist.
// If exp > 0, uses the provided expiration in seconds; otherwise uses the default urlExpTime.
// Returns true if the URL was stored (key was new), false if the code already exists.
func (s *urlStorage) StoreUrlIfNotExists(ctx context.Context, code, url string, exp int) (bool, error) {
	expDuration := urlExpTime
	if exp > 0 {
		expDuration = time.Duration(exp) * time.Second
	}
	return s.c.SetNX(ctx, code, url, expDuration).Result()
}

// GetUrl retrieves the original URL from Redis using the provided code.
func (s *urlStorage) GetUrl(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}

// Exists checks if a code exists in Redis.
func (s *urlStorage) Exists(ctx context.Context, code string) (bool, error) {
	result, err := s.c.Exists(ctx, code).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
