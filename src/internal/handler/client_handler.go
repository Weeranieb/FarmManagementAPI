package handler

import (
	"fmt"
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=ClientHandler --output=./mocks --outpkg=handler --filename=client_handler.go --structname=MockClientHandler --with-expecter=false
type ClientHandler interface {
	AddClient(c *fiber.Ctx) error
	GetClient(c *fiber.Ctx) error
	UpdateClient(c *fiber.Ctx) error
}

type clientHandlerImpl struct {
	clientService service.ClientService
}

func NewClientHandler(clientService service.ClientService) ClientHandler {
	return &clientHandlerImpl{
		clientService: clientService,
	}
}

// POST /client
// Add a new client.
// @Summary      Add a new client
// @Description  Add a new client with the provided details. Only super admin can create clients.
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body dto.CreateClientRequest true "Client data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client [post]
func (h *clientHandlerImpl) AddClient(c *fiber.Ctx) error {
	var createClientRequest dto.CreateClientRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Check if user is super admin (skip check for system setup - no JWT)
	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil {
		// No JWT token - this is system setup, allow it
		// Use "system" as the creator
	} else if !isSuperAdmin {
		// Has JWT but not super admin - deny
		return http.Error(c, 403, "Only super admin can create clients")
	}

	if err := validateAndParse(c, &createClientRequest, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	// Get username (for system setup, use "system" if no JWT)
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	newClient, err := h.clientService.Create(createClientRequest, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newClient)
}

// GET /client/:id
// Get a client by ID.
// @Summary      Get a client by ID
// @Description  Retrieve a client by its ID. Super admin can access any client, others can only access their own client.
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path int true "Client ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client/{id} [get]
func (h *clientHandlerImpl) GetClient(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get id from param
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid client ID")
	}

	// Check if user can access this client
	canAccess, err := utils.CanAccessClient(c.UserContext(), id)
	if err != nil || !canAccess {
		return http.Error(c, 403, "Access denied to this client")
	}

	client, err := h.clientService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, client)
}

// PUT /client
// Update client.
// @Summary      Update client
// @Description  Update details of a client. Super admin can update any client, others can only update their own client.
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body model.Client true "Client data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client [put]
func (h *clientHandlerImpl) UpdateClient(c *fiber.Ctx) error {
	var updateClient *model.Client

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := c.BodyParser(&updateClient); err != nil {
		return http.Error(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Message)
	}

	// Check if user can access this client
	canAccess, err := utils.CanAccessClient(c.UserContext(), updateClient.Id)
	if err != nil || !canAccess {
		return http.Error(c, 403, "Access denied to this client")
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.clientService.Update(updateClient, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
