package user

import (
	"errors"
	"net/http"

	"github.com/HadesHo3820/ebvn-golang-course/internal/dto"
	"github.com/HadesHo3820/ebvn-golang-course/internal/handler/utils"
	"github.com/HadesHo3820/ebvn-golang-course/internal/service"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/dbutils"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/response"
	"github.com/gin-gonic/gin"
)

// loginInputBody represents the expected JSON request body for user login.
type loginInputBody struct {
	// Username is the user's login identifier.
	Username string `json:"username" validate:"required"`
	// Password is the user's password.
	Password string `json:"password" validate:"required,gte=8"`
}

// Login handles user login requests.
// It validates credentials, authenticates the user, and returns a JWT token.
//
// @Summary Login a user
// @Description Authenticate a user with username and password, returns a JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param body body loginInputBody true "User login credentials"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} response.Message "Invalid username or password"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/users/login [post]
func (u *userHandler) Login(c *gin.Context) {
	inputBody, err := utils.BindInputFromRequest[loginInputBody](c)
	if err != nil {
		return
	}

	// call service
	token, err := u.svc.Login(c, inputBody.Username, inputBody.Password)
	switch {
	case errors.Is(err, service.ErrClientErr):
		c.JSON(http.StatusBadRequest, response.Message{
			Message: err.Error(),
		})
		return
	case errors.Is(err, dbutils.ErrNotFoundType):
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid username or password"})
		return
	case errors.Is(err, nil):
	default:
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	// return token
	c.JSON(http.StatusOK, dto.SuccessResponse[string]{
		Message: "Logged in successfully!",
		Data:    token,
	})
}
