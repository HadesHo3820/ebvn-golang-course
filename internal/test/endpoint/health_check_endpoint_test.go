package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/api"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		name            string
		setupTestHTTP   func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: service.HealthCheckOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := tc.setupTestHTTP(api.New(cfg))

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]any
			json.Unmarshal(rec.Body.Bytes(), &resp)

			assert.Equal(t, tc.expectedMessage, resp["message"])
		})
	}
}
