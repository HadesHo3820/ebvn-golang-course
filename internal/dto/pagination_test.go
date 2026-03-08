package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_GetOffset(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		request        *Request
		expectedOffset int
		expectedPage   int
		expectedLimit  int
	}{
		{
			name:           "valid page and limit",
			request:        &Request{Page: 2, Limit: 10},
			expectedOffset: 10,
			expectedPage:   2,
			expectedLimit:  10,
		},
		{
			name:           "first page",
			request:        &Request{Page: 1, Limit: 10},
			expectedOffset: 0,
			expectedPage:   1,
			expectedLimit:  10,
		},
		{
			name:           "page less than 1 - defaults to 1",
			request:        &Request{Page: 0, Limit: 10},
			expectedOffset: 0,
			expectedPage:   1,
			expectedLimit:  10,
		},
		{
			name:           "negative page - defaults to 1",
			request:        &Request{Page: -5, Limit: 10},
			expectedOffset: 0,
			expectedPage:   1,
			expectedLimit:  10,
		},
		{
			name:           "limit less than 1 - defaults to DefaultLimit",
			request:        &Request{Page: 1, Limit: 0},
			expectedOffset: 0,
			expectedPage:   1,
			expectedLimit:  DefaultLimit,
		},
		{
			name:           "limit exceeds MaxLimit - caps at MaxLimit",
			request:        &Request{Page: 1, Limit: 200},
			expectedOffset: 0,
			expectedPage:   1,
			expectedLimit:  MaxLimit,
		},
		{
			name:           "large page number",
			request:        &Request{Page: 100, Limit: 25},
			expectedOffset: 2475, // (100-1) * 25
			expectedPage:   100,
			expectedLimit:  25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			offset := tc.request.GetOffset()

			assert.Equal(t, tc.expectedOffset, offset)
			assert.Equal(t, tc.expectedPage, tc.request.Page)
			assert.Equal(t, tc.expectedLimit, tc.request.Limit)
		})
	}
}

func TestRequest_GetLimit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		request       *Request
		expectedLimit int
	}{
		{
			name:          "valid limit",
			request:       &Request{Limit: 20},
			expectedLimit: 20,
		},
		{
			name:          "limit less than 1 - defaults to DefaultLimit",
			request:       &Request{Limit: 0},
			expectedLimit: DefaultLimit,
		},
		{
			name:          "negative limit - defaults to DefaultLimit",
			request:       &Request{Limit: -10},
			expectedLimit: DefaultLimit,
		},
		{
			name:          "limit exceeds MaxLimit - caps at MaxLimit",
			request:       &Request{Limit: 500},
			expectedLimit: MaxLimit,
		},
		{
			name:          "limit at MaxLimit - returns MaxLimit",
			request:       &Request{Limit: MaxLimit},
			expectedLimit: MaxLimit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			limit := tc.request.GetLimit()

			assert.Equal(t, tc.expectedLimit, limit)
		})
	}
}

func TestCalculateMetadata(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		total    int64
		page     int
		limit    int
		expected Metadata
	}{
		{
			name:  "standard case",
			total: 100,
			page:  1,
			limit: 10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     10,
				TotalRecords: 100,
			},
		},
		{
			name:  "partial last page",
			total: 25,
			page:  2,
			limit: 10,
			expected: Metadata{
				CurrentPage:  2,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     3, // 25/10 = 2.5, ceil = 3
				TotalRecords: 25,
			},
		},
		{
			name:  "zero total records",
			total: 0,
			page:  1,
			limit: 10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     1, // min is 1
				TotalRecords: 0,
			},
		},
		{
			name:  "page less than 1 - defaults to 1",
			total: 50,
			page:  0,
			limit: 10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     5,
				TotalRecords: 50,
			},
		},
		{
			name:  "limit less than 1 - defaults to DefaultLimit",
			total: 50,
			page:  1,
			limit: 0,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     DefaultLimit,
				FirstPage:    1,
				LastPage:     5, // 50/10 = 5
				TotalRecords: 50,
			},
		},
		{
			name:  "single item",
			total: 1,
			page:  1,
			limit: 10,
			expected: Metadata{
				CurrentPage:  1,
				PageSize:     10,
				FirstPage:    1,
				LastPage:     1,
				TotalRecords: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := CalculateMetadata(tc.total, tc.page, tc.limit)

			assert.Equal(t, tc.expected, result)
		})
	}
}
