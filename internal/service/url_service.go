// Package service provides business logic implementations for the application.
// This file contains the URL shortening service which generates unique codes
// for URLs and stores them in a repository for later retrieval.
package service

import (
	"context"
	"fmt"

	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
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
//go:generate mockery --name ShortenUrl --output ./mocks --filename shorten_url.go
type ShortenUrl interface {
	// ShortenUrl generates a unique short code for the given URL
	// and stores the mapping in the repository.
	ShortenUrl(ctx context.Context, url string, exp int) (string, error)
}

// shortenUrl is the concrete implementation of the ShortenUrl interface.
// It uses a UrlStorage repository for persisting URL-to-code mappings.
type shortenUrl struct {
	repo repository.UrlStorage
}

// NewShortenUrl creates a new instance of the ShortenUrl service.
// It requires a UrlStorage repository for storing shortened URL mappings.
func NewShortenUrl(repo repository.UrlStorage) ShortenUrl {
	return &shortenUrl{repo: repo}
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
		urlCode, err := stringutils.GenerateCode(urlCodeLength)
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
