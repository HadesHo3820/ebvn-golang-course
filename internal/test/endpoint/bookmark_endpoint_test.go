package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	jwtMocks "github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestBookmarkEndpoint_Create validates the POST /v1/bookmarks endpoint.
func TestBookmarkEndpoint_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		authToken      string
		requestBody    map[string]any
		setupMock      func(*jwtMocks.JWTValidator) jwt.MapClaims
		expectedStatus int
		expectedFields []string
	}{
		{
			name:      "success - create bookmark",
			authToken: testValidAuthToken,
			requestBody: map[string]any{
				"description": "Integration Test Bookmark",
				"url":         "https://integration-test.com",
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"id", "description", "url", "code", "user_id", "created_at"},
		},
		{
			name:      "error - unauthorized (invalid token)",
			authToken: "Bearer invalid",
			requestBody: map[string]any{
				"description": "My Bookmark",
				"url":         "https://example.com",
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				m.On("ValidateToken", mock.Anything).Return(nil, jwt.ErrTokenInvalidClaims)
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:      "error - missing authorization header",
			authToken: "",
			requestBody: map[string]any{
				"description": "My Bookmark",
				"url":         "https://example.com",
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				// No validation call expected as header is missing
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:      "error - invalid input (missing URL)",
			authToken: testValidAuthToken,
			requestBody: map[string]any{
				"description": "My Bookmark",
				// URL missing
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "error - invalid input (URL too long)",
			authToken: testValidAuthToken,
			requestBody: map[string]any{
				"description": "My Bookmark",
				"url":         "https://example.com/" + strings.Repeat("a", 2050),
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper and Bookmark fixture
			testEngine := NewTestEngine(&TestEngineOpts{
				T:       t,
				Fixture: &fixture.BookmarkCommonTestDB{},
			})

			// Setup mock expectations
			tc.setupMock(testEngine.JwtValidator)

			// Create request
			bodyBytes, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewReader(bodyBytes))
			req.Header.Set(contentTypeHeader, contentTypeJSON)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response fields
			if tc.expectedStatus == http.StatusOK && tc.expectedFields != nil {
				var body map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				for _, field := range tc.expectedFields {
					assert.Contains(t, body, field)
				}
				// Verify specific values
				assert.Equal(t, tc.requestBody["description"], body["description"])
				assert.Equal(t, tc.requestBody["url"], body["url"])
				assert.NotEmpty(t, body["code"])
			}
		})
	}
}

// TestBookmarkEndpoint_GetBookmarks validates the GET /v1/bookmarks endpoint.
func TestBookmarkEndpoint_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		authToken            string
		queryParams          string
		setupMock            func(*jwtMocks.JWTValidator) jwt.MapClaims
		expectedStatus       int
		expectedTotalRecords int
		expectedDataLength   int
	}{
		{
			name:        "success - get bookmarks for User One",
			authToken:   testValidAuthToken,
			queryParams: "", // Default pagination
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				// User 1 has 1 bookmark seeded in fixture
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus:       http.StatusOK,
			expectedTotalRecords: 1,
			expectedDataLength:   1,
		},
		{
			name:        "success - get bookmarks for User Two",
			authToken:   testValidAuthToken,
			queryParams: "",
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				// User 2 has 1 bookmark seeded in fixture
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserTwoID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus:       http.StatusOK,
			expectedTotalRecords: 1,
			expectedDataLength:   1,
		},
		{
			name:        "error - unauthorized",
			authToken:   "Bearer invalid",
			queryParams: "",
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				m.On("ValidateToken", mock.Anything).Return(nil, jwt.ErrTokenInvalidClaims)
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:        "success - custom pagination (limit 0)",
			authToken:   testValidAuthToken,
			queryParams: "?limit=0", // Should default to 10 effectively
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus:       http.StatusOK,
			expectedTotalRecords: 1,
			expectedDataLength:   1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper and Bookmark fixture
			testEngine := NewTestEngine(&TestEngineOpts{
				T:       t,
				Fixture: &fixture.BookmarkCommonTestDB{},
			})

			// Setup mock expectations
			tc.setupMock(testEngine.JwtValidator)

			// Create request
			testURL := "/v1/bookmarks"
			if tc.queryParams != "" {
				testURL += tc.queryParams
			}
			req := httptest.NewRequest(http.MethodGet, testURL, nil)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response data
			if tc.expectedStatus == http.StatusOK {
				var body map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)

				// Verify metadata
				metadata, ok := body["metadata"].(map[string]any)
				assert.True(t, ok)
				assert.EqualValues(t, tc.expectedTotalRecords, metadata["total_records"])

				// Verify data
				data, ok := body["data"].([]any)
				assert.True(t, ok)
				assert.Len(t, data, tc.expectedDataLength)

				if tc.expectedDataLength > 0 {
					firstItem := data[0].(map[string]any)
					assert.NotEmpty(t, firstItem["code"])
				}
			}
		})
	}
}

