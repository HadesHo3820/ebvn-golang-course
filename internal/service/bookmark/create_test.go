package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	repoMocks "github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testBookmarkDesc = "Test Bookmark"
	testBookmarkURL  = "https://example.com"
	testUserID       = "user-123"
	testCode         = "abc123456"
	testBookmarkID   = "bookmark-1"
)

func TestBookmarkSvc_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		inputDescription string
		inputURL         string
		inputUserID      string
		setupMock        func(mockRepo *repoMocks.Repository, mockCodeGen *mocks.KeyGenerator, ctx context.Context)
		expectedErr      error
		expectedOutput   *model.Bookmark
	}{
		{
			name:             "Success",
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			inputUserID:      testUserID,
			setupMock: func(mockRepo *repoMocks.Repository, mockCodeGen *mocks.KeyGenerator, ctx context.Context) {
				mockCodeGen.On("GenerateCode", 9).Return(testCode, nil)
				mockRepo.On("CreateBookmark", ctx, mock.Anything).
					Return(&model.Bookmark{
						Base:        model.Base{ID: testBookmarkID},
						Description: testBookmarkDesc,
						URL:         testBookmarkURL,
						Code:        testCode,
						UserID:      testUserID,
					}, nil)
			},
			expectedOutput: &model.Bookmark{
				Base:        model.Base{ID: testBookmarkID},
				Description: testBookmarkDesc,
				URL:         testBookmarkURL,
				Code:        testCode,
				UserID:      testUserID,
			},
		},
		{
			name:             "Error - Key Generation Failed",
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			inputUserID:      testUserID,
			setupMock: func(mockRepo *repoMocks.Repository, mockCodeGen *mocks.KeyGenerator, ctx context.Context) {
				mockCodeGen.On("GenerateCode", 9).Return("", errors.New("code gen error"))
			},
			expectedErr: errors.New("code gen error"),
		},
		{
			name:             "Error - Repository Creation Failed",
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			inputUserID:      testUserID,
			setupMock: func(mockRepo *repoMocks.Repository, mockCodeGen *mocks.KeyGenerator, ctx context.Context) {
				mockCodeGen.On("GenerateCode", 9).Return(testCode, nil)
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
			mockCodeGen := mocks.NewKeyGenerator(t)
			tc.setupMock(mockRepo, mockCodeGen, ctx)

			// Create service
			svc := NewBookmarkSvc(mockRepo, mockCodeGen)

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
