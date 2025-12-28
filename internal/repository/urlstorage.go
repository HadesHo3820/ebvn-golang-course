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
type UrlStorage interface {
	// StoreUrl associates a code with a URL.
	StoreUrl(ctx context.Context, code, url string) error
	// GetUrl retrieves the URL for a given code.
	GetUrl(ctx context.Context, code string) (string, error)
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

// GetUrl retrieves the original URL from Redis using the provided code.
func (s *urlStorage) GetUrl(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}
