package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Test structs for binding tests
type testJSONBody struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type testURIParams struct {
	ID string `uri:"id" validate:"required"`
}

type testQueryParams struct {
	Page  int    `form:"page" validate:"gte=1"`
	Limit int    `form:"limit" validate:"gte=1,lte=100"`
	Sort  string `form:"sort"`
}

type testHeaderParams struct {
	APIKey string `header:"X-API-Key" validate:"required"`
}

type testPasswordInput struct {
	Password string `json:"password" validate:"required,password"`
}

type testCombinedInput struct {
	Username string `json:"username" validate:"required"`
	ID       string `uri:"id"`
	Page     int    `form:"page"`
}

// TestBindInputFromRequest_JSONBinding tests JSON body binding.
func TestBindInputFromRequest_JSONBinding(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
		validateResult func(t *testing.T, result *testJSONBody)
	}{
		{
			name:           "success - valid JSON",
			requestBody:    `{"username": "testuser", "email": "test@example.com"}`,
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResult: func(t *testing.T, result *testJSONBody) {
				assert.Equal(t, "testuser", result.Username)
				assert.Equal(t, "test@example.com", result.Email)
			},
		},
		{
			name:           "error - invalid JSON syntax",
			requestBody:    `{"username": "testuser", "email": }`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "error - missing required field",
			requestBody:    `{"username": "testuser"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "error - invalid email format",
			requestBody:    `{"username": "testuser", "email": "not-an-email"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gctx, _ := gin.CreateTestContext(rec)

			gctx.Request = httptest.NewRequest(
				http.MethodPost,
				"/test",
				bytes.NewBufferString(tc.requestBody),
			)
			gctx.Request.Header.Set("Content-Type", "application/json")

			result, err := BindInputFromRequest[testJSONBody](gctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.validateResult != nil {
					tc.validateResult(t, result)
				}
			}
		})
	}
}

// TestBindInputFromRequest_URIBinding tests URI parameter binding.
func TestBindInputFromRequest_URIBinding(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		uriParams      gin.Params
		expectedStatus int
		expectError    bool
		validateResult func(t *testing.T, result *testURIParams)
	}{
		{
			name: "success - valid URI param",
			uriParams: gin.Params{
				{Key: "id", Value: "user-123"},
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResult: func(t *testing.T, result *testURIParams) {
				assert.Equal(t, "user-123", result.ID)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gctx, _ := gin.CreateTestContext(rec)

			gctx.Request = httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString("{}"))
			gctx.Request.Header.Set("Content-Type", "application/json")
			gctx.Params = tc.uriParams

			result, err := BindInputFromRequest[testURIParams](gctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.validateResult != nil {
					tc.validateResult(t, result)
				}
			}
		})
	}
}

// TestBindInputFromRequest_QueryBinding tests query parameter binding.
func TestBindInputFromRequest_QueryBinding(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		queryString    string
		expectedStatus int
		expectError    bool
		validateResult func(t *testing.T, result *testQueryParams)
	}{
		{
			name:           "success - valid query params",
			queryString:    "?page=1&limit=10&sort=name",
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResult: func(t *testing.T, result *testQueryParams) {
				assert.Equal(t, 1, result.Page)
				assert.Equal(t, 10, result.Limit)
				assert.Equal(t, "name", result.Sort)
			},
		},
		{
			name:           "error - page less than 1",
			queryString:    "?page=0&limit=10",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "error - limit exceeds max",
			queryString:    "?page=1&limit=101",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gctx, _ := gin.CreateTestContext(rec)

			gctx.Request = httptest.NewRequest(http.MethodGet, "/test"+tc.queryString, bytes.NewBufferString("{}"))
			gctx.Request.Header.Set("Content-Type", "application/json")

			result, err := BindInputFromRequest[testQueryParams](gctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.validateResult != nil {
					tc.validateResult(t, result)
				}
			}
		})
	}
}

