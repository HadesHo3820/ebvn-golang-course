package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	repoMocks "github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testBookmarkDesc                 = "Test Bookmark"
	testBookmarkURL                  = "https://example.com"
	testUserID                       = "user-123"
	testBookmarkSequenceCodeID int64 = 42
	testBookmarkID                   = "bookmark-1"
)

func TestBookmarkSvc_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		inputDescription string
		inputURL         string
		inputUserID      string
		setupMock        func(mockRepo *repoMocks.Repository, ctx context.Context)
		expectedErr      error
		expectedOutput   *model.Bookmark
	}{
		{
			name:             "Success",
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			inputUserID:      testUserID,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("CreateBookmark", ctx, mock.Anything).
					Return(func() *model.Bookmark {
						encoded := "G" // base62.Encode(42)
						return &model.Bookmark{
							Base:                model.Base{ID: testBookmarkID},
							Description:         testBookmarkDesc,
							URL:                 testBookmarkURL,
							SequenceCodeID:      testBookmarkSequenceCodeID,
							EncodedBookmarkCode: &encoded,
							UserID:              testUserID,
						}
					}(), nil)
			},
			expectedOutput: func() *model.Bookmark {
				encoded := "G" // base62.Encode(42)
				return &model.Bookmark{
					Base:                model.Base{ID: testBookmarkID},
					Description:         testBookmarkDesc,
					URL:                 testBookmarkURL,
					SequenceCodeID:      testBookmarkSequenceCodeID,
					EncodedBookmarkCode: &encoded,
					UserID:              testUserID,
				}
			}(),
		},
		{
			name:             "Error - Repository Creation Failed",
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			inputUserID:      testUserID,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("CreateBookmark", ctx, mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// Setup mocks
			mockRepo := repoMocks.NewRepository(t)
			tc.setupMock(mockRepo, ctx)

			// Create service (no codeGen needed — codes are auto-incremented by DB)
			svc := NewBookmarkSvc(mockRepo)

			// Execute
			got, err := svc.CreateBookmark(ctx, tc.inputDescription, tc.inputURL, tc.inputUserID)

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
