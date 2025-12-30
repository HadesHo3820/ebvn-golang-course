// Package repository provides the data access layer for the application.
// This file contains the health check repository for verifying external dependencies.
package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthChecker defines the interface for health check operations on external dependencies.
// This abstraction allows the service layer to verify connectivity without knowing
// the specific implementation details (Clean Architecture / Hexagonal pattern).
//
//go:generate mockery --name HealthChecker --dir ../../internal/repository --output ../../internal/service/mocks --filename health_checker.go
type HealthChecker interface {
	// Ping checks if the dependency is reachable and healthy.
	// Returns nil if healthy, error otherwise.
	Ping(ctx context.Context) error
}

// redisHealthChecker is a Redis-backed implementation of HealthChecker.
type redisHealthChecker struct {
	client *redis.Client
}

// NewRedisHealthChecker creates a new HealthChecker that verifies Redis connectivity.
func NewRedisHealthChecker(client *redis.Client) HealthChecker {
	return &redisHealthChecker{client: client}
}

// Ping checks if Redis is reachable by sending a PING command.
func (r *redisHealthChecker) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
