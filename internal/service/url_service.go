// Package service provides business logic implementations for the application.
// This file contains the URL shortening service which generates unique codes
// for URLs and stores them in a repository for later retrieval.
package service

import (
	"context"
	"errors"
	"fmt"

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
	// Returns ErrCodeNotFound if the code does not exist in the repository.
	GetUrl(ctx context.Context, code string) (string, error)
}

// shortenUrl is the concrete implementation of the ShortenUrl interface.
// It uses a UrlStorage repository for persisting URL-to-code mappings.
type shortenUrl struct {
	repo   repository.UrlStorage
	keyGen stringutils.KeyGenerator
}

// NewShortenUrl creates a new instance of the ShortenUrl service.
// It requires a UrlStorage repository for storing shortened URL mappings.
func NewShortenUrl(repo repository.UrlStorage, keyGen stringutils.KeyGenerator) ShortenUrl {
	return &shortenUrl{repo: repo, keyGen: keyGen}
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
// It queries the repository and translates redis.Nil errors to ErrCodeNotFound
// for a cleaner abstraction that doesn't leak storage implementation details.
//
// Returns:
//   - The original URL if the code exists.
//   - ErrCodeNotFound if the code does not exist.
//   - Other errors for repository/connection failures.
func (s *shortenUrl) GetUrl(ctx context.Context, code string) (string, error) {
	url, err := s.repo.GetUrl(ctx, code)
	// redis.Nil is returned when the key does not exist
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}
	return url, err
}
