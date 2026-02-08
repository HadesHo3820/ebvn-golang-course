package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	repoMocks "github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkSvc_GetBookmarks(t *testing.T) {
	t.Parallel()

	testPage := 1
	testLimit := 10
	testOffset := 0
	testTotal := int64(20)

	testCases := []struct {
		name           string
		inputUserID    string
		inputReq       *dto.Request
		setupMock      func(mockRepo *repoMocks.Repository, ctx context.Context)
		expectedErr    error
		expectedOutput *dto.Response[*model.Bookmark]
	}{
		{
			name:        "Success",
			inputUserID: testUserID,
			inputReq: &dto.Request{
				Page:  testPage,
				Limit: testLimit,
			},
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("GetBookmarks", ctx, testUserID, testLimit, testOffset).
					Return([]*model.Bookmark{
						{Base: model.Base{ID: "bm-1"}},
						{Base: model.Base{ID: "bm-2"}},
					}, nil)
				mockRepo.On("GetBookmarksCount", ctx, testUserID).
					Return(testTotal, nil)
			},
			expectedOutput: &dto.Response[*model.Bookmark]{
				Data: []*model.Bookmark{
					{Base: model.Base{ID: "bm-1"}},
					{Base: model.Base{ID: "bm-2"}},
				},
				Metadata: dto.Metadata{
					CurrentPage:  testPage,
					PageSize:     testLimit,
					FirstPage:    1,
					LastPage:     2, // 20 / 10 = 2
					TotalRecords: testTotal,
				},
			},
		},
		{
			name:        "Error - GetBookmarks Failed",
			inputUserID: testUserID,
			inputReq: &dto.Request{
				Page:  testPage,
				Limit: testLimit,
			},
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("GetBookmarks", ctx, testUserID, testLimit, testOffset).
					Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
		{
			name:        "Error - GetBookmarksCount Failed",
			inputUserID: testUserID,
			inputReq: &dto.Request{
				Page:  testPage,
				Limit: testLimit,
			},
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("GetBookmarks", ctx, testUserID, testLimit, testOffset).
					Return([]*model.Bookmark{}, nil)
				mockRepo.On("GetBookmarksCount", ctx, testUserID).
					Return(int64(0), errors.New("count error"))
			},
			expectedErr: errors.New("count error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup mocks
			mockRepo := repoMocks.NewRepository(t)
			mockCodeGen := mocks.NewKeyGenerator(t)
			tc.setupMock(mockRepo, ctx)

			// Create service
			svc := NewBookmarkSvc(mockRepo, mockCodeGen)

			// Execute
			got, err := svc.GetBookmarks(ctx, tc.inputUserID, tc.inputReq)

			// Assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, got)
		})
	}
}
