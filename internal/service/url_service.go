// Package service provides business logic implementations for the application.
// This file contains the URL shortening service which generates unique codes
// for URLs and stores them in a repository for later retrieval.
// It also supports fallback to the bookmark database for codes not found in Redis.
package service

import (
	"context"
	"errors"
	"fmt"

	bookmarkSvc "github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/base62"

	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
	"github.com/redis/go-redis/v9"
)

// urlCodeLength defines the length of the generated short code for URLs.
// A 7-character alphanumeric code provides ~3.5 trillion unique combinations.
const (
	urlCodeLength = 7
	maxRetries    = 5 // Maximum attempts to generate a unique code
)

// ShortenUrl defines the interface for URL shortening operations.
// Implementations of this interface handle the generation of short codes
// and persistence of URL mappings.
//
//go:generate mockery --name ShortenUrl --filename shorten_url.go
type ShortenUrl interface {
	// ShortenUrl generates a unique short code for the given URL
	// and stores the mapping in the repository.
	ShortenUrl(ctx context.Context, url string, exp int) (string, error)

	// GetUrl retrieves the original URL associated with the given short code.
	// It first checks Redis for the code. If not found in Redis, it attempts
	// to base62-decode the code and look up the bookmark in the database.
	// Returns ErrCodeNotFound if the code does not exist in either store.
	GetUrl(ctx context.Context, code string) (string, error)
}

// shortenUrl is the concrete implementation of the ShortenUrl interface.
// It uses a UrlStorage repository for persisting URL-to-code mappings
// and falls back to the bookmark service for database-stored codes.
type shortenUrl struct {
	repo        repository.UrlStorage
	keyGen      stringutils.KeyGenerator
	bookmarkSvc bookmarkSvc.Service
}

// NewShortenUrl creates a new instance of the ShortenUrl service.
// It requires a UrlStorage repository for storing shortened URL mappings
// and a bookmark service for fallback database lookups.
func NewShortenUrl(repo repository.UrlStorage, keyGen stringutils.KeyGenerator, bmSvc bookmarkSvc.Service) ShortenUrl {
	return &shortenUrl{repo: repo, keyGen: keyGen, bookmarkSvc: bmSvc}
}

// ShortenUrl generates a unique alphanumeric code for the given URL,
// stores the code-to-URL mapping in the repository, and returns the code.
//
// The method attempts to generate a unique code up to maxRetries times.
// For each attempt, it uses an atomic SETNX operation to store the URL
// only if the code doesn't already exist. If a collision is detected
// (code already exists), it retries with a new code.
//
// The generated code is urlCodeLength characters long and uses a
// cryptographically secure random number generator.
//
// Returns:
//   - The generated short code on success.
//   - An error if code generation fails, storage fails, or max retries exceeded.
func (s *shortenUrl) ShortenUrl(ctx context.Context, url string, exp int) (string, error) {
	for range maxRetries {
		// generate random code
		urlCode, err := s.keyGen.GenerateCode(urlCodeLength)
		if err != nil {
			return "", err
		}

		// atomically store url if code doesn't exist (SETNX)
		stored, err := s.repo.StoreUrlIfNotExists(ctx, urlCode, url, exp)
		if err != nil {
			return "", err
		}
		if !stored {
			continue // collision detected, retry with new code
		}

		return urlCode, nil
	}

	return "", fmt.Errorf("failed to generate unique code after %d attempts", maxRetries)
}

// ErrCodeNotFound is a sentinel error returned when a short code
// does not exist in the repository. Callers should use errors.Is()
// to check for this specific error condition.
var ErrCodeNotFound = errors.New("code not found")

// GetUrl retrieves the original URL for a given short code.
// It first queries Redis for the code. If the code is not found in Redis
// (redis.Nil), it attempts to base62-decode the code and look it up
// in the bookmark database as a fallback.
//
// This enables the redirect endpoint to serve both:
//   - Redis-stored URL codes (7-char random codes from URL shortening)
//   - Database-stored bookmark codes (auto-incremented integers encoded as base62)
//
// Returns:
//   - The original URL if the code exists in either store.
//   - ErrCodeNotFound if the code does not exist in either store.
//   - Other errors for repository/connection failures.
func (s *shortenUrl) GetUrl(ctx context.Context, code string) (string, error) {
	// 1. Try Redis first (URL shortener codes)
	url, err := s.repo.GetUrl(ctx, code)
	if err == nil {
		return url, nil
	}

	// If the error is not "key not found", it's a real error
	if !errors.Is(err, redis.Nil) {
		return "", err
	}

	// 2. Redis miss: Try to decode as base62 and look up in bookmark DB
	decodedCode, decodeErr := base62.Decode(code)
	if decodeErr != nil {
		// Code is not a valid base62 string, treat as not found
		return "", ErrCodeNotFound
	}

	bookmarkUrl, bmErr := s.bookmarkSvc.GetUrlByCode(ctx, decodedCode)
	if errors.Is(bmErr, bookmarkSvc.ErrBookmarkNotFound) {
		return "", ErrCodeNotFound
	}
	if bmErr != nil {
		return "", bmErr
	}

	return bookmarkUrl, nil
}
