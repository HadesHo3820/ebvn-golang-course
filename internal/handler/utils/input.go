// Package utils provides utility functions for HTTP handlers.
// It contains reusable helper functions that simplify common handler operations
// such as request binding and validation.
package utils

import (
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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

	if err := c.ShouldBindJSON(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindUri(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindQuery(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindHeader(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	return reqInput, nil
}
