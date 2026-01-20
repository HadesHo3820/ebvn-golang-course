package url

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUrlShortenHandler_ShortenUrl validates the ShortenUrl handler.
// It uses table-driven tests to cover the following scenarios:
//   - Success cases: valid URL with default and custom expiration times
//   - Validation errors: missing URL, invalid URL format, negative expiration
//   - Service errors: handling failures from the service layer
//
// Each test case sets up an HTTP request and a mock service, then verifies
// that the handler returns the expected status code and response body.
func TestUrlShortenHandler_ShortenUrl(t *testing.T) {
	t.Parallel()

	// Set Gin to test mode to reduce noise in test output
	gin.SetMode(gin.TestMode)

	// testCases defines a table of test scenarios for the ShortenUrl handler.
	// Each test case contains:
	//   - name: descriptive name for the test scenario
	//   - setupRequest: factory function to create the HTTP request with specific body
	//   - setupMockSvc: factory function to configure mock service expectations
	//   - expectedStatus: the expected HTTP status code
	//   - expectedBody: the expected JSON response body (nil if not checked)
	testCases := []struct {
		name           string
		setupRequest   func() *http.Request
		setupMockSvc   func(ctx context.Context) *mocks.ShortenUrl
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success - shorten URL with valid expiration",
			setupRequest: func() *http.Request {
				body := map[string]any{
					"url": "https://example.com",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					// context matcher
					ctx,
					"https://example.com",
					3600,
				).Return("abc1234", nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Shorten URL generated successfully!",
				"code":    "abc1234",
			},
		},
		{
			name: "success - shorten URL with custom expiration",
			setupRequest: func() *http.Request {
				body := map[string]any{
					"url": "https://google.com",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					ctx,
					"https://google.com",
					3600,
				).Return("xyz7890", nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Shorten URL generated successfully!",
				"code":    "xyz7890",
			},
		},
		{
			name: "bad request - missing URL",
			setupRequest: func() *http.Request {
				body := map[string]any{
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				// No mock expectations set - request fails validation before reaching service layer.
				// The mock is still created to satisfy the handler constructor.
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Will contain error message
		},
		{
			name: "bad request - invalid URL format",
			setupRequest: func() *http.Request {
				body := map[string]any{
					"url": "not-a-valid-url",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				// No mock expectations set - request fails validation before reaching service layer.
				// The mock is still created to satisfy the handler constructor.
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Will contain error message
		},
		{
			name: "bad request - negative expiration",
			setupRequest: func() *http.Request {
				body := map[string]any{
					"url": "https://example.com",
					"exp": -100,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				// No mock expectations set - request fails validation before reaching service layer.
				// The mock is still created to satisfy the handler constructor.
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Will contain error message
		},
		{
			name: "internal server error - service failure",
			setupRequest: func() *http.Request {
				body := map[string]any{
					"url": "https://example.com",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					ctx,
					"https://example.com",
					3600,
				).Return("", errors.New("redis connection failed")).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"message": "internal server error",
			},
		},
	}

	// Execute each test case in a subtest for better isolation and reporting.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run subtests in parallel for faster execution

			// Create a ResponseRecorder to capture the HTTP response from the handler.
			// This allows us to inspect status codes and response bodies after the call.
			rec := httptest.NewRecorder()

			// Create a Gin test context bound to the recorder.
			// This simulates the Gin framework's request handling environment.
			gctx, _ := gin.CreateTestContext(rec)

			// Setup the request
			req := tc.setupRequest()
			req.Header.Set("Content-Type", "application/json")
			gctx.Request = req

			// Setup the mock service
			svcMock := tc.setupMockSvc(gctx)

			// Create the handler with the mock service
			handler := NewUrlHandler(svcMock)

			// Call the handler
			handler.ShortenUrl(gctx)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body if expected
			if tc.expectedBody != nil {
				var actualBody map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, actualBody)
			}
		})
	}
}
