package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertStatusCode asserts that the response has the expected HTTP status code.
//
// Parameters:
//   - t: Testing instance
//   - recorder: HTTP response recorder from the handler call
//   - expected: Expected HTTP status code
func AssertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()
	assert.Equal(t, expected, recorder.Code)
}

// AssertJSONResponse unmarshals the response body and compares it to the expected map.
// It first asserts the status code, then compares the JSON body.
//
// Parameters:
//   - t: Testing instance
//   - recorder: HTTP response recorder from the handler call
//   - expectedStatus: Expected HTTP status code
//   - expectedBody: Expected response body as a map (nil to skip body assertion)
func AssertJSONResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]any) {
	t.Helper()
	AssertStatusCode(t, recorder, expectedStatus)

	if expectedBody != nil {
		var actualBody map[string]any
		err := json.Unmarshal(recorder.Body.Bytes(), &actualBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, actualBody)
	}
}

// AssertJSONContainsFields asserts that the response body contains the specified fields.
// This is useful when you only care about certain fields being present.
//
// Parameters:
//   - t: Testing instance
//   - recorder: HTTP response recorder from the handler call
//   - fields: List of field names that should be present in the response
func AssertJSONContainsFields(t *testing.T, recorder *httptest.ResponseRecorder, fields []string) {
	t.Helper()
	var body map[string]any
	err := json.Unmarshal(recorder.Body.Bytes(), &body)
	assert.NoError(t, err)

	for _, field := range fields {
		assert.Contains(t, body, field)
	}
}
