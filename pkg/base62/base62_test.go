package base62

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       int64
		expected    string
		expectedErr error
	}{
		{name: "zero", input: 0, expected: "0"},
		{name: "single digit", input: 1, expected: "1"},
		{name: "nine", input: 9, expected: "9"},
		{name: "ten maps to a", input: 10, expected: "a"},
		{name: "thirty five maps to z", input: 35, expected: "z"},
		{name: "thirty six maps to A", input: 36, expected: "A"},
		{name: "sixty one maps to Z", input: 61, expected: "Z"},
		{name: "sixty two maps to 10", input: 62, expected: "10"},
		{name: "thousand", input: 1000, expected: "g8"},
		{name: "large number", input: 123456789, expected: "8m0Kx"},
		{name: "negative number", input: -1, expected: "", expectedErr: ErrNegativeNumber},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := Encode(tc.input)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		expected    int64
		expectedErr error
	}{
		{name: "zero", input: "0", expected: 0},
		{name: "single digit", input: "1", expected: 1},
		{name: "nine", input: "9", expected: 9},
		{name: "a maps to ten", input: "a", expected: 10},
		{name: "z maps to thirty five", input: "z", expected: 35},
		{name: "A maps to thirty six", input: "A", expected: 36},
		{name: "Z maps to sixty one", input: "Z", expected: 61},
		{name: "10 maps to sixty two", input: "10", expected: 62},
		{name: "g8 maps to thousand", input: "g8", expected: 1000},
		{name: "large number", input: "8m0Kx", expected: 123456789},
		{name: "empty string", input: "", expected: 0, expectedErr: ErrInvalidCharacter},
		{name: "invalid character", input: "abc!def", expected: 0, expectedErr: ErrInvalidCharacter},
		{name: "special character", input: "#", expected: 0, expectedErr: ErrInvalidCharacter},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := Decode(tc.input)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// TestRoundTrip ensures Encode and Decode are inverse operations.
func TestRoundTrip(t *testing.T) {
	t.Parallel()

	values := []int64{0, 1, 10, 61, 62, 100, 1000, 123456789, math.MaxInt32}

	for _, v := range values {
		encoded, err := Encode(v)
		assert.NoError(t, err)

		decoded, err := Decode(encoded)
		assert.NoError(t, err)
		assert.Equal(t, v, decoded)
	}
}