// TestBookmarkEndpoint_Update validates the PUT /v1/bookmarks/{id} endpoint.
func TestBookmarkEndpoint_Update(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		authToken      string
		bookmarkID     string
		requestBody    map[string]any
		setupMock      func(*jwtMocks.JWTValidator) jwt.MapClaims
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:       "success - update own bookmark",
			authToken:  testValidAuthToken,
			bookmarkID: fixture.FixtureBookmarkOneID,
			requestBody: map[string]any{
				"description": "Updated Description",
				"url":         "https://updated-example.com",
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Success",
			},
		},
		{
			name:       "error - update another user's bookmark",
			authToken:  testValidAuthToken,
			bookmarkID: fixture.FixtureBookmarkOneID, // Belongs to User One
			requestBody: map[string]any{
				"description": "Malicious Update",
				"url":         "https://malicious.com",
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				// User Two trying to update User One's bookmark
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserTwoID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]any{
				"message": "Bookmark not found",
			},
		},
		{
			name:       "error - bookmark not found",
			authToken:  testValidAuthToken,
			bookmarkID: "00000000-0000-0000-0000-000000000000",
			requestBody: map[string]any{
				"description": "Description",
				"url":         "https://example.com",
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]any{
				"message": "Bookmark not found",
			},
		},
		{
			name:       "error - invalid input (missing URL)",
			authToken:  testValidAuthToken,
			bookmarkID: fixture.FixtureBookmarkOneID,
			requestBody: map[string]any{
				"description": "Description",
				// URL missing
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "error - unauthorized",
			authToken:  "Bearer invalid",
			bookmarkID: fixture.FixtureBookmarkOneID,
			requestBody: map[string]any{
				"description": "Description",
				"url":         "https://example.com",
			},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				m.On("ValidateToken", mock.Anything).Return(nil, jwt.ErrTokenInvalidClaims)
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper and Bookmark fixture
			testEngine := NewTestEngine(&TestEngineOpts{
				T:       t,
				Fixture: &fixture.BookmarkCommonTestDB{},
			})

			// Setup mock expectations
			tc.setupMock(testEngine.JwtValidator)

			// Create request
			bodyBytes, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+tc.bookmarkID, bytes.NewReader(bodyBytes))
			req.Header.Set(contentTypeHeader, contentTypeJSON)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body if expected
			if tc.expectedBody != nil {
				var body map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				for key, expected := range tc.expectedBody {
					assert.Equal(t, expected, body[key])
				}
			}
		})
	}
}

// TestBookmarkEndpoint_Delete validates the DELETE /v1/bookmarks/{id} endpoint.
func TestBookmarkEndpoint_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		authToken      string
		bookmarkID     string
		setupMock      func(*jwtMocks.JWTValidator) jwt.MapClaims
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:       "success - delete own bookmark",
			authToken:  testValidAuthToken,
			bookmarkID: fixture.FixtureBookmarkOneID,
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Success",
			},
		},
		{
			name:       "error - delete another user's bookmark",
			authToken:  testValidAuthToken,
			bookmarkID: fixture.FixtureBookmarkOneID, // Belongs to User One
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				// User Two trying to delete User One's bookmark
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserTwoID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]any{
				"message": "Bookmark not found",
			},
		},
		{
			name:       "error - bookmark not found",
			authToken:  testValidAuthToken,
			bookmarkID: "00000000-0000-0000-0000-000000000000",
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]any{
				"message": "Bookmark not found",
			},
		},
		{
			name:       "error - unauthorized",
			authToken:  "Bearer invalid",
			bookmarkID: fixture.FixtureBookmarkOneID,
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				m.On("ValidateToken", mock.Anything).Return(nil, jwt.ErrTokenInvalidClaims)
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper and Bookmark fixture
			testEngine := NewTestEngine(&TestEngineOpts{
				T:       t,
				Fixture: &fixture.BookmarkCommonTestDB{},
			})

			// Setup mock expectations
			tc.setupMock(testEngine.JwtValidator)

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+tc.bookmarkID, nil)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response body if expected
			if tc.expectedBody != nil {
				var body map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				for key, expected := range tc.expectedBody {
					assert.Equal(t, expected, body[key])
				}
			}
		})
	}
}
