package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	repoMocks "github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkSvc_GetUrlByCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		inputCode   int64
		setupMock   func(mockRepo *repoMocks.Repository, ctx context.Context)
		expectedUrl string
		expectedErr error
	}{
		{
			name:      "Success - returns URL",
			inputCode: 42,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("GetBookmarkByCode", ctx, int64(42)).
					Return(&model.Bookmark{
						URL: "https://example.com/found",
					}, nil)
			},
			expectedUrl: "https://example.com/found",
		},
		{
			name:      "Not found - returns ErrBookmarkNotFound",
			inputCode: 999,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("GetBookmarkByCode", ctx, int64(999)).
					Return(nil, gorm.ErrRecordNotFound)
			},
			expectedUrl: "",
			expectedErr: ErrBookmarkNotFound,
		},
		{
			name:      "DB error - propagates error",
			inputCode: 42,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("GetBookmarkByCode", ctx, int64(42)).
					Return(nil, errors.New("db connection error"))
			},
			expectedUrl: "",
			expectedErr: errors.New("db connection error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := repoMocks.NewRepository(t)
			tc.setupMock(mockRepo, ctx)

			svc := NewBookmarkSvc(mockRepo)
			url, err := svc.GetUrlByCode(ctx, tc.inputCode)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				assert.Empty(t, url)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUrl, url)
		})
	}
}
