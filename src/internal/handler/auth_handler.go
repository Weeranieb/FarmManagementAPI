package handler

import (
	"fmt"

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
}

type authHandlerImpl struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandlerImpl{
		authService: authService,
	}
}

// POST /api/v1/auth/register
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
// @Router       /api/v1/auth/register [post]
func (h *authHandlerImpl) Register(c *fiber.Ctx) error {
	var registerRequest dto.RegisterRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &registerRequest, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	newUser, err := h.authService.Register(registerRequest)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newUser)
}

// POST /api/v1/auth/login
// Login user and return JWT token.
// @Summary      Login user
// @Description  Login user with provided credentials and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginRequest true "Login data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      401  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/auth/login [post]
func (h *authHandlerImpl) Login(c *fiber.Ctx) error {
	var loginRequest dto.LoginRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &loginRequest, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	token, user, expDate, err := h.authService.Login(loginRequest)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	loginResponse := dto.LoginResponse{
		AccessToken: token,
		ExpiredAt:   expDate,
		User:        user,
	}

	return http.Success(c, loginResponse)
}
