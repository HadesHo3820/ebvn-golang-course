package bookmark_test

import (
	"context"
	"errors"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/repository/bookmark/mocks"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service/bookmark"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkSvc_DeleteBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		bookmarkID  string
		userID      string
		setupMock   func(m *mocks.Repository)
		expectedErr error
	}{
		{
			name:       "Success",
			bookmarkID: "bookmark-uuid-123",
			userID:     "user-uuid-456",
			setupMock: func(m *mocks.Repository) {
				m.On("DeleteBookmark", mock.Anything, "bookmark-uuid-123", "user-uuid-456").
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:       "Error - Repository Not Found",
			bookmarkID: "nonexistent-id",
			userID:     "user-uuid-456",
			setupMock: func(m *mocks.Repository) {
				m.On("DeleteBookmark", mock.Anything, "nonexistent-id", "user-uuid-456").
					Return(dbutils.ErrNotFoundType)
			},
			expectedErr: dbutils.ErrNotFoundType,
		},
		{
			name:       "Error - Repository Error",
			bookmarkID: "bookmark-uuid-123",
			userID:     "user-uuid-456",
			setupMock: func(m *mocks.Repository) {
				m.On("DeleteBookmark", mock.Anything, "bookmark-uuid-123", "user-uuid-456").
					Return(errors.New("database connection error"))
			},
			expectedErr: errors.New("database connection error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock
			mockRepo := mocks.NewRepository(t)
			tc.setupMock(mockRepo)

			// Create service with mock
			svc := bookmark.NewBookmarkSvc(mockRepo, nil)

			// Execute
			err := svc.DeleteBookmark(context.Background(), tc.bookmarkID, tc.userID)

			// Assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				if errors.Is(tc.expectedErr, dbutils.ErrNotFoundType) {
					assert.ErrorIs(t, err, dbutils.ErrNotFoundType)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
