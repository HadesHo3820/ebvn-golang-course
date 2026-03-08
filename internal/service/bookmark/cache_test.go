package bookmark_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	mock_cache "github.com/HadesHo3820/ebvn-golang-course/internal/repository/cache/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreateBookmark ensures bookmarks are created and cache invalidated.
func TestCreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		description   string
		url           string
		userID        string
		setupMocks    func(context.Context, *mocks.Service, *mock_cache.DB)
		expectedError error
		checkResult   func(*testing.T, *model.Bookmark)
	}{
		{
			name:        "success - bookmark created and cache invalidated",
			description: "Test bookmark",
			url:         "https://example.com",
			userID:      "user-123",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				expectedBookmark := &model.Bookmark{
					Base: model.Base{
						ID: "bookmark-456",
					},
					Description: "Test bookmark",
					URL:         "https://example.com",
					UserID:      "user-123",
				}
				mockService.On("CreateBookmark", ctx, "Test bookmark", "https://example.com", "user-123").
					Return(expectedBookmark, nil)
				mockCache.On("DeleteCacheData", mock.Anything, "get_bookmarks_user-123").
					Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, bm *model.Bookmark) {
				assert.NotNil(t, bm)
				assert.Equal(t, "bookmark-456", bm.ID)
				assert.Equal(t, "Test bookmark", bm.Description)
			},
		},
		{
			name:        "db error - creation fails",
			description: "Test bookmark",
			url:         "https://example.com",
			userID:      "user-123",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockService.On("CreateBookmark", ctx, "Test bookmark", "https://example.com", "user-123").
					Return(nil, errors.New("db error"))
				// Cache should not be called when DB fails
			},
			expectedError: errors.New("db error"),
			checkResult: func(t *testing.T, bm *model.Bookmark) {
				assert.Nil(t, bm)
			},
		},
		{
			name:        "cache invalidation fails - bookmark still created",
			description: "Test bookmark",
			url:         "https://example.com",
			userID:      "user-123",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				expectedBookmark := &model.Bookmark{
					Base: model.Base{
						ID: "bookmark-456",
					},
					Description: "Test bookmark",
					URL:         "https://example.com",
					UserID:      "user-123",
				}
				mockService.On("CreateBookmark", ctx, "Test bookmark", "https://example.com", "user-123").
					Return(expectedBookmark, nil)
				mockCache.On("DeleteCacheData", mock.Anything, "get_bookmarks_user-123").
					Return(errors.New("cache error"))
			},
			expectedError: nil, // Cache errors should not fail the operation
			checkResult: func(t *testing.T, bm *model.Bookmark) {
				assert.NotNil(t, bm)
				assert.Equal(t, "bookmark-456", bm.ID)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// Setup
			mockService := mocks.NewService(t)
			mockCache := mock_cache.NewDB(t)
			tc.setupMocks(ctx, mockService, mockCache)

			svc := bookmark.NewServiceWithCache(mockService, mockCache)

			// Execute
			result, err := svc.CreateBookmark(ctx, tc.description, tc.url, tc.userID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			tc.checkResult(t, result)
		})
	}
}

