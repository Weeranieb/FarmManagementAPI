package handler

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v3"
)

type FarmGroupHandler interface {
	AddFarmGroup(c fiber.Ctx) error
	GetFarmGroup(c fiber.Ctx) error
	UpdateFarmGroup(c fiber.Ctx) error
	ListFarmGroup(c fiber.Ctx) error
	GetFarmGroupDropdown(c fiber.Ctx) error
}

type farmGroupHandlerImpl struct {
	farmGroupService service.FarmGroupService
}

func NewFarmGroupHandler(farmGroupService service.FarmGroupService) FarmGroupHandler {
	return &farmGroupHandlerImpl{
		farmGroupService: farmGroupService,
	}
}

func (h *farmGroupHandlerImpl) AddFarmGroup(c fiber.Ctx) error {
	var request dto.CreateFarmGroupRequest

	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	if err := requireClientAdmin(c); err != nil {
		return err
	}
	if err := requireClientAccess(c, request.ClientId); err != nil {
		return err
	}

	result, err := h.farmGroupService.Create(c.Context(), request)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

func (h *farmGroupHandlerImpl) GetFarmGroup(c fiber.Ctx) error {
	id, err := parseParamInt(c, "id", "Invalid farm group ID")
	if err != nil {
		return err
	}

	result, err := h.farmGroupService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	if err := requireClientAccess(c, result.ClientId); err != nil {
		return err
	}

	return http.Success(c, result)
}

func (h *farmGroupHandlerImpl) UpdateFarmGroup(c fiber.Ctx) error {
	var request dto.UpdateFarmGroupRequest

	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	if err := requireClientAdmin(c); err != nil {
		return err
	}

	existing, err := h.farmGroupService.Get(request.Id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	if err := requireClientAccess(c, existing.ClientId); err != nil {
		return err
	}

	if err := h.farmGroupService.Update(c.Context(), request); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func (h *farmGroupHandlerImpl) ListFarmGroup(c fiber.Ctx) error {
	clientId, err := resolveListClientId(c)
	if err != nil {
		return err
	}

	list, err := h.farmGroupService.List(clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, list)
}

func (h *farmGroupHandlerImpl) GetFarmGroupDropdown(c fiber.Ctx) error {
	clientId, err := resolveListClientId(c)
	if err != nil {
		return err
	}

	items, err := h.farmGroupService.GetDropdown(clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, items)
}
