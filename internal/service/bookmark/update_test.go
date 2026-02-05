package bookmark

import (
	"context"
	"errors"
	"testing"

	repoMocks "github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkSvc_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		inputBookmarkID  string
		inputUserID      string
		inputDescription string
		inputURL         string
		setupMock        func(mockRepo *repoMocks.Repository, ctx context.Context)
		expectedErr      error
	}{
		{
			name:             "Success",
			inputBookmarkID:  testBookmarkID,
			inputUserID:      testUserID,
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("UpdateBookmark", ctx, testBookmarkID, testUserID, testBookmarkDesc, testBookmarkURL).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:             "Error - Repository Not Found",
			inputBookmarkID:  testBookmarkID,
			inputUserID:      testUserID,
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("UpdateBookmark", ctx, testBookmarkID, testUserID, testBookmarkDesc, testBookmarkURL).
					Return(dbutils.ErrNotFoundType)
			},
			expectedErr: dbutils.ErrNotFoundType,
		},
		{
			name:             "Error - Repository Error",
			inputBookmarkID:  testBookmarkID,
			inputUserID:      testUserID,
			inputDescription: testBookmarkDesc,
			inputURL:         testBookmarkURL,
			setupMock: func(mockRepo *repoMocks.Repository, ctx context.Context) {
				mockRepo.On("UpdateBookmark", ctx, testBookmarkID, testUserID, testBookmarkDesc, testBookmarkURL).
					Return(errors.New("db error"))
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
			tc.setupMock(mockRepo, ctx)

			// Create service
			svc := NewBookmarkSvc(mockRepo, mockCodeGen)

			// Execute
			err := svc.UpdateBookmark(ctx, tc.inputBookmarkID, tc.inputUserID, tc.inputDescription, tc.inputURL)

			// Assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
