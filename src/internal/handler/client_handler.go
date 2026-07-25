package handler

import (

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v3"
)

type ClientHandler interface {
	AddClient(c fiber.Ctx) error
	GetClient(c fiber.Ctx) error
	GetClientList(c fiber.Ctx) error
	GetClientSummaries(c fiber.Ctx) error
	UpdateClient(c fiber.Ctx) error
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
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.CreateClientRequest true "Client data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client [post]
func (h *clientHandlerImpl) AddClient(c fiber.Ctx) error {
	var createClientRequest dto.CreateClientRequest

	if err := validateAndParse(c, &createClientRequest); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	newClient, err := h.clientService.Create(c.Context(), createClientRequest, username)
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
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "Client ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client/{id} [get]
func (h *clientHandlerImpl) GetClient(c fiber.Ctx) error {

	id, err := parseParamInt(c, "id", "Invalid client ID")
	if err != nil {
		return err
	}

	// Check if user can access this client
	if err := requireClientAccess(c, id); err != nil {
		return err
	}

	client, err := h.clientService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, client)
}

// GET /client/list
// Get list of clients for dropdown (id + name). Super admin only.
// @Summary      Get client list for dropdown
// @Description  Returns a list of clients with id and name for dropdown/select. Only super admin can access.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  http.ResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client/list [get]
func (h *clientHandlerImpl) GetClientList(c fiber.Ctx) error {

	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	dropdown, err := h.clientService.GetClientDropdown()
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, dropdown)
}

// GET /client/summaries
// Get list of clients with aggregate counts (farms, ponds, users). Super admin only.
// @Summary      Get client summaries
// @Description  Returns each client with farmCount, pondCount, userCount. Super admin only.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  http.ResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client/summaries [get]
func (h *clientHandlerImpl) GetClientSummaries(c fiber.Ctx) error {

	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	summaries, err := h.clientService.GetSummaries()
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, summaries)
}

// PUT /client
// Update client.
// @Summary      Update client
// @Description  Update details of a client. Super admin can update any client, others can only update their own client.
// @Tags         client
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body model.Client true "Client data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /client [put]
func (h *clientHandlerImpl) UpdateClient(c fiber.Ctx) error {
	var updateClient dto.UpdateClientRequest

	if err := validateAndParse(c, &updateClient); err != nil {
		return err
	}

	// Check if user can access this client
	if err := requireClientAccess(c, updateClient.Id); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.clientService.Update(c.Context(), updateClient, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
