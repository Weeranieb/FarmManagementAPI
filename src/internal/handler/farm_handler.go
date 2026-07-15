package handler

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

type FarmHandler interface {
	AddFarm(c *fiber.Ctx) error
	GetFarm(c *fiber.Ctx) error
	GetFarmList(c *fiber.Ctx) error
	GetFarmHierarchy(c *fiber.Ctx) error
	UpdateFarm(c *fiber.Ctx) error
}

type farmHandlerImpl struct {
	farmService service.FarmService
}

func NewFarmHandler(farmService service.FarmService) FarmHandler {
	return &farmHandlerImpl{
		farmService: farmService,
	}
}

// Add a new farm entry. Allowed for super admin (any client) or client admin (their own client only).
// @Summary      Add a new farm entry
// @Description  Add a new farm entry with the provided details. Requires client-admin role or above.
// @Tags         farm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.CreateFarmRequest true "Farm data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm [post]
func (h *farmHandlerImpl) AddFarm(c *fiber.Ctx) error {
	var createFarmRequest dto.CreateFarmRequest

	if err := validateAndParse(c, &createFarmRequest); err != nil {
		return err
	}

	if err := requireClientAdmin(c); err != nil {
		return err
	}
	if err := requireClientAccess(c, createFarmRequest.ClientId); err != nil {
		return err
	}

	newFarm, err := h.farmService.Create(c.UserContext(), createFarmRequest, createFarmRequest.ClientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newFarm)
}

// Get farm by ID.
// @Summary      Get farm by ID
// @Description  Retrieve details of a specific farm by its ID
// @Tags         farm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "Farm ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm/{id} [get]
func (h *farmHandlerImpl) GetFarm(c *fiber.Ctx) error {

	id, err := parseParamInt(c, "id", "Invalid farm ID")
	if err != nil {
		return err
	}

	clientIdPtr, canAccess := utils.GetClientIdForAccess(c.UserContext())
	if !canAccess {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	farm, err := h.farmService.Get(id, clientIdPtr)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, farm)
}

// Get list of farms associated with the current client.
// @Summary      Get list of farms
// @Description  Retrieve a list of farms associated with the current client
// @Tags         farm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        clientId query int false "Client ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm [get]
func (h *farmHandlerImpl) GetFarmList(c *fiber.Ctx) error {

	clientId, err := resolveListClientId(c)
	if err != nil {
		return err
	}

	farmList, err := h.farmService.GetList(clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, farmList)
}

// Get farms with nested ponds for the current client (Existing Data view).
// @Summary      Get farm hierarchy with ponds
// @Description  Retrieve all farms for the client with their nested ponds (for Existing Data view). Super admin may pass clientId query param.
// @Tags         farm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        clientId query int false "Client ID (optional for super admin)"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm/hierarchy [get]
func (h *farmHandlerImpl) GetFarmHierarchy(c *fiber.Ctx) error {

	clientId, err := resolveListClientId(c)
	if err != nil {
		return err
	}

	list, err := h.farmService.GetHierarchy(clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, list)
}

// Update farm entry. Super admin only.
// @Summary      Update farm entry
// @Description  Update details of a farm entry. Super admin only. Id in path; body contains name.
// @Tags         farm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id   path int true "Farm ID"
// @Param        body body dto.UpdateFarmBody true "Farm data (name)"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm/{id} [put]
func (h *farmHandlerImpl) UpdateFarm(c *fiber.Ctx) error {
	id, err := parseParamInt(c, "id", "Invalid farm ID")
	if err != nil {
		return err
	}

	var body dto.UpdateFarmBody
	if err := validateAndParse(c, &body); err != nil {
		return err
	}

	if err := requireSuperAdmin(c); err != nil {
		return err
	}

	updateReq := dto.UpdateFarmRequest{Id: id, Name: body.Name}
	if err = h.farmService.Update(c.UserContext(), updateReq); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
