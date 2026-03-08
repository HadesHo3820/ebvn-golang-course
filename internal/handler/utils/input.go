// Package utils provides utility functions for HTTP handlers.
// It contains reusable helper functions that simplify common handler operations
// such as request binding and validation.
package utils

import (
	"net/http"
	"regexp"

	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Password validation patterns - Go's regexp doesn't support Perl lookahead (?=),
// so we check each requirement separately.
var (
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasDigit   = regexp.MustCompile(`\d`)
	hasSpecial = regexp.MustCompile(`[@$!%*?&]`)
)

// validatePassword is a custom validator for password fields.
// It checks that the password contains:
//   - At least one lowercase letter
//   - At least one uppercase letter
//   - At least one digit
//   - At least one special character (@$!%*?&)
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return hasLower.MatchString(password) &&
		hasUpper.MatchString(password) &&
		hasDigit.MatchString(password) &&
		hasSpecial.MatchString(password)
}

// BindInputFromRequest is a generic utility function that binds request data from multiple
// sources into a single struct and validates it. This eliminates boilerplate code in handlers
// by consolidating all binding and validation logic into one reusable function.
//
// The function binds data from the following sources in order:
//  1. JSON body - using struct tags `json:"fieldName"`
//  2. URI path parameters - using struct tags `uri:"paramName"`
//  3. Query string parameters - using struct tags `form:"paramName"`
//  4. HTTP headers - using struct tags `header:"headerName"`
//
// After binding, it validates the struct using go-playground/validator with the
// `validate:"rule"` struct tags.
//
// Custom validators available:
//   - password: Validates password contains uppercase, lowercase, digit, and special character
//
// Type Parameters:
//   - T: The type of the input struct to bind to. Must have appropriate struct tags
//     for the desired binding sources and validation rules.
//
// Parameters:
//   - c: The Gin context containing the HTTP request
//
// Returns:
//   - *T: A pointer to the populated and validated struct, or nil if binding/validation fails
//   - error: The binding or validation error, or nil on success
//
// On error, this function automatically:
//   - Responds with HTTP 400 Bad Request and error details
//   - Calls c.Abort() to prevent subsequent handlers from executing
//
// Example usage:
//
//	type UserInput struct {
//	    Username string `json:"username" validate:"required,min=3"`
//	    Password string `json:"password" validate:"required,gte=8,password"`
//	    UserID   string `uri:"id"`
//	}
//
//	func (h *Handler) GetUser(c *gin.Context) {
//	    input, err := utils.BindInputFromRequest[UserInput](c)
//	    if err != nil {
//	        return // Response already sent
//	    }
//	    // Use input.Username, input.UserID...
//	}
func BindInputFromRequest[T any](c *gin.Context) (*T, error) {
	reqInput := new(T)

	// Skip JSON binding for GET requests to avoid EOF error on empty body
	if c.Request.Method != http.MethodGet {
		if err := c.ShouldBindJSON(reqInput); err != nil && err.Error() != "EOF" {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
			return nil, err
		}
	}

	if err := c.ShouldBindUri(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := c.ShouldBindQuery(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := c.ShouldBindHeader(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	// Register custom password validator
	validate.RegisterValidation("password", validatePassword)

	if err := validate.Struct(reqInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	return reqInput, nil
}

// BindInputFromRequestWithAuth is a convenience function that combines input binding
// and user authentication in a single call. It first binds and validates the request
// data using BindInputFromRequest, then extracts the user ID from the JWT token.
//
// This function is useful for handlers that require both input validation and
// authenticated user identification, reducing boilerplate code.
//
// Type Parameters:
//   - T: The type of the input struct to bind to. Must have appropriate struct tags
//     for the desired binding sources and validation rules.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and JWT claims
//
// Returns:
//   - *T: A pointer to the populated and validated struct, or nil if binding/validation/auth fails
//   - string: The user ID extracted from the JWT token, or empty string on failure
//   - error: The binding, validation, or authentication error, or nil on success
//
// Possible errors:
//   - Any error from BindInputFromRequest (binding or validation failures)
//   - ErrInvalidToken: If the JWT claims are not present or invalid in the context
//   - ErrEmptyUID: If the user ID is missing or empty in the token claims
func BindInputFromRequestWithAuth[T any](c *gin.Context) (*T, string, error) {
	input, err := BindInputFromRequest[T](c)
	if err != nil {
		return nil, "", err
	}

	uid, err := GetUIDFromRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid token",
		})
		return nil, "", err
	}

	return input, uid, nil
}
