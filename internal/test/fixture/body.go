// Package fixture provides test fixtures for setting up test data.
// This file contains helpers for creating request body maps used in handler tests.
package fixture

import "github.com/golang-jwt/jwt/v5"

// BodyMod is a function that modifies a request body map with string values.
// Use with the Default*Body functions to customize test data.
type BodyMod func(body map[string]string)

// BodyModAny is a function that modifies a request body map with any values.
// Use with functions like DefaultShortenURLBody to customize test data.
type BodyModAny func(body map[string]any)

// ClaimsMod is a function that modifies JWT claims.
// Use with DefaultJWTClaims to customize test claims.
type ClaimsMod func(claims jwt.MapClaims)

// WithField returns a modifier that sets a specific field value.
// Use empty string value to delete the field from the body.
//
// Example usage:
//
//	DefaultRegisterBody(WithField("email", "custom@example.com")) // override
//	DefaultRegisterBody(WithField("username", ""))               // delete field
func WithField(key, value string) BodyMod {
	return func(body map[string]string) {
		if value == "" {
			delete(body, key)
		} else {
			body[key] = value
		}
	}
}

// WithFieldAny returns a modifier that sets a specific field value for map[string]any.
// Use nil value to delete the field from the body.
func WithFieldAny(key string, value any) BodyModAny {
	return func(body map[string]any) {
		if value == nil {
			delete(body, key)
		} else {
			body[key] = value
		}
	}
}

// WithClaim returns a modifier that sets a specific claim value.
// Use nil value to delete the claim.
func WithClaim(key string, value any) ClaimsMod {
	return func(claims jwt.MapClaims) {
		if value == nil {
			delete(claims, key)
		} else {
			claims[key] = value
		}
	}
}

// DefaultRegisterBody returns a valid user registration request body.
// Pass optional BodyMod functions to customize specific fields.
//
// Default values:
//   - username: "testuser"
//   - password: "Password1!" (meets complexity requirements)
//   - display_name: "Test User"
//   - email: "test@example.com"
func DefaultRegisterBody(mods ...BodyMod) map[string]string {
	body := map[string]string{
		"username":     "testuser",
		"password":     "Password1!",
		"display_name": "Test User",
		"email":        "test@example.com",
	}
	for _, mod := range mods {
		mod(body)
	}
	return body
}

// DefaultLoginBody returns a valid login request body.
// Pass optional BodyMod functions to customize specific fields.
//
// Default values:
//   - username: "testuser"
//   - password: "password123" (8+ characters)
func DefaultLoginBody(mods ...BodyMod) map[string]string {
	body := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	for _, mod := range mods {
		mod(body)
	}
	return body
}

// DefaultUpdateUserBody returns a valid update user request body.
// Pass optional BodyMod functions to customize specific fields.
//
// Default values:
//   - display_name: "New Display Name"
//   - email: "new@example.com"
func DefaultUpdateUserBody(mods ...BodyMod) map[string]string {
	body := map[string]string{
		"display_name": "New Display Name",
		"email":        "new@example.com",
	}
	for _, mod := range mods {
		mod(body)
	}
	return body
}

// DefaultShortenURLBody returns a valid URL shorten request body.
// Pass optional BodyModAny functions to customize specific fields.
//
// Default values:
//   - url: "https://example.com"
//   - exp: 3600 (1 hour in seconds)
func DefaultShortenURLBody(mods ...BodyModAny) map[string]any {
	body := map[string]any{
		"url": "https://example.com",
		"exp": 3600,
	}
	for _, mod := range mods {
		mod(body)
	}
	return body
}

// DefaultJWTClaims returns default JWT claims for testing authenticated endpoints.
// Pass optional ClaimsMod functions to customize specific claims.
//
// Default values:
//   - sub: "test-user-id"
func DefaultJWTClaims(mods ...ClaimsMod) jwt.MapClaims {
	claims := jwt.MapClaims{
		"sub": "test-user-id",
	}
	for _, mod := range mods {
		mod(claims)
	}
	return claims
}
