// Package endpoint provides integration tests for API endpoints.
//
// This file contains integration tests for the URL shortening endpoint,
// validating the full HTTP stack including routing, handlers, and real Redis integration.
package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestUrlShortenEndpoint validates the /links/shorten endpoint through the full HTTP stack.
//
// This is an integration test that exercises:
//   - HTTP routing configuration
//   - Request handling through the Gin engine
//   - Handler-to-service-to-repository delegation with real Redis
//   - JSON request parsing and response serialization
//
// Prerequisites:
//   - Redis must be running and accessible
//
// Test coverage includes:
//   - Verifying successful URL shortening with valid input
//   - Validating error response for invalid URL format
//   - Validating error response for missing required field
func TestUrlShortenEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		requestBody    map[string]any
		expectedStatus int
		validateBody   func(t *testing.T, body map[string]interface{})
	}{
		{
			name:           "success - shorten valid URL",
			requestBody:    fixture.DefaultShortenURLBody(),
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Shorten URL generated successfully!", body["message"])
				assert.NotEmpty(t, body["data"])
				// Code should be 7 characters
				code, ok := body["data"].(string)
				assert.True(t, ok)
				assert.Len(t, code, 7)
			},
		},
		{
			name:           "bad request - invalid URL format",
			requestBody:    fixture.DefaultShortenURLBody(fixture.WithFieldAny("url", "not-a-valid-url")),
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, response.InputErrMessage, body["message"])
			},
		},
		{
			name:           "bad request - missing URL",
			requestBody:    fixture.DefaultShortenURLBody(fixture.WithFieldAny("url", nil)),
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, response.InputErrMessage, body["message"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testEngine := NewTestEngine(&TestEngineOpts{
				T:       t,
				Fixture: &fixture.BookmarkCommonTestDB{},
			})

			jsonBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]any
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tc.validateBody != nil {
				tc.validateBody(t, resp)
			}
		})
	}
}

// TestGetUrlEndpoint validates the /links/redirect/:code endpoint through the full HTTP stack.
//
// This is an integration test that exercises:
//   - HTTP routing configuration for GET /v1/links/redirect/:code
//   - Request handling through the Gin engine
//   - Handler-to-service-to-repository delegation with real Redis
//   - Redirect responses (HTTP 302) for successful lookups
//
// Test coverage includes:
//   - Verifying successful URL retrieval and redirect after shortening
//   - Validating error response for non-existent code
func TestGetUrlEndpoint(t *testing.T) {
	t.Parallel()

	redirectURI := "/v1/links/redirect/"

	testCases := []struct {
		name           string
		code           string                    // Code to request
		setupRedis     func(redis *redis.Client) // Optional: pre-populate Redis
		expectedStatus int
		validateBody   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success - redirects to original URL",
			code: "preload1",
			setupRedis: func(r *redis.Client) {
				// Pre-populate Redis with a code-URL mapping
				r.Set(context.Background(), "preload1", "https://preloaded-url.com", 0)
			},
			expectedStatus: http.StatusFound,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, "https://preloaded-url.com", rec.Header().Get("Location"))
			},
		},
		{
			name:           "success - Redis miss, DB fallback hit",
			code:           "1", // base62.Encode(FixtureBookmarkOneCode=1) = "1"
			setupRedis:     nil, // Redis empty — forces DB fallback
			expectedStatus: http.StatusFound,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				// Should redirect to the bookmark URL seeded by BookmarkCommonTestDB fixture
				assert.Equal(t, fixture.FixtureBookmarkURL, rec.Header().Get("Location"))
			},
		},
		{
			name:           "bad request - code not found in Redis or DB",
			code:           "notexist",
			setupRedis:     nil, // No pre-population needed
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "url not found", resp["message"])
			},
		},
		{
			name: "internal server error - redis connection failure",
			code: "anycode1",
			setupRedis: func(r *redis.Client) {
				// Close the Redis connection to simulate a connection failure
				r.Close()
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, response.InternalErrMessage, resp["message"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testEngine := NewTestEngine(&TestEngineOpts{
				T:       t,
				Fixture: &fixture.BookmarkCommonTestDB{},
			})

			// Pre-populate Redis if needed
			if tc.setupRedis != nil {
				tc.setupRedis(testEngine.RedisClient)
			}

			// Execute request
			req := httptest.NewRequest(http.MethodGet, redirectURI+tc.code, nil)
			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Validate response
			if tc.validateBody != nil {
				tc.validateBody(t, rec)
			}
		})
	}
}
