package handler

import (
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v3"
)

type UserHandler interface {
	AddUser(c fiber.Ctx) error
	GetUser(c fiber.Ctx) error
	UpdateUser(c fiber.Ctx) error
	AdminUpdateUser(c fiber.Ctx) error
	AdminResetPassword(c fiber.Ctx) error
	ChangePassword(c fiber.Ctx) error
	DeleteUser(c fiber.Ctx) error
	GetUserList(c fiber.Ctx) error
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
func (h *userHandlerImpl) AddUser(c fiber.Ctx) error {
	var addUser dto.CreateUserRequest

	if err := validateAndParse(c, &addUser); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}
	clientId := utils.GetClientId(c.Context())

	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	newUser, err := h.userService.Create(c.Context(), addUser, username, clientId)
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
func (h *userHandlerImpl) GetUser(c fiber.Ctx) error {

	id, err := utils.GetUserId(c.Context())
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
// Update the current user (self-update).
// @Summary      Update the current user
// @Description  Update the details of the currently authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.UpdateUserRequest true "Updated user data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user [put]
func (h *userHandlerImpl) UpdateUser(c fiber.Ctx) error {
	var updateUser dto.UpdateUserRequest

	if err := validateAndParse(c, &updateUser); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrGeneric.Code, errors.ErrGeneric.Message)
	}

	userId, err := utils.GetUserId(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	updated, err := h.userService.Update(c.Context(), userId, updateUser, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, updated)
}

// PUT /api/v1/user/:id
// Super-admin update of any user.
// @Summary      Admin-update a user
// @Description  Update any user. Super-admin only. Cannot promote to SuperAdmin or modify an existing SuperAdmin.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "User ID"
// @Param        body body dto.AdminUpdateUserRequest true "Updated user data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user/{id} [put]
func (h *userHandlerImpl) AdminUpdateUser(c fiber.Ctx) error {
	var body dto.AdminUpdateUserRequest

	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	userId, err := parseParamInt(c, "id", "Invalid user ID")
	if err != nil {
		return err
	}

	if err := validateAndParse(c, &body); err != nil {
		return err
	}

	actor, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if err := h.userService.AdminUpdate(c.Context(), userId, body, actor); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// PUT /api/v1/user/:id/password
// Super-admin reset of another user's password.
// @Summary      Admin-reset a user's password
// @Description  Reset another user's password. Super-admin only. Cannot reset a SuperAdmin.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "User ID"
// @Param        body body dto.AdminResetPasswordRequest true "New password"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user/{id}/password [put]
func (h *userHandlerImpl) AdminResetPassword(c fiber.Ctx) error {
	var body dto.AdminResetPasswordRequest

	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	userId, err := parseParamInt(c, "id", "Invalid user ID")
	if err != nil {
		return err
	}

	if err := validateAndParse(c, &body); err != nil {
		return err
	}

	actor, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if err := h.userService.AdminResetPassword(c.Context(), userId, body, actor); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// PUT /api/v1/user/password
// Change the current user's own password.
// @Summary      Change own password
// @Description  Authenticated user changes their own password. Requires current password.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.ChangePasswordRequest true "Current and new password"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      401  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user/password [put]
func (h *userHandlerImpl) ChangePassword(c fiber.Ctx) error {
	var body dto.ChangePasswordRequest

	userId, err := utils.GetUserId(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if err := validateAndParse(c, &body); err != nil {
		return err
	}

	actor, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if err := h.userService.ChangePassword(c.Context(), userId, body, actor); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// DELETE /api/v1/user/:id
// Soft-delete a user. Super-admin only.
// @Summary      Delete a user
// @Description  Soft-delete a user. Super-admin only. Cannot delete self or another super admin.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "User ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user/{id} [delete]
func (h *userHandlerImpl) DeleteUser(c fiber.Ctx) error {

	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	userId, err := parseParamInt(c, "id", "Invalid user ID")
	if err != nil {
		return err
	}

	actor, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if err := h.userService.Delete(c.Context(), userId, actor); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// GET /api/v1/user/list
// Get a list of users.
// @Summary      Get a list of users
// @Description  Retrieve users. Super admins see all users and can filter; non-admins are scoped to their own client.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        search query string false "Substring match across username, email, firstName, lastName"
// @Param        userLevel query int false "Filter by user level"
// @Param        clientId query int false "Filter by client id (super admin only)"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /user/list [get]
func (h *userHandlerImpl) GetUserList(c fiber.Ctx) error {

	isSuperAdmin, err := utils.IsSuperAdmin(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	filters := dto.UserListQuery{}

	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}
	if levelStr := c.Query("userLevel"); levelStr != "" {
		if level, parseErr := strconv.Atoi(levelStr); parseErr == nil {
			filters.UserLevel = &level
		}
	}
	if clientStr := c.Query("clientId"); clientStr != "" && isSuperAdmin {
		if cid, parseErr := strconv.Atoi(clientStr); parseErr == nil {
			filters.ClientId = &cid
		}
	}

	// Non-super-admins are forcibly scoped to their own client, ignoring any query param.
	if !isSuperAdmin {
		filters.ClientId = utils.GetClientId(c.Context())
	}

	users, err := h.userService.GetUserList(c.Context(), filters)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, users)
}
