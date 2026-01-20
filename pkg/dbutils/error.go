// Package dbutils provides database utility functions for error handling and translation.
// It abstracts database-specific errors into application-level errors, allowing the
// service layer to handle errors without coupling to specific database implementations.
package dbutils

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// errorFilters is a list of filter functions that check for specific database error types.
// Each filter returns true if the error matches its criteria, along with a translated error.
// Filters are applied in order until one matches.
var errorFilters = []func(err error) (bool, error){
	filterDuplicationType,
	filterRecordNotFound,
}

// CatchDBErr translates database-specific errors into application-level errors.
// It iterates through registered error filters and returns the first matching
// translated error. If no filter matches, the original error is returned unchanged.
//
// This function enables the repository layer to return consistent, database-agnostic
// errors that the service layer can handle without knowing the underlying database.
//
// Parameters:
//   - err: The database error to translate
//
// Returns:
//   - error: A translated application-level error, or the original error if no filter matches
//
// Example usage:
//
//	err := db.Create(&user).Error
//	if err != nil {
//	    return nil, dbutils.CatchDBErr(err) // Returns ErrDuplicationType for unique constraint violations
//	}
func CatchDBErr(err error) error {
	for _, filter := range errorFilters {
		if ok, newErr := filter(err); ok {
			return newErr
		}
	}
	return err
}

// Sentinel errors for common database error conditions.
// These errors can be used with errors.Is() for type-safe error checking.
var (
	// ErrDuplicationType indicates a unique constraint violation occurred,
	// typically when inserting a record with a duplicate key or unique field value.
	ErrDuplicationType = errors.New("duplication type")

	// ErrNotFoundType indicates the requested record was not found in the database.
	ErrNotFoundType = errors.New("not found type")
)

// filterDuplicationType checks if the error is a unique constraint violation.
// It performs a case-insensitive search for "unique constraint" in the error message,
// which works across different database implementations (PostgreSQL, MySQL, etc.).
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - bool: true if the error indicates a unique constraint violation
//   - error: ErrDuplicationType if matched
func filterDuplicationType(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint"), ErrDuplicationType
}

// filterRecordNotFound checks if the error is a GORM record not found error.
// This filter specifically handles gorm.ErrRecordNotFound which is returned
// when a query expects to find a record but none exists (e.g., First(), Take()).
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - bool: true if the error is gorm.ErrRecordNotFound
//   - error: ErrNotFoundType if matched
func filterRecordNotFound(err error) (bool, error) {
	return errors.Is(err, gorm.ErrRecordNotFound), ErrNotFoundType
}