// TestGetBookmarks validates bookmark retrieval logic.
// Specifically, it uses mock.Anything for context in SetCacheData expectations
// to handle the context.WithTimeout used in the implementation.
func TestGetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		userID        string
		request       *dto.Request
		setupMocks    func(context.Context, *mocks.Service, *mock_cache.DB)
		expectedError error
		checkResult   func(*testing.T, *dto.Response[*model.Bookmark])
	}{
		{
			name:    "nil request - returns error",
			userID:  "user-123",
			request: nil,
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				// No mocks needed, validation happens first
			},
			expectedError: bookmark.ErrNilRequest,
			checkResult: func(t *testing.T, resp *dto.Response[*model.Bookmark]) {
				assert.Nil(t, resp)
			},
		},
		{
			name:   "cache hit - returns cached data",
			userID: "user-123",
			request: &dto.Request{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				cachedResponse := &dto.Response[*model.Bookmark]{
					Data: []*model.Bookmark{
						{Base: model.Base{ID: "cached-1"}, Description: "Cached bookmark"},
					},
					Metadata: dto.Metadata{
						CurrentPage:  1,
						PageSize:     10,
						TotalRecords: 1,
					},
				}
				cachedBytes, _ := json.Marshal(cachedResponse)
				mockCache.On("GetCacheData", ctx, "get_bookmarks_user-123", "1_10").
					Return(cachedBytes, nil)
				// Service should not be called on cache hit
			},
			expectedError: nil,
			checkResult: func(t *testing.T, resp *dto.Response[*model.Bookmark]) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, "cached-1", resp.Data[0].ID)
			},
		},
		{
			name:   "cache hit but unmarshal fails - falls back to DB",
			userID: "user-123",
			request: &dto.Request{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockCache.On("GetCacheData", ctx, "get_bookmarks_user-123", "1_10").
					Return([]byte("invalid json"), nil)

				dbResponse := &dto.Response[*model.Bookmark]{
					Data: []*model.Bookmark{
						{Base: model.Base{ID: "db-1"}, Description: "DB bookmark"},
					},
					Metadata: dto.Metadata{
						CurrentPage:  1,
						PageSize:     10,
						TotalRecords: 1,
					},
				}
				mockService.On("GetBookmarks", ctx, "user-123", mock.Anything).
					Return(dbResponse, nil)

				// Cache population after fetching from DB
				mockCache.On("SetCacheData", mock.Anything, "get_bookmarks_user-123", "1_10", mock.Anything, 24*time.Hour).
					Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, resp *dto.Response[*model.Bookmark]) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, "db-1", resp.Data[0].ID)
			},
		},
		{
			name:   "cache miss - fetches from DB and caches",
			userID: "user-123",
			request: &dto.Request{
				Page:  2,
				Limit: 20,
			},
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockCache.On("GetCacheData", ctx, "get_bookmarks_user-123", "2_20").
					Return(nil, errors.New("cache miss"))

				dbResponse := &dto.Response[*model.Bookmark]{
					Data: []*model.Bookmark{
						{Base: model.Base{ID: "db-2"}, Description: "Fresh DB bookmark"},
					},
					Metadata: dto.Metadata{
						CurrentPage:  2,
						PageSize:     20,
						TotalRecords: 1,
					},
				}
				mockService.On("GetBookmarks", ctx, "user-123", mock.Anything).
					Return(dbResponse, nil)

				mockCache.On("SetCacheData", mock.Anything, "get_bookmarks_user-123", "2_20", mock.Anything, 24*time.Hour).
					Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, resp *dto.Response[*model.Bookmark]) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, "db-2", resp.Data[0].ID)
			},
		},
		{
			name:   "cache miss and DB error",
			userID: "user-123",
			request: &dto.Request{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockCache.On("GetCacheData", ctx, "get_bookmarks_user-123", "1_10").
					Return(nil, errors.New("cache miss"))

				mockService.On("GetBookmarks", ctx, "user-123", mock.Anything).
					Return(nil, errors.New("database connection failed"))
			},
			expectedError: errors.New("database connection failed"),
			checkResult: func(t *testing.T, resp *dto.Response[*model.Bookmark]) {
				assert.Nil(t, resp)
			},
		},
		{
			name:   "cache population fails after DB fetch - still returns data",
			userID: "user-123",
			request: &dto.Request{
				Page:  1,
				Limit: 10,
			},
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockCache.On("GetCacheData", ctx, "get_bookmarks_user-123", "1_10").
					Return(nil, errors.New("cache miss"))

				dbResponse := &dto.Response[*model.Bookmark]{
					Data: []*model.Bookmark{
						{Base: model.Base{ID: "db-3"}, Description: "DB bookmark"},
					},
					Metadata: dto.Metadata{
						CurrentPage:  1,
						PageSize:     10,
						TotalRecords: 1,
					},
				}
				mockService.On("GetBookmarks", ctx, "user-123", mock.Anything).
					Return(dbResponse, nil)

				mockCache.On("SetCacheData", mock.Anything, "get_bookmarks_user-123", "1_10", mock.Anything, 24*time.Hour).
					Return(errors.New("cache write failed"))
			},
			expectedError: nil, // Cache errors should not fail the operation
			checkResult: func(t *testing.T, resp *dto.Response[*model.Bookmark]) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, "db-3", resp.Data[0].ID)
			},
		},
		{
			name:   "sanitization applied - zero page becomes 1",
			userID: "user-123",
			request: &dto.Request{
				Page:  0, // Should be sanitized to 1
				Limit: 10,
			},
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockCache.On("GetCacheData", ctx, "get_bookmarks_user-123", "1_10").
					Return(nil, errors.New("cache miss"))

				dbResponse := &dto.Response[*model.Bookmark]{
					Data:     []*model.Bookmark{{Base: model.Base{ID: "test"}}},
					Metadata: dto.Metadata{CurrentPage: 1, PageSize: 10, TotalRecords: 1},
				}
				mockService.On("GetBookmarks", ctx, "user-123", mock.Anything).
					Return(dbResponse, nil)

				mockCache.On("SetCacheData", mock.Anything, "get_bookmarks_user-123", "1_10", mock.Anything, 24*time.Hour).
					Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, resp *dto.Response[*model.Bookmark]) {
				assert.NotNil(t, resp)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup
			mockService := mocks.NewService(t)
			mockCache := mock_cache.NewDB(t)
			tc.setupMocks(ctx, mockService, mockCache)

			svc := bookmark.NewServiceWithCache(mockService, mockCache)

			// Execute
			result, err := svc.GetBookmarks(ctx, tc.userID, tc.request)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			tc.checkResult(t, result)
		})
	}
}

