package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HadesHo3820/ebvn-golang-course/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPasswordHandler_GenPass(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func() *mocks.Password

		expectedStatus int
		expectedResp   string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				// Create virtual request
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				// Mock the GeneratePassword method that is called in the handler
				svcMock.On("GeneratePassword").Return("123456789", nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedResp:   "123456789",
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				// Create virtual request
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				// Mock the GeneratePassword method that is called in the handler
				svcMock.On("GeneratePassword").Return("", errors.New("something"))
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "err",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a httptest.NewRecorder to capture status code and response for that request
			rec := httptest.NewRecorder()
			// Create a Gin test context to simulate a request
			gctx, _ := gin.CreateTestContext(rec)

			// Setup the request
			tc.setupRequest(gctx)

			// Setup the mock service
			svcMock := tc.setupMockSvc()

			// Create the handler with the mock service
			handler := NewPassword(svcMock)

			// Call the handler
			handler.GenPass(gctx)

			// Check the response and status code
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
