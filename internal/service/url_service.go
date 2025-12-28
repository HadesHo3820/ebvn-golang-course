// Package service provides business logic implementations for the application.
// This file contains the URL shortening service which generates unique codes
// for URLs and stores them in a repository for later retrieval.
package service

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/repository"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
)

// urlCodeLength defines the length of the generated short code for URLs.
// A 7-character alphanumeric code provides ~3.5 trillion unique combinations.
const (
	urlCodeLength = 7
)

// ShortenUrl defines the interface for URL shortening operations.
// Implementations of this interface handle the generation of short codes
// and persistence of URL mappings.
type ShortenUrl interface {
	// ShortenUrl generates a unique short code for the given URL
	// and stores the mapping in the repository.
	ShortenUrl(ctx context.Context, url string) (string, error)
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

// ShortenUrl generates a random alphanumeric code for the given URL,
// stores the code-to-URL mapping in the repository, and returns the code.
//
// The generated code is urlCodeLength characters long and uses a
// cryptographically secure random number generator.
//
// Returns:
//   - The generated short code on success.
//   - An error if code generation or storage fails.
func (s *shortenUrl) ShortenUrl(ctx context.Context, url string) (string, error) {
	// generate random code
	urlCode, err := stringutils.GenerateCode(urlCodeLength)
	if err != nil {
		return "", err
	}
	// store url in repository
	if err := s.repo.StoreUrl(ctx, urlCode, url); err != nil {
		return "", err
	}
	// return code
	return urlCode, nil
}
