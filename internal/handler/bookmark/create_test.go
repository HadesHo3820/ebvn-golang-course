package bookmark

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	serviceMocks "github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark/mocks"
	handlertest "github.com/HadesHo3820/ebvn-golang-course/internal/test/handler"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

const (
	testUserID       = "test-user-id"
	testBookmarkDesc = "My Bookmark"
	testBookmarkURL  = "https://example.com"
	testBookmarkCode = "abc123456"
)

var (
	testBookmarkLongURL = "https://example.com/" + strings.Repeat("a", 2050)
)

func TestBookmarkHandler_CreateBookmark(t *testing.T) {
	t.Parallel()

	// Disable Gin debug mode
	gin.SetMode(gin.TestMode)

	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		jwtClaims      jwt.MapClaims
		inputBody      any
		setupMockSvc   func(t *testing.T, ctx context.Context) *serviceMocks.Service
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success - create bookmark",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			inputBody: map[string]any{"description": "My Bookmark", "url": "https://example.com"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("CreateBookmark",
					ctx,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(&model.Bookmark{
					Base: model.Base{
						ID:        "bm-1",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					},
					Description: testBookmarkDesc,
					URL:         testBookmarkURL,
					Code:        testBookmarkCode,
					UserID:      testUserID,
				}, nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Bookmark created successfully",
				"data": map[string]any{
					"id":          "bm-1",
					"description": testBookmarkDesc,
					"url":         testBookmarkURL,
					"code":        testBookmarkCode,
					"user_id":     testUserID,
					"created_at":  fixedTime.Format(time.RFC3339Nano),
					"updated_at":  fixedTime.Format(time.RFC3339Nano),
				},
			},
		},
		{
			name:      "error - missing JWT claims",
			jwtClaims: nil,
			inputBody: map[string]any{"description": "My Bookmark", "url": "https://example.com"},
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
			inputBody: map[string]any{"description": "My Bookmark"}, // Missing URL
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
			name: "error - invalid input (invalid URL format)",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			inputBody: map[string]any{"description": "My Bookmark", "url": "not-a-url"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": response.InputErrMessage,
				"details": []any{"URL is invalid (url)"},
			},
		},
		{
			name: "error - invalid input (description too long)",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			inputBody: map[string]any{
				"description": strings.Repeat("a", 256),
				"url":         "https://example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": response.InputErrMessage,
				"details": []any{"Description is invalid (lte)"},
			},
		},
		{
			name: "error - invalid input (URL too long)",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			inputBody: map[string]any{
				"description": "My Bookmark",
				"url":         testBookmarkLongURL,
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": response.InputErrMessage,
				"details": []any{"URL is invalid (lte)"},
			},
		},
		{
			name: "error - service failure",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			inputBody: map[string]any{"description": "My Bookmark", "url": "https://example.com"},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("CreateBookmark",
					ctx,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil, errors.New("service error"))
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

			// Create test context with JWT claims and JSON body
			testCtx := handlertest.NewTestContext(http.MethodPost, "/v1/bookmarks").
				WithJWTClaims(tc.jwtClaims).
				WithJSONBody(tc.inputBody)

			// Setup mock service
			svcMock := tc.setupMockSvc(t, testCtx.Ctx)

			// Create handler with mock service
			handler := NewHandler(svcMock)

			// Call the handler
			handler.CreateBookmark(testCtx.Ctx)

			// Assert response
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
