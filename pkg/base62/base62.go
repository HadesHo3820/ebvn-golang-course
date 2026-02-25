// Package base62 provides encoding and decoding of integers using base62 representation.
// Base62 uses the character set: 0-9, a-z, A-Z (62 characters total).
//
// This is useful for generating short, human-readable, URL-safe codes from
// auto-incrementing database IDs. For example:
//
//	base62.Encode(1)    → "1"
//	base62.Encode(62)   → "10"
//	base62.Encode(1000) → "g8"
package base62

import (
	"errors"
	"math"
	"strings"
)

// charset defines the 62-character alphabet used for encoding.
// Characters are ordered: digits (0-9), lowercase (a-z), uppercase (A-Z).
const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const base = 62

// ErrInvalidCharacter is returned when Decode encounters a character not in the base62 charset.
var ErrInvalidCharacter = errors.New("base62: invalid character")

// ErrNegativeNumber is returned when Encode receives a negative number.
var ErrNegativeNumber = errors.New("base62: negative number")

// ErrOverflow is returned when Decode would produce a value exceeding int64 max.
var ErrOverflow = errors.New("base62: value overflow")

// Encode converts a non-negative int64 to its base62 string representation.
// Returns an error if the input is negative.
//
// Algorithm: Repeated division by 62 (similar to converting decimal → any base).
//
// How it works:
//  1. Divide the number by 62 repeatedly.
//  2. At each step, the remainder (num % 62) gives an index into the charset.
//  3. The character at that index becomes part of the encoded string.
//  4. The quotient (num / 62) becomes the new number for the next iteration.
//  5. Since remainders produce digits from least-significant to most-significant,
//     the final string is reversed to get the correct order.
//
// Example: Encode(1000)
//
//	Step 1: 1000 % 62 = 8  → charset[8]  = '8'    | 1000 / 62 = 16
//	Step 2:   16 % 62 = 16 → charset[16] = 'g'    |   16 / 62 = 0   (stop)
//	Built string (reverse order): "8g"
//	Reversed (final result):      "g8"
func Encode(num int64) (string, error) {
	if num < 0 {
		return "", ErrNegativeNumber
	}
	if num == 0 {
		return "0", nil
	}

	// Build the encoded string by extracting digits from least-significant to most-significant.
	// Each iteration extracts one base62 digit:
	//   remainder = num % 62  → index into charset (the current digit)
	//   num       = num / 62  → shift right by one base62 position
	var result strings.Builder
	for num > 0 {
		remainder := num % base
		result.WriteByte(charset[remainder])
		num /= base
	}

	// The loop above produces digits in reverse order (least-significant first),
	// so we reverse to get the correct most-significant-first representation.
	// Example: 1000 → "8g" (reversed) → "g8" (correct)
	encoded := result.String()
	runes := []rune(encoded)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}

// Decode converts a base62-encoded string back to its int64 representation.
// Returns an error if the string contains invalid characters or if the result
// would overflow int64.
//
// Algorithm: Horner's method — process characters left-to-right, accumulating
// the value by multiplying by the base and adding the next digit's value.
//
// How it works:
//  1. Start with result = 0.
//  2. For each character (left to right), find its index in the charset.
//     That index IS the digit's numeric value (0-61).
//  3. Multiply the current result by 62 (shift left by one base62 position),
//     then add the digit's value.
//  4. Repeat until all characters are processed.
//
// Example: Decode("g8")
//
//	Step 1: result = 0*62 + indexOf('g') = 0  + 16 = 16
//	Step 2: result = 16*62 + indexOf('8') = 992 +  8 = 1000
//	Final result: 1000
//
// This is the inverse of Encode: Encode(1000) = "g8", Decode("g8") = 1000.
func Decode(encoded string) (int64, error) {
	if len(encoded) == 0 {
		return 0, ErrInvalidCharacter
	}

	var result int64
	for _, c := range encoded {
		// Find the numeric value of this character (its position in the charset).
		// For example: '0'→0, '9'→9, 'a'→10, 'z'→35, 'A'→36, 'Z'→61
		idx := strings.IndexRune(charset, c)
		if idx == -1 {
			return 0, ErrInvalidCharacter
		}

		// Overflow check: ensure result*62 + idx won't exceed math.MaxInt64.
		// Rearranged as: result > (MaxInt64 - idx) / 62
		if result > (math.MaxInt64-int64(idx))/base {
			return 0, ErrOverflow
		}

		// Shift the accumulated result left by one base62 position and add the new digit.
		// This is equivalent to: result = result * 62 + digitValue
		result = result*base + int64(idx)
	}

	return result, nil
}
