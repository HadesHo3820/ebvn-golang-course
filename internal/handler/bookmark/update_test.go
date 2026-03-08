package bookmark

import (
	"context"
	"errors"
	"net/http"
	"testing"

	serviceMocks "github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark/mocks"
	handlertest "github.com/HadesHo3820/ebvn-golang-course/internal/test/handler"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

const testBookmarkIDUpdate = "f47ac10b-58cc-4372-a567-0e02b2c3d479"

func TestBookmarkHandler_UpdateBookmark(t *testing.T) {
	t.Parallel()

	// Disable Gin debug mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		jwtClaims      jwt.MapClaims
		uriParams      map[string]string
		inputBody      any
		setupMockSvc   func(t *testing.T, ctx context.Context) *serviceMocks.Service
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success - update bookmark",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			uriParams: map[string]string{"id": testBookmarkIDUpdate},
			inputBody: map[string]any{"description": "Updated Description", "url": "https://updated.com"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("UpdateBookmark",
					ctx,
					testBookmarkIDUpdate,
					testUserID,
					"Updated Description",
					"https://updated.com",
				).Return(nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"data":    nil,
				"message": "Bookmark updated successfully",
			},
		},
		{
			name:      "error - missing JWT claims",
			jwtClaims: nil,
			uriParams: map[string]string{"id": testBookmarkIDUpdate},
			inputBody: map[string]any{"description": "Description", "url": "https://example.com"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"message": "Invalid token",
			},
		},
		{
			name: "error - invalid input (missing URL)",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			uriParams: map[string]string{"id": testBookmarkIDUpdate},
			inputBody: map[string]any{"description": "Description"}, // Missing URL
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": response.InputErrMessage,
				"details": []any{"URL is invalid (required)"},
			},
		},
		{
			name: "error - invalid input (invalid UUID)",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			uriParams: map[string]string{"id": "not-a-uuid"},
			inputBody: map[string]any{"description": "Description", "url": "https://example.com"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": response.InputErrMessage,
				"details": []any{"ID is invalid (uuid)"},
			},
		},
		{
			name: "error - bookmark not found",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			uriParams: map[string]string{"id": testBookmarkIDUpdate},
			inputBody: map[string]any{"description": "Description", "url": "https://example.com"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("UpdateBookmark",
					ctx,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(dbutils.ErrNotFoundType)
				return svcMock
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]any{
				"message": "Bookmark not found",
			},
		},
		{
			name: "error - service failure",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			uriParams: map[string]string{"id": testBookmarkIDUpdate},
			inputBody: map[string]any{"description": "Description", "url": "https://example.com"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("UpdateBookmark",
					ctx,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(errors.New("service error"))
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

			// Create test context with JWT claims, URI params, and JSON body
			testCtx := handlertest.NewTestContext(http.MethodPut, "/v1/bookmarks/:id").
				WithJWTClaims(tc.jwtClaims).
				WithURIParams(tc.uriParams).
				WithJSONBody(tc.inputBody)

			// Setup mock service
			svcMock := tc.setupMockSvc(t, testCtx.Ctx)

			// Create handler with mock service
			handler := NewHandler(svcMock)

			// Call the handler
			handler.UpdateBookmark(testCtx.Ctx)

			// Assert response
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
