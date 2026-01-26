package bookmark

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/pagination"
)

// GetBookmarks retrieves a paginated list of bookmarks for the specified user.
// It handles the calculation of offset/limit from the request, fetches data from the repository,
// and constructs the final paginated response with metadata.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: The ID of the owner
//   - req: Pointer to Pagination request with Page and Limit
//
// Returns:
//   - *pagination.Response: Standard paginated response wrapper
//   - error: Database or internal error
func (s *BookmarkSvc) GetBookmarks(ctx context.Context, userID string, req *pagination.Request) (*pagination.Response[*model.Bookmark], error) {
	limit := req.GetLimit()
	offset := req.GetOffset()

	bookmarks, total, err := s.repo.GetBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	meta := pagination.CalculateMetadata(total, req.Page, limit)

	return &pagination.Response[*model.Bookmark]{
		Data:     bookmarks,
		Metadata: meta,
	}, nil
}
