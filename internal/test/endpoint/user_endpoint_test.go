// Package endpoint provides integration tests for API endpoints.
//
// Unlike unit tests that mock dependencies, endpoint tests validate the full HTTP stack
// including routing, middleware, handlers, and real service implementations.
package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/test/fixture"
	jwtMocks "github.com/HadesHo3820/ebvn-golang-course/pkg/jwtutils/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// testValidUserDisplayName is the display name used for valid user test cases.
const testValidUserDisplayName = "Valid User"

// testValidAuthToken is the Authorization header value used for authenticated test cases.
const testValidAuthToken = "Bearer valid.jwt.token"

// testValidUsername is the username used for valid user test cases.
const testValidUsername = "validuser"

// testValidEmail is the email used for valid user test cases.
const testValidEmail = "valid@example.com"

// HTTP header constants
const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

// TestUserEndpoint_Register validates the POST /v1/users/register endpoint.
func TestUserEndpoint_Register(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedFields []string // Fields expected in response
	}{
		{
			name: "success - register new user",
			requestBody: fixture.DefaultRegisterBody(
				fixture.WithField("username", "newuser"),
				fixture.WithField("display_name", "New User"),
				fixture.WithField("email", "newuser@example.com"),
			),
			expectedStatus: http.StatusOK,
			expectedFields: []string{"message", "data"},
		},
		{
			name: "error - duplicate username",
			requestBody: fixture.DefaultRegisterBody(
				fixture.WithField("username", "johnny.ho"), // Already exists in fixture
				fixture.WithField("display_name", "Johnny Ho"),
				fixture.WithField("email", "another@example.com"),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - invalid email format",
			requestBody: fixture.DefaultRegisterBody(
				fixture.WithField("username", testValidUsername),
				fixture.WithField("display_name", testValidUserDisplayName),
				fixture.WithField("email", "invalid-email"),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - password too short",
			requestBody: fixture.DefaultRegisterBody(
				fixture.WithField("username", testValidUsername),
				fixture.WithField("password", "Pass1!"),
				fixture.WithField("display_name", testValidUserDisplayName),
				fixture.WithField("email", testValidEmail),
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - password missing special character",
			requestBody: fixture.DefaultRegisterBody(
				fixture.WithField("username", testValidUsername),
				fixture.WithField("password", "Password123"),
				fixture.WithField("display_name", testValidUserDisplayName),
				fixture.WithField("email", testValidEmail),
			),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper
			testEngine := NewTestEngine(&TestEngineOpts{T: t})

			// Create request
			bodyBytes, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(bodyBytes))
			req.Header.Set(contentTypeHeader, contentTypeJSON)

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Assert response contains expected fields for success cases
			if tc.expectedStatus == http.StatusOK && tc.expectedFields != nil {
				var body map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				for _, field := range tc.expectedFields {
					assert.Contains(t, body, field)
				}
			}
		})
	}
}

// TestUserEndpoint_Login validates the POST /v1/users/login endpoint.
func TestUserEndpoint_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func(*jwtMocks.JWTGenerator)
		expectedStatus int
		expectedToken  string
	}{
		{
			name:           "error - user not found",
			requestBody:    fixture.DefaultLoginBody(fixture.WithField("username", "nonexistent")),
			setupMock:      func(m *jwtMocks.JWTGenerator) {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "error - password too short",
			requestBody: fixture.DefaultLoginBody(
				fixture.WithField("username", "johnny.ho"),
				fixture.WithField("password", "short"),
			),
			setupMock:      func(m *jwtMocks.JWTGenerator) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper
			testEngine := NewTestEngine(&TestEngineOpts{T: t})

			// Setup mock expectations
			tc.setupMock(testEngine.JwtGen)

			// Create request
			bodyBytes, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(bodyBytes))
			req.Header.Set(contentTypeHeader, contentTypeJSON)

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

// TestUserEndpoint_GetSelfInfo validates the GET /v1/self/info endpoint.
func TestUserEndpoint_GetSelfInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		authToken      string
		setupMock      func(*jwtMocks.JWTValidator) jwt.MapClaims
		expectedStatus int
	}{
		{
			name:      "success - get user profile",
			authToken: testValidAuthToken,
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID)) // johnny.ho from fixture
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "error - user not found",
			authToken: testValidAuthToken,
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", "non-existent-user-id"))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - missing authorization",
			authToken:      "",
			setupMock:      func(m *jwtMocks.JWTValidator) jwt.MapClaims { return nil },
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper
			testEngine := NewTestEngine(&TestEngineOpts{T: t})

			// Setup mock expectations
			tc.setupMock(testEngine.JwtValidator)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

// TestUserEndpoint_UpdateSelfInfo validates the PUT /v1/self/info endpoint.
func TestUserEndpoint_UpdateSelfInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		authToken      string
		requestBody    map[string]string
		setupMock      func(*jwtMocks.JWTValidator) jwt.MapClaims
		expectedStatus int
	}{
		{
			name:        "success - update display name",
			authToken:   testValidAuthToken,
			requestBody: fixture.DefaultUpdateUserBody(fixture.WithField("email", "")),
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID)) // johnny.ho from fixture
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success - update email",
			authToken:   testValidAuthToken,
			requestBody: fixture.DefaultUpdateUserBody(fixture.WithField("display_name", "")),
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "error - no update data",
			authToken:   testValidAuthToken,
			requestBody: map[string]string{},
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "error - invalid email format",
			authToken: testValidAuthToken,
			requestBody: fixture.DefaultUpdateUserBody(
				fixture.WithField("display_name", ""),
				fixture.WithField("email", "invalid-email"),
			),
			setupMock: func(m *jwtMocks.JWTValidator) jwt.MapClaims {
				claims := fixture.DefaultJWTClaims(fixture.WithClaim("sub", fixture.FixtureUserOneID))
				m.On("ValidateToken", mock.Anything).Return(claims, nil)
				return claims
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error - missing authorization",
			authToken:      "",
			requestBody:    fixture.DefaultUpdateUserBody(fixture.WithField("email", "")),
			setupMock:      func(m *jwtMocks.JWTValidator) jwt.MapClaims { return nil },
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create test engine with helper
			testEngine := NewTestEngine(&TestEngineOpts{T: t})

			// Setup mock expectations
			tc.setupMock(testEngine.JwtValidator)

			// Create request
			bodyBytes, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewReader(bodyBytes))
			req.Header.Set(contentTypeHeader, contentTypeJSON)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}

			rec := httptest.NewRecorder()
			testEngine.Engine.ServeHTTP(rec, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
