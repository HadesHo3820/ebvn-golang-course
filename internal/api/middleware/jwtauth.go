// Package middleware provides HTTP middleware components for the API layer.
// It includes authentication, logging, and other cross-cutting concerns
// that are applied to routes via Gin's middleware chain.
package middleware

import (
	"net/http"
	"strings"

	"github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils"
	"github.com/gin-gonic/gin"
)

// JWTAuth defines the interface for JWT-based authentication middleware.
// It provides a method to create a Gin handler function that validates
// JWT tokens in incoming requests.
type JWTAuth interface {
	// JWTAuth returns a Gin middleware handler that validates JWT tokens.
	// It extracts the token from the Authorization header, validates it,
	// and stores the claims in the Gin context for downstream handlers.
	JWTAuth() gin.HandlerFunc
}

// jwtAuth is the concrete implementation of the JWTAuth interface.
// It uses a JWTValidator to verify token signatures and expiration.
type jwtAuth struct {
	jwtValidator jwtutils.JWTValidator
}

// NewJWTAuth creates a new JWTAuth middleware instance.
//
// Parameters:
//   - jwtValidator: The validator used to verify JWT tokens
//
// Returns:
//   - JWTAuth: A new middleware instance ready to be used with Gin routes
//
// Example:
//
//	jwtMiddleware := middleware.NewJWTAuth(jwtValidator)
//	router.Use(jwtMiddleware.JWTAuth())
func NewJWTAuth(jwtValidator jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{
		jwtValidator: jwtValidator,
	}
}

// JWTAuth returns a Gin middleware handler function that performs JWT authentication.
//
// The middleware performs the following steps:
//  1. Extracts the Authorization header from the request
//  2. Validates the header format (must be "Bearer <token>")
//  3. Validates the JWT token using the configured validator
//  4. Stores the token claims in the Gin context under the key "claims"
//  5. Calls the next handler in the chain if validation succeeds
//
// On failure, the middleware aborts the request with HTTP 401 Unauthorized
// and returns a JSON error response.
//
// Downstream handlers can access the claims using:
//
//	claims, exists := c.Get("claims")
//	if exists {
//	    tokenClaims := claims.(*jwtutils.TokenClaims)
//	    userID := tokenClaims.UserID
//	}
func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header from the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Parse the header to extract the Bearer token
		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}
		tokenString := parts[1]

		// Validate the token signature, expiration, and claims
		tokenClaims, err := j.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Store claims in context for downstream handlers to access
		c.Set("claims", tokenClaims)

		// Continue to the next handler in the chain
		c.Next()
	}
}
