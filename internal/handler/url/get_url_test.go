package url

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUrlShortenHandler_GetUrl validates the GetUrl handler.
// It uses table-driven tests to cover the following scenarios:
//   - Success case: valid code returns 302 redirect
//   - Validation errors: empty code returns 400
//   - Service errors: code not found (ErrCodeNotFound) returns 400, other errors return 500
//
// Each test case sets up an HTTP request with a path parameter and a mock service,
// then verifies that the handler returns the expected status code and response.
func TestUrlShortenHandler_GetUrl(t *testing.T) {
	t.Parallel()

	// Set Gin to test mode to reduce noise in test output
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		code           string                                      // Path parameter value
		setupMockSvc   func(ctx context.Context) *mocks.ShortenUrl // Mock service setup
		expectedStatus int
		expectedBody   map[string]any // nil for redirect responses
		expectedHeader string         // Expected Location header for redirects
	}{
		{
			name: "success - redirects to original URL",
			code: "abc1234",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("GetUrl", ctx, "abc1234").
					Return("https://example.com", nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusFound,
			expectedBody:   nil,
			expectedHeader: "https://example.com",
		},
		{
			name: "bad request - empty code",
			code: "",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				// No mock expectations - request fails validation before reaching service
				return mocks.NewShortenUrl(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": "wrong format",
			},
		},
		{
			name: "bad request - code not found",
			code: "notfound",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("GetUrl", ctx, "notfound").
					Return("", service.ErrCodeNotFound).Once()
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": "url not found",
			},
		},
		{
			name: "internal server error - service failure",
			code: "abc1234",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("GetUrl", ctx, "abc1234").
					Return("", errors.New("redis connection failed")).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"message": response.InternalErrMessage,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a ResponseRecorder to capture the HTTP response
			rec := httptest.NewRecorder()

			// Create a Gin test context
			gctx, _ := gin.CreateTestContext(rec)

			// Setup the request with the code path parameter
			req := httptest.NewRequest(http.MethodGet, "/v1/links/"+tc.code, nil)
			gctx.Request = req

			// Set the path parameter for Gin to read via c.Param("code")
			gctx.Params = gin.Params{
				{Key: "code", Value: tc.code},
			}

			// Setup the mock service
			svcMock := tc.setupMockSvc(gctx)

			// Create the handler with the mock service
			handler := NewUrlHandler(svcMock)

			// Call the handler
			handler.GetUrl(gctx)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body if expected (for error responses)
			if tc.expectedBody != nil {
				var actualBody map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, actualBody)
			}

			// Assert Location header for redirect responses
			if tc.expectedHeader != "" {
				assert.Equal(t, tc.expectedHeader, rec.Header().Get("Location"))
			}
		})
	}
}
