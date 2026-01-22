// Package handler provides test utilities for HTTP handler testing.
// It offers helper functions for setting up Gin test contexts, creating
// JSON request bodies, and asserting JSON responses.
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	// Disable Gin debug mode for cleaner test output
	gin.SetMode(gin.TestMode)
}

// TestContext holds the test context and response recorder for handler testing.
// It provides a fluent API for configuring requests and simulating middleware behavior.
type TestContext struct {
	Ctx      *gin.Context
	Recorder *httptest.ResponseRecorder
}

// NewTestContext creates a new Gin test context for handler testing.
// It initializes a response recorder and creates a basic HTTP request.
//
// Parameters:
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - path: Request URL path
//
// Returns a TestContext that can be further configured with fluent methods.
func NewTestContext(method string, path string) *TestContext {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(method, path, nil)

	return &TestContext{
		Ctx:      ctx,
		Recorder: rec,
	}
}

// WithJSONBody sets a JSON request body on the test context.
// It marshals the provided data to JSON and sets the Content-Type header.
//
// Parameters:
//   - data: Any value that can be marshaled to JSON
//
// Returns the TestContext for method chaining.
func (tc *TestContext) WithJSONBody(data any) *TestContext {
	bodyBytes, _ := json.Marshal(data)
	tc.Ctx.Request = httptest.NewRequest(
		tc.Ctx.Request.Method,
		tc.Ctx.Request.URL.Path,
		bytes.NewReader(bodyBytes),
	)
	tc.Ctx.Request.Header.Set("Content-Type", "application/json")
	return tc
}

// WithJWTClaims sets JWT claims in the context, simulating JWT middleware behavior.
// This is useful for testing handlers that require authenticated users.
//
// Parameters:
//   - claims: JWT claims to set in the context (typically contains "sub" for user ID)
//
// Returns the TestContext for method chaining.
func (tc *TestContext) WithJWTClaims(claims jwt.MapClaims) *TestContext {
	if claims != nil {
		tc.Ctx.Set("claims", claims)
	}
	return tc
}

// WithHeader sets a header on the request.
//
// Parameters:
//   - key: Header name
//   - value: Header value
//
// Returns the TestContext for method chaining.
func (tc *TestContext) WithHeader(key, value string) *TestContext {
	if tc.Ctx.Request.Header == nil {
		tc.Ctx.Request.Header = make(http.Header)
	}
	tc.Ctx.Request.Header.Set(key, value)
	return tc
}
