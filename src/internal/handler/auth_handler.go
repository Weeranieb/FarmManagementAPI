package handler

import (
	"fmt"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=AuthHandler --output=./mocks --outpkg=handler --filename=auth_handler.go --structname=MockAuthHandler --with-expecter=false
type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authHandlerImpl struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandlerImpl{
		authService: authService,
	}
}

// POST /auth/register
// Register a new user.
// @Summary      Register a new user
// @Description  Register a new user with the provided details
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RegisterRequest true "User data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /auth/register [post]
func (h *authHandlerImpl) Register(c *fiber.Ctx) error {
	var registerRequest dto.RegisterRequest

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &registerRequest); err != nil {
		return err
	}

	newUser, err := h.authService.Register(registerRequest)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newUser)
}

// POST /auth/login
// Login user and return JWT token.
// @Summary      Login user
// @Description  Login user with provided credentials and return JWT token. Token is also set as HTTP-only cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginRequest true "Login data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      401  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /auth/login [post]
func (h *authHandlerImpl) Login(c *fiber.Ctx) error {
	var loginRequest dto.LoginRequest

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &loginRequest); err != nil {
		return err
	}

	token, user, expDate, err := h.authService.Login(loginRequest)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	// Set JWT token as HTTP-only cookie
	cookie := &fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Protocol() == "https", // Only send over HTTPS in production
		SameSite: "Strict",
	}

	if expDate != nil {
		cookie.Expires = *expDate
		cookie.MaxAge = int(time.Until(*expDate).Seconds())
	} else {
		// Default to 24 hours if no expiration date
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.MaxAge = 86400 // 24 hours in seconds
	}

	c.Cookie(cookie)

	loginResponse := dto.LoginResponse{
		AccessToken: token,
		ExpiredAt:   expDate,
		User:        user,
	}

	return http.Success(c, loginResponse)
}

// POST /auth/logout
// Logout user and clear JWT token cookie.
// @Summary      Logout user
// @Description  Logout user and clear the authentication cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  http.ResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /auth/logout [post]
func (h *authHandlerImpl) Logout(c *fiber.Ctx) error {
	// Clear the JWT token cookie
	cookie := &fiber.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Strict",
		Expires:  time.Now().Add(-1 * time.Hour), // Set to past date to delete
		MaxAge:   -1,
	}

	c.Cookie(cookie)

	return http.SuccessWithoutData(c)
}
