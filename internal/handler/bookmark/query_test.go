package bookmark

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	serviceMocks "github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark/mocks"
	handlertest "github.com/HadesHo3820/ebvn-golang-course/internal/test/handler"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/pagination"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

const (
	testQueryBookmarkDesc = "Bookmark 1"
	testQueryBookmarkURL  = "https://example.com/1"
	testQueryBookmarkCode = "test-code"
)

func TestBookmarkHandler_GetBookmarks(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		jwtClaims      jwt.MapClaims
		queryParams    string // e.g., "?page=1&limit=10"
		setupMockSvc   func(t *testing.T, ctx context.Context) *serviceMocks.Service
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success - get bookmarks with defaults",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			queryParams: "",
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("GetBookmarks",
					mock.Anything,
					testUserID,
					mock.MatchedBy(func(req pagination.Request) bool {
						return req.Page == 0 && req.Limit == 0 // Defaults before validation/sanitization in service/repo layer
					}),
				).Return(&pagination.Response[*model.Bookmark]{
					Data: []*model.Bookmark{
						{
							Base: model.Base{
								ID:        "bm-1",
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
							},
							Description: testQueryBookmarkDesc,
							URL:         testQueryBookmarkURL,
							Code:        testQueryBookmarkCode,
							UserID:      testUserID,
						},
					},
					Metadata: pagination.Metadata{
						CurrentPage:  1,
						PageSize:     10,
						TotalRecords: 1,
						FirstPage:    1,
						LastPage:     1,
					},
				}, nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"data": []any{
					map[string]any{
						"id":          "bm-1",
						"description": testQueryBookmarkDesc,
						"url":         testQueryBookmarkURL,
						"code":        testQueryBookmarkCode,
						"user_id":     testUserID,
						"created_at":  fixedTime.Format(time.RFC3339Nano),
						"updated_at":  fixedTime.Format(time.RFC3339Nano),
					},
				},
				"metadata": map[string]any{
					"current_page":  float64(1),
					"page_size":     float64(10),
					"total_records": float64(1),
					"first_page":    float64(1),
					"last_page":     float64(1),
				},
			},
		},
		{
			name: "success - get bookmarks with custom pagination",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			queryParams: "?page=2&limit=5",
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("GetBookmarks",
					mock.Anything,
					testUserID,
					pagination.Request{Page: 2, Limit: 5},
				).Return(&pagination.Response[*model.Bookmark]{
					Data: []*model.Bookmark{},
					Metadata: pagination.Metadata{
						CurrentPage:  2,
						PageSize:     5,
						TotalRecords: 20,
					},
				}, nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"data": []any{},
				"metadata": map[string]any{
					"current_page":  float64(2),
					"page_size":     float64(5),
					"total_records": float64(20),
					"first_page":    float64(0), // Default zero value if not set in response mock
					"last_page":     float64(0),
				},
			},
		},
		{
			name:      "error - missing JWT claims",
			jwtClaims: nil,
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"message": "Invalid token",
			},
		},
		{
			name: "error - service failure",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			queryParams: "",
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("GetBookmarks",
					mock.Anything,
					testUserID,
					mock.Anything,
				).Return(nil, errors.New("db error"))
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

			// Create test context with JWT claims and path with query params
			path := "/v1/bookmarks"
			if tc.queryParams != "" {
				path += tc.queryParams
			}

			testCtx := handlertest.NewTestContext(http.MethodGet, path).
				WithJWTClaims(tc.jwtClaims)

			// Setup mock service
			svcMock := tc.setupMockSvc(t, testCtx.Ctx)

			// Create handler with mock service
			handler := NewHandler(svcMock)

			// Call the handler
			handler.GetBookmarks(testCtx.Ctx)

			// Assert response
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
