// Package handler contains unit tests for the URL shortening handler.
// These tests use mocks to isolate the handler logic from the service layer.
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUrlShortenHandler_ShortenUrl validates the ShortenUrl handler.
// It uses table-driven tests to cover success, validation errors, and service errors.
func TestUrlShortenHandler_ShortenUrl(t *testing.T) {
	t.Parallel()

	// Set Gin to test mode to reduce noise in test output
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupRequest   func() *http.Request
		setupMockSvc   func() *mocks.ShortenUrl
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success - shorten URL with default expiration",
			setupRequest: func() *http.Request {
				body := map[string]interface{}{
					"url": "https://example.com",
					"exp": 0,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func() *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					// context matcher
					mock.Anything,
					"https://example.com",
					0,
				).Return("abc1234", nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Shorten URL generated successfully!",
				"code":    "abc1234",
			},
		},
		{
			name: "success - shorten URL with custom expiration",
			setupRequest: func() *http.Request {
				body := map[string]interface{}{
					"url": "https://google.com",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func() *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					mock.Anything,
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
				body := map[string]interface{}{
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func() *mocks.ShortenUrl {
				// No mock setup needed - validation fails before service call
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Will contain error message
		},
		{
			name: "bad request - invalid URL format",
			setupRequest: func() *http.Request {
				body := map[string]interface{}{
					"url": "not-a-valid-url",
					"exp": 0,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func() *mocks.ShortenUrl {
				// No mock setup needed - validation fails before service call
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Will contain error message
		},
		{
			name: "bad request - negative expiration",
			setupRequest: func() *http.Request {
				body := map[string]interface{}{
					"url": "https://example.com",
					"exp": -100,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func() *mocks.ShortenUrl {
				// No mock setup needed - validation fails before service call
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil, // Will contain error message
		},
		{
			name: "internal server error - service failure",
			setupRequest: func() *http.Request {
				body := map[string]interface{}{
					"url": "https://example.com",
					"exp": 0,
				}
				jsonBody, _ := json.Marshal(body)
				return httptest.NewRequest(http.MethodPost, "/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func() *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					mock.Anything,
					"https://example.com",
					0,
				).Return("", errors.New("redis connection failed")).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "redis connection failed",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a recorder to capture the response
			rec := httptest.NewRecorder()

			// Create a Gin test context
			gctx, _ := gin.CreateTestContext(rec)

			// Setup the request
			req := tc.setupRequest()
			req.Header.Set("Content-Type", "application/json")
			gctx.Request = req

			// Setup the mock service
			svcMock := tc.setupMockSvc()

			// Create the handler with the mock service
			handler := NewUrlShorten(svcMock)

			// Call the handler
			handler.ShortenUrl(gctx)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body if expected
			if tc.expectedBody != nil {
				var actualBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, actualBody)
			}

			// Verify mock expectations
			svcMock.AssertExpectations(t)
		})
	}
}