// TestBindInputFromRequest_PasswordValidation tests custom password validation.
func TestBindInputFromRequest_PasswordValidation(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "success - valid password",
			requestBody:    `{"password": "Password1!"}`,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "error - missing uppercase",
			requestBody:    `{"password": "password1!"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "error - missing lowercase",
			requestBody:    `{"password": "PASSWORD1!"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "error - missing digit",
			requestBody:    `{"password": "Password!"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "error - missing special character",
			requestBody:    `{"password": "Password1"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gctx, _ := gin.CreateTestContext(rec)

			gctx.Request = httptest.NewRequest(
				http.MethodPost,
				"/test",
				bytes.NewBufferString(tc.requestBody),
			)
			gctx.Request.Header.Set("Content-Type", "application/json")

			result, err := BindInputFromRequest[testPasswordInput](gctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestValidatePassword tests the password validation function directly.
func TestValidatePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		password string
		expected bool
	}{
		{"valid - all requirements", "Password1!", true},
		{"valid - complex password", "P@ssw0rd123!", true},
		{"invalid - no lowercase", "PASSWORD1!", false},
		{"invalid - no uppercase", "password1!", false},
		{"invalid - no digit", "Password!", false},
		{"invalid - no special char", "Password1", false},
		{"invalid - empty", "", false},
		{"invalid - only lowercase", "password", false},
		{"invalid - only numbers", "12345678", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := hasLower.MatchString(tc.password) &&
				hasUpper.MatchString(tc.password) &&
				hasDigit.MatchString(tc.password) &&
				hasSpecial.MatchString(tc.password)

			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestBindInputFromRequest_HeaderBinding tests header parameter binding.
func TestBindInputFromRequest_HeaderBinding(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
		expectError    bool
		validateResult func(t *testing.T, result *testHeaderParams)
	}{
		{
			name: "success - valid header",
			headers: map[string]string{
				"X-API-Key": "my-secret-key",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResult: func(t *testing.T, result *testHeaderParams) {
				assert.Equal(t, "my-secret-key", result.APIKey)
			},
		},
		{
			name:           "error - missing required header",
			headers:        map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gctx, _ := gin.CreateTestContext(rec)

			gctx.Request = httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString("{}"))
			gctx.Request.Header.Set("Content-Type", "application/json")
			for key, value := range tc.headers {
				gctx.Request.Header.Set(key, value)
			}

			result, err := BindInputFromRequest[testHeaderParams](gctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.validateResult != nil {
					tc.validateResult(t, result)
				}
			}
		})
	}
}

// TestBindInputFromRequest_CombinedBinding tests binding from multiple sources.
func TestBindInputFromRequest_CombinedBinding(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		requestBody    string
		uriParams      gin.Params
		queryString    string
		expectedStatus int
		expectError    bool
		validateResult func(t *testing.T, result *testCombinedInput)
	}{
		{
			name:        "success - all sources combined",
			requestBody: `{"username": "testuser"}`,
			uriParams: gin.Params{
				{Key: "id", Value: "user-456"},
			},
			queryString:    "?page=5",
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResult: func(t *testing.T, result *testCombinedInput) {
				assert.Equal(t, "testuser", result.Username)
				assert.Equal(t, "user-456", result.ID)
				assert.Equal(t, 5, result.Page)
			},
		},
		{
			name:        "error - missing required JSON field",
			requestBody: `{}`,
			uriParams: gin.Params{
				{Key: "id", Value: "user-456"},
			},
			queryString:    "?page=5",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gctx, _ := gin.CreateTestContext(rec)

			gctx.Request = httptest.NewRequest(
				http.MethodPost,
				"/test"+tc.queryString,
				bytes.NewBufferString(tc.requestBody),
			)
			gctx.Request.Header.Set("Content-Type", "application/json")
			gctx.Params = tc.uriParams

			result, err := BindInputFromRequest[testCombinedInput](gctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.validateResult != nil {
					tc.validateResult(t, result)
				}
			}
		})
	}
}

// TestBindInputFromRequestWithAuth tests the combined binding and authentication function.
func TestBindInputFromRequestWithAuth(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		requestBody    string
		claims         any
		expectedStatus int
		expectError    bool
		expectedUID    string
		validateResult func(t *testing.T, result *testJSONBody, uid string)
	}{
		{
			name:        "success - valid input and auth",
			requestBody: `{"username": "testuser", "email": "test@example.com"}`,
			claims: jwt.MapClaims{
				"sub": "user-123",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedUID:    "user-123",
			validateResult: func(t *testing.T, result *testJSONBody, uid string) {
				assert.Equal(t, "testuser", result.Username)
				assert.Equal(t, "test@example.com", result.Email)
				assert.Equal(t, "user-123", uid)
			},
		},
		{
			name:           "error - invalid JSON binding",
			requestBody:    `{"username": "testuser"}`, // missing required email
			claims:         jwt.MapClaims{"sub": "user-123"},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "error - missing claims (invalid token)",
			requestBody:    `{"username": "testuser", "email": "test@example.com"}`,
			claims:         nil,           // no claims set
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "error - invalid claims type",
			requestBody:    `{"username": "testuser", "email": "test@example.com"}`,
			claims:         "invalid-claims", // wrong type
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:        "error - empty uid in claims",
			requestBody: `{"username": "testuser", "email": "test@example.com"}`,
			claims: jwt.MapClaims{
				"sub": "", // empty uid
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:        "error - missing sub claim",
			requestBody: `{"username": "testuser", "email": "test@example.com"}`,
			claims: jwt.MapClaims{
				"exp": 1234567890, // no sub claim
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gctx, _ := gin.CreateTestContext(rec)

			gctx.Request = httptest.NewRequest(
				http.MethodPost,
				"/test",
				bytes.NewBufferString(tc.requestBody),
			)
			gctx.Request.Header.Set("Content-Type", "application/json")

			// Set claims in context if provided
			if tc.claims != nil {
				gctx.Set("claims", tc.claims)
			}

			result, uid, err := BindInputFromRequestWithAuth[testJSONBody](gctx)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Empty(t, uid)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedUID, uid)
				if tc.validateResult != nil {
					tc.validateResult(t, result, uid)
				}
			}
		})
	}
}
