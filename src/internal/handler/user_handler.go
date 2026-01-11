package handler

import (
	"fmt"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
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
// @Description  Create a new user with the provided details. For system setup, this endpoint is public and allows creating users without authentication. When authenticated, uses JWT context for creator and clientId.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "Bearer token (optional for system setup)"
// @Param        body body dto.CreateUserRequest true "User data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/user [post]
func (h *userHandlerImpl) AddUser(c *fiber.Ctx) error {
	var addUser dto.CreateUserRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &addUser, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	var username string
	var clientId *int

	// Try to get username and clientId from JWT context (for authenticated requests)
	jwtUsername, jwtErr := utils.GetUsername(c)
	jwtClientId, jwtClientErr := utils.GetClientId(c)

	if jwtErr == nil && jwtClientErr == nil {
		// Authenticated request - use JWT context
		username = jwtUsername
		clientId = &jwtClientId
	} else {
		// System setup - bypass authentication
		// Use "system" as the creator for setup operations
		username = "system"

		// For system setup, use ClientId from request if provided, otherwise allow NULL
		// (Super admin can be created without a client during initial setup)
		if addUser.ClientId != nil {
			clientId = addUser.ClientId
		} else {
			// Allow NULL clientId for system setup (super admin)
			clientId = nil
		}
	}

	newUser, err := h.userService.Create(addUser, username, clientId)
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
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/user [get]
func (h *userHandlerImpl) GetUser(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// get userId
	id, err := utils.GetUserId(c)
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
// @Param        body body model.User true "Updated user data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/user [put]
func (h *userHandlerImpl) UpdateUser(c *fiber.Ctx) error {
	var updateUser *model.User

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := c.BodyParser(&updateUser); err != nil {
		return http.Error(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Message)
	}

	// get username
	username, err := utils.GetUsername(c)
	if err != nil {
		return http.Error(c, errors.ErrGeneric.Code, errors.ErrGeneric.Message)
	}

	err = h.userService.Update(updateUser, username)
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
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/user/list [get]
func (h *userHandlerImpl) GetUserList(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// get clientId
	clientId, err := utils.GetClientId(c)
	if err != nil {
		return http.Error(c, errors.ErrGeneric.Code, errors.ErrGeneric.Message)
	}

	users, err := h.userService.GetUserList(clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, users)
}
