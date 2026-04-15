package handler

import (
	"fmt"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserHandler --output=./mocks --outpkg=handler --filename=user_handler.go --structname=MockUserHandler --with-expecter=false
type UserHandler interface {
	AddUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	GetUserList(c *fiber.Ctx) error
}

type userHandlerImpl struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandlerImpl{
		userService: userService,
	}
}

// POST /api/v1/user
// Add a new user.
// @Summary      Add a new user
// @Description  Create a new user with the provided details. Only super admin can create users.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.CreateUserRequest true "User data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user [post]
func (h *userHandlerImpl) AddUser(c *fiber.Ctx) error {
	var addUser dto.CreateUserRequest

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &addUser); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}
	clientId := utils.GetClientId(c.UserContext())

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil || !isSuperAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	newUser, err := h.userService.Create(c.UserContext(), addUser, username, clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newUser)
}

// GET /api/v1/user
// Get the current user.
// @Summary      Get the current user
// @Description  Retrieve the user details of the currently authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user [get]
func (h *userHandlerImpl) GetUser(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// get userId
	id, err := utils.GetUserId(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrGeneric.Code, errors.ErrGeneric.Message)
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, user)
}

// PUT /api/v1/user
// Update the current user.
// @Summary      Update the current user
// @Description  Update the details of the currently authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body model.User true "Updated user data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user [put]
func (h *userHandlerImpl) UpdateUser(c *fiber.Ctx) error {
	var updateUser dto.UpdateUserRequest

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &updateUser); err != nil {
		return err
	}

	// get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrGeneric.Code, errors.ErrGeneric.Message)
	}

	userId, err := utils.GetUserId(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.userService.Update(c.UserContext(), userId, updateUser, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// GET /api/v1/user/list
// Get a list of users.
// @Summary      Get a list of users
// @Description  Retrieve a list of users associated with the current client ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user/list [get]
func (h *userHandlerImpl) GetUserList(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// get clientId
	var clientId *int
	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if !isSuperAdmin {
		clientId = utils.GetClientId(c.UserContext())
	}

	users, err := h.userService.GetUserList(c.UserContext(), clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, users)
}
