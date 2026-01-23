// Package common provides shared utility functions used across the application.
// It contains helper functions for common operations like error handling.
package common

// HandleError is a utility function that panics if the provided error is non-nil.
// This function is intended for use during application startup or initialization
// where an error represents an unrecoverable condition that should halt execution.
//
// WARNING: This function should NOT be used in request handlers or business logic
// where errors should be handled gracefully and returned to the caller.
//
// Parameters:
//   - err: The error to check. If non-nil, the function panics with this error.
//
// Example usage:
//
//	db, err := sql.Open("postgres", connectionString)
//	common.HandleError(err) // Panics if database connection fails
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
