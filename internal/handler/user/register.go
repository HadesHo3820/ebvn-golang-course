package user

import (
	"errors"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
)

// registerInputBody represents the expected JSON request body for user registration.
// All fields are required for successful registration.
type registerInputBody struct {
	// Username is the unique login identifier for the user.
	Username string `json:"username" validate:"required,gte=2,lte=20"`
	// Password must contain at least one uppercase, one lowercase, one digit, and one special character.
	Password string `json:"password" validate:"required,gte=8,lte=20,password"`
	// DisplayName is the user's display name shown in the UI.
	DisplayName string `json:"display_name" validate:"required,gte=2,lte=50"`
	// Email is the user's email address.
	Email string `json:"email" validate:"required,email"`
}

// registerUserData is a custom DTO for the registration response.
// It excludes CreatedAt to control exactly which fields are returned.
type registerUserData struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	UpdatedAt   string `json:"updated_at"`
}

type registerResBody struct {
	Data    *registerUserData `json:"data"`
	Message string            `json:"message"`
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
// @Success 200 {object} registerResBody
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Router /v1/users/register [post]
func (u *userHandler) Register(c *gin.Context) {
	inputBody, err := utils.BindInputFromRequest[registerInputBody](c)
	if err != nil {
		return
	}

	// call service
	res, err := u.svc.CreateUser(c, inputBody.Username, inputBody.Password, inputBody.DisplayName, inputBody.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationType):
		c.JSON(http.StatusBadRequest, response.Message{
			Message: "username or email is already taken",
		})
		return
	case errors.Is(err, nil):
	default:
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, &registerResBody{
		Message: "Register an user successfully!",
		Data: &registerUserData{
			ID:          res.ID,
			Username:    res.Username,
			DisplayName: res.DisplayName,
			Email:       res.Email,
			UpdatedAt:   res.UpdatedAt.String(),
		},
	})
}
