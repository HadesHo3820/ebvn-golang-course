package url

import (
	"errors"
	"net/http"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	handlertest "github.com/HadesHo3820/ebvn-golang-course/internal/test/handler"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
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
		requestBody    map[string]any
		setupMockSvc   func(ctx *gin.Context) *mocks.ShortenUrl
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:        "success - shorten URL with valid expiration",
			requestBody: fixture.DefaultShortenURLBody(),
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenUrl {
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
			name:        "success - shorten URL with custom expiration",
			requestBody: fixture.DefaultShortenURLBody(fixture.WithFieldAny("url", "https://google.com")),
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					ctx,
					"https://google.com",
					3600,
				).Return("xyz7890", nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Shorten URL generated successfully!",
				"code":    "xyz7890",
			},
		},
		{
			name:        "bad request - missing URL",
			requestBody: fixture.DefaultShortenURLBody(fixture.WithFieldAny("url", nil)),
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenUrl {
				// No mock expectations set - request fails validation before reaching service layer.
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:        "bad request - invalid URL format",
			requestBody: fixture.DefaultShortenURLBody(fixture.WithFieldAny("url", "not-a-valid-url")),
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenUrl {
				// No mock expectations set - request fails validation before reaching service layer.
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:        "bad request - negative expiration",
			requestBody: fixture.DefaultShortenURLBody(fixture.WithFieldAny("exp", -100)),
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenUrl {
				// No mock expectations set - request fails validation before reaching service layer.
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:        "internal server error - service failure",
			requestBody: fixture.DefaultShortenURLBody(),
			setupMockSvc: func(ctx *gin.Context) *mocks.ShortenUrl {
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
				"message": response.InternalErrMessage,
			},
		},
	}

	// Execute each test case in a subtest for better isolation and reporting.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test context with JSON body using helper
			testCtx := handlertest.NewTestContext(http.MethodPost, "/v1/links/shorten").
				WithJSONBody(tc.requestBody)

			// Setup the mock service
			svcMock := tc.setupMockSvc(testCtx.Ctx)

			// Create the handler with the mock service
			handler := NewUrlHandler(svcMock)

			// Call the handler
			handler.ShortenUrl(testCtx.Ctx)

			// Assert response using helper
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
