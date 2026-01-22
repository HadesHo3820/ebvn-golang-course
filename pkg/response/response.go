// Package response provides standardized structures and helper functions for HTTP API responses.
// It defines a consistent JSON format for sending messages and errors back to the client.
package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Message represents the structure of a response message.
// It is used to return success notifications or error details to the client.
type Message struct {
	// Message is a brief summary of the response (e.g., "Input error").
	Message string `json:"message"`
	// Details contains a list of specific error messages or additional information, if any.
	// If empty, this field is omitted from the JSON response.
	Details any `json:"details,omitempty"`
}

const (
	InternalErrMessage = "Processing error"
	InputErrMessage    = "Input error"
)

// Common response messages used throughout the application.
var (
	// InternalErrResponse is a generic response for internal server errors.
	// It should be used when the server encounters an unexpected condition.
	InternalErrResponse = Message{
		Message: InternalErrMessage,
		Details: nil,
	}
	// InputErrResponse is a generic response for client-side input errors.
	// It serves as a fallback when specific validation details are not available.
	InputErrResponse = Message{
		Message: InputErrMessage,
		Details: nil,
	}
)

// InputFieldError processes an error to construct a detailed validation error response.
// If the error is of type validator.ValidationErrors, it extracts and formats
// validation messages for each invalid field.
// For other error types, it returns a generic InputErrResponse.
func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}

	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}

	return Message{
		Message: InputErrMessage,
		Details: errs,
	}
}
