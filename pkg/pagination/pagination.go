package pagination

import "math"

// DefaultLimit is the default number of items per page.
const DefaultLimit = 10

// MaxLimit be used to limit the number of items per page.
const MaxLimit = 100

// Request represents the standard pagination query parameters.
//
// The 'json' tags are included to allow this struct to be embedded in
// larger request structs used for POST requests (e.g., search with filters),
// where pagination parameters are sent as part of the JSON body.
type Request struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}

// GetOffset calculates the database offset based on Page and Limit.
// It also applies default values if Page or Limit are invalid.
func (r *Request) GetOffset() int {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = DefaultLimit
	}
	if r.Limit > MaxLimit {
		r.Limit = MaxLimit
	}
	return (r.Page - 1) * r.Limit
}

// GetLimit returns the sanitized limit.
func (r *Request) GetLimit() int {
	if r.Limit < 1 {
		r.Limit = DefaultLimit
	}
	if r.Limit > MaxLimit {
		r.Limit = MaxLimit
	}
	return r.Limit
}

// Metadata contains pagination details to be returned in the API response.
type Metadata struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	FirstPage    int   `json:"first_page"`
	LastPage     int   `json:"last_page"`
	TotalRecords int64 `json:"total_records"`
}

// CalculateMetadata constructs the Metadata struct based on total records and current page settings.
func CalculateMetadata(total int64, page int, limit int) Metadata {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = DefaultLimit
	}

	// Calculate last page:
	// 1. Divide total by limit and round UP (Ceil) to cover partial pages.
	// 2. Ensure value is at least 1 (max) so we never return "Page 0 of 0".
	lastPage := max(int(math.Ceil(float64(total)/float64(limit))), 1)

	return Metadata{
		CurrentPage:  page,
		PageSize:     limit,
		FirstPage:    1,
		LastPage:     lastPage,
		TotalRecords: total,
	}
}

// Response is a generic wrapper for paginated API responses.
type Response[T any] struct {
	Data     []T      `json:"data"`
	Metadata Metadata `json:"metadata"`
}
