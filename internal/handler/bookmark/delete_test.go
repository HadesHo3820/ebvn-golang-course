package bookmark

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"

	serviceMocks "github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark/mocks"
	handlertest "github.com/HadesHo3820/ebvn-golang-course/internal/test/handler"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
)

const testBookmarkIDDelete = "550e8400-e29b-41d4-a716-446655440000"

func TestBookmarkHandler_DeleteBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		jwtClaims      jwt.MapClaims
		uriParams      map[string]string
		setupMockSvc   func(t *testing.T, ctx context.Context) *serviceMocks.Service
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success - delete bookmark",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			uriParams: map[string]string{"id": testBookmarkIDDelete},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("DeleteBookmark",
					ctx,
					testBookmarkIDDelete,
					testUserID,
				).Return(nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Bookmark deleted successfully",
			},
		},
		{
			name:      "error - missing JWT claims",
			jwtClaims: nil,
			uriParams: map[string]string{"id": testBookmarkIDDelete},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"message": "Invalid token",
			},
		},
		{
			name: "error - invalid UUID",
			jwtClaims: jwt.MapClaims{
				"sub": testUserID,
			},
			uriParams: map[string]string{"id": "not-a-valid-uuid"},
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
			uriParams: map[string]string{"id": testBookmarkIDDelete},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("DeleteBookmark",
					ctx,
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
			uriParams: map[string]string{"id": testBookmarkIDDelete},
			setupMockSvc: func(t *testing.T, ctx context.Context) *serviceMocks.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("DeleteBookmark",
					ctx,
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

			// Create test context with JWT claims and URI params
			testCtx := handlertest.NewTestContext(http.MethodDelete, "/v1/bookmarks/:id").
				WithJWTClaims(tc.jwtClaims).
				WithURIParams(tc.uriParams)

			// Setup mock service
			svcMock := tc.setupMockSvc(t, testCtx.Ctx)

			// Create handler with mock service
			handler := NewHandler(svcMock)

			// Call the handler
			handler.DeleteBookmark(testCtx.Ctx)

			// Assert response
			handlertest.AssertJSONResponse(t, testCtx.Recorder, tc.expectedStatus, tc.expectedBody)
		})
	}
}
