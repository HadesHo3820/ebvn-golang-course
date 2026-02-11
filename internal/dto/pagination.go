package dto

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

// Sanitize normalizes the Page and Limit values to ensure they are within valid ranges.
// This method mutates the struct in place and should be called before using Page/Limit
// in any operations (e.g., cache key generation, database queries).
//
// Rules:
//   - Page < 1 → set to 1
//   - Limit < 1 → set to DefaultLimit (10)
//   - Limit > MaxLimit → set to MaxLimit (100)
func (r *Request) Sanitize() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = DefaultLimit
	}
	if r.Limit > MaxLimit {
		r.Limit = MaxLimit
	}
}

// GetOffset calculates the database offset based on Page and Limit.
// It automatically sanitizes the values before calculation.
func (r *Request) GetOffset() int {
	r.Sanitize()
	return (r.Page - 1) * r.Limit
}

// GetLimit returns the sanitized limit.
// It automatically sanitizes the values before returning.
func (r *Request) GetLimit() int {
	r.Sanitize()
	return r.Limit
}

// Metadata contains pagination details to be returned in the API response.
type Metadata struct {
	CurrentPage  int   `json:"current_page" example:"1"`
	PageSize     int   `json:"page_size" example:"10"`
	FirstPage    int   `json:"first_page" example:"1"`
	LastPage     int   `json:"last_page" example:"1"`
	TotalRecords int64 `json:"total_records" example:"1"`
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
