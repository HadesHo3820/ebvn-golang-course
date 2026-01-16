// Package handler contains HTTP request handlers that process incoming requests,
// validate input, delegate to service layer, and format responses.
// Handlers are the entry point for API endpoints.
package handler

import (
	"net/http"

	// Blank import for swag to generate Swagger documentation for model.User type.
	_ "github.com/HadesHo3820/ebvn-golang-course/internal/model"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
)

// User defines the interface for user-related HTTP handlers.
// This abstraction allows for easy testing and swapping of implementations.
type User interface {
	// RegisterUser handles HTTP POST requests for user registration.
	RegisterUser(c *gin.Context)
}

// user is the concrete implementation of the User handler interface.
// It delegates business logic to the service layer.
type user struct {
	svc service.User
}

// NewUser creates and returns a new User handler instance.
// It requires a User service for business logic operations.
//
// Parameters:
//   - svc: A service.User implementation for user operations
//
// Returns:
//   - User: An implementation of the User handler interface
func NewUser(svc service.User) User {
	return &user{svc: svc}
}

// registerInputBody represents the expected JSON request body for user registration.
// All fields are required for successful registration.
type registerInputBody struct {
	// Username is the unique login identifier for the user.
	Username string `json:"username"`
	// Password is the user's plain-text password (will be hashed by service).
	Password string `json:"password"`
	// DisplayName is the user's display name shown in the UI.
	DisplayName string `json:"display_name"`
	// Email is the user's email address.
	Email string `json:"email"`
}

// RegisterUser handles user registration requests.
// It validates the JSON input, delegates to the service layer for user creation,
// and returns the created user or an appropriate error response.
//
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags User
// @Accept json
// @Produce json
// @Param body body registerInputBody true "User registration details"
// @Success 200 {object} model.User
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /v1/users/register [post]
func (u *user) RegisterUser(c *gin.Context) {
	inputBody := &registerInputBody{}
	if err := c.ShouldBindJSON(inputBody); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	// call service
	res, err := u.svc.CreateUser(c, inputBody.Username, inputBody.Password, inputBody.DisplayName, inputBody.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
