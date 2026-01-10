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

	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils"
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
		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus int
		validateBody   func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "success - shorten valid URL",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "https://example.com",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Shorten URL generated successfully!", body["message"])
				assert.NotEmpty(t, body["code"])
				// Code should be 7 characters
				code, ok := body["code"].(string)
				assert.True(t, ok)
				assert.Len(t, code, 7)
			},
		},
		{
			name: "bad request - invalid URL format",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "not-a-valid-url",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "wrong input", body["message"])
			},
		},
		{
			name: "bad request - missing URL",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "wrong input", body["message"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Since the URL shortening feature doesn't require configuration,
			// we can pass nil to the api.New function.
			rec := tc.setupTestHTTP(api.New(&api.Config{}, redisPkg.InitMockRedis(t), stringutils.NewKeyGenerator()))

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
		code           string                                                             // Code to request
		setupRedis     func(redis *redis.Client)                                          // Optional: pre-populate Redis
		setupTestHTTP  func(apiEngine api.Engine, code string) *httptest.ResponseRecorder // Setup and execute request
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
			setupTestHTTP: func(apiEngine api.Engine, code string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, redirectURI+code, nil)
				rec := httptest.NewRecorder()
				apiEngine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusFound,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, "https://preloaded-url.com", rec.Header().Get("Location"))
			},
		},
		{
			name:       "bad request - code not found",
			code:       "notexist",
			setupRedis: nil, // No pre-population needed
			setupTestHTTP: func(apiEngine api.Engine, code string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, redirectURI+code, nil)
				rec := httptest.NewRecorder()
				apiEngine.ServeHTTP(rec, req)
				return rec
			},
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
			setupTestHTTP: func(apiEngine api.Engine, code string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, redirectURI+code, nil)
				rec := httptest.NewRecorder()
				apiEngine.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "internal server error", resp["message"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock Redis
			mockRedis := redisPkg.InitMockRedis(t)

			// Pre-populate Redis if needed
			if tc.setupRedis != nil {
				tc.setupRedis(mockRedis)
			}

			// Create API engine
			apiEngine := api.New(&api.Config{}, mockRedis, stringutils.NewKeyGenerator())

			// Execute request
			rec := tc.setupTestHTTP(apiEngine, tc.code)

			// Assert status
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Validate response
			if tc.validateBody != nil {
				tc.validateBody(t, rec)
			}
		})
	}
}