// TestUpdateBookmark validates bookmark updates and cache invalidation.
func TestUpdateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		bookmarkID    string
		userID        string
		description   string
		url           string
		setupMocks    func(context.Context, *mocks.Service, *mock_cache.DB)
		expectedError error
	}{
		{
			name:        "success - bookmark updated and cache invalidated",
			bookmarkID:  "bookmark-123",
			userID:      "user-456",
			description: "Updated description",
			url:         "https://updated.com",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockService.On("UpdateBookmark", ctx, "bookmark-123", "user-456", "Updated description", "https://updated.com").
					Return(nil)
				mockCache.On("DeleteCacheData", mock.Anything, "get_bookmarks_user-456").
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "db error - update fails",
			bookmarkID:  "bookmark-123",
			userID:      "user-456",
			description: "Updated description",
			url:         "https://updated.com",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockService.On("UpdateBookmark", ctx, "bookmark-123", "user-456", "Updated description", "https://updated.com").
					Return(errors.New("not found"))
				// Cache should not be called when DB fails
			},
			expectedError: errors.New("not found"),
		},
		{
			name:        "cache invalidation fails - update still succeeds",
			bookmarkID:  "bookmark-123",
			userID:      "user-456",
			description: "Updated description",
			url:         "https://updated.com",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockService.On("UpdateBookmark", ctx, "bookmark-123", "user-456", "Updated description", "https://updated.com").
					Return(nil)
				mockCache.On("DeleteCacheData", mock.Anything, "get_bookmarks_user-456").
					Return(errors.New("cache error"))
			},
			expectedError: nil, // Cache errors should not fail the operation
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup
			mockService := mocks.NewService(t)
			mockCache := mock_cache.NewDB(t)
			tc.setupMocks(ctx, mockService, mockCache)

			svc := bookmark.NewServiceWithCache(mockService, mockCache)

			// Execute
			err := svc.UpdateBookmark(ctx, tc.bookmarkID, tc.userID, tc.description, tc.url)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDeleteBookmark validates bookmark deletion and cache invalidation.
func TestDeleteBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		bookmarkID    string
		userID        string
		setupMocks    func(context.Context, *mocks.Service, *mock_cache.DB)
		expectedError error
	}{
		{
			name:       "success - bookmark deleted and cache invalidated",
			bookmarkID: "bookmark-123",
			userID:     "user-456",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockService.On("DeleteBookmark", ctx, "bookmark-123", "user-456").
					Return(nil)
				mockCache.On("DeleteCacheData", mock.Anything, "get_bookmarks_user-456").
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "db error - deletion fails",
			bookmarkID: "bookmark-123",
			userID:     "user-456",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockService.On("DeleteBookmark", ctx, "bookmark-123", "user-456").
					Return(errors.New("not found"))
				// Cache should not be called when DB fails
			},
			expectedError: errors.New("not found"),
		},
		{
			name:       "cache invalidation fails - deletion still succeeds",
			bookmarkID: "bookmark-123",
			userID:     "user-456",
			setupMocks: func(ctx context.Context, mockService *mocks.Service, mockCache *mock_cache.DB) {
				mockService.On("DeleteBookmark", ctx, "bookmark-123", "user-456").
					Return(nil)
				mockCache.On("DeleteCacheData", mock.Anything, "get_bookmarks_user-456").
					Return(errors.New("cache error"))
			},
			expectedError: nil, // Cache errors should not fail the operation
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup
			mockService := mocks.NewService(t)
			mockCache := mock_cache.NewDB(t)
			tc.setupMocks(ctx, mockService, mockCache)

			svc := bookmark.NewServiceWithCache(mockService, mockCache)

			// Execute
			err := svc.DeleteBookmark(ctx, tc.bookmarkID, tc.userID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
