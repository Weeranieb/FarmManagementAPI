package handler

import (
	"fmt"
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

type FarmGroupHandler interface {
	AddFarmGroup(c *fiber.Ctx) error
	GetFarmGroup(c *fiber.Ctx) error
	UpdateFarmGroup(c *fiber.Ctx) error
	ListFarmGroup(c *fiber.Ctx) error
	GetFarmGroupDropdown(c *fiber.Ctx) error
}

type farmGroupHandlerImpl struct {
	farmGroupService service.FarmGroupService
}

func NewFarmGroupHandler(farmGroupService service.FarmGroupService) FarmGroupHandler {
	return &farmGroupHandlerImpl{
		farmGroupService: farmGroupService,
	}
}

func (h *farmGroupHandlerImpl) AddFarmGroup(c *fiber.Ctx) error {
	var request dto.CreateFarmGroupRequest

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	isAdmin, err := utils.IsClientAdminOrAbove(c.UserContext())
	if err != nil || !isAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	if err := validateClientAccess(c, request.ClientId); err != nil {
		return err
	}

	result, err := h.farmGroupService.Create(c.UserContext(), request)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

func (h *farmGroupHandlerImpl) GetFarmGroup(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid farm group ID")
	}

	result, err := h.farmGroupService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	canAccess, accessErr := utils.CanAccessClient(c.UserContext(), result.ClientId)
	if accessErr != nil || !canAccess {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	return http.Success(c, result)
}

func (h *farmGroupHandlerImpl) UpdateFarmGroup(c *fiber.Ctx) error {
	var request dto.UpdateFarmGroupRequest

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	isAdmin, err := utils.IsClientAdminOrAbove(c.UserContext())
	if err != nil || !isAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	existing, err := h.farmGroupService.Get(request.Id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	canAccess, accessErr := utils.CanAccessClient(c.UserContext(), existing.ClientId)
	if accessErr != nil || !canAccess {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	if err := h.farmGroupService.Update(c.UserContext(), request); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func (h *farmGroupHandlerImpl) ListFarmGroup(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	var clientId int

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if isSuperAdmin {
		clientIdStr := c.Query("clientId")
		if clientIdStr != "" {
			clientIdVal, err := strconv.Atoi(clientIdStr)
			if err != nil {
				return http.Error(c, errors.ErrValidationFailed.Code, "Invalid clientId parameter")
			}
			clientId = clientIdVal
		}
	} else {
		clientIdPtr := utils.GetClientId(c.UserContext())
		if clientIdPtr == nil {
			return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
		}
		clientId = *clientIdPtr
	}

	list, err := h.farmGroupService.List(clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, list)
}

func (h *farmGroupHandlerImpl) GetFarmGroupDropdown(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	var clientId int

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	if isSuperAdmin {
		clientIdStr := c.Query("clientId")
		if clientIdStr != "" {
			clientIdVal, err := strconv.Atoi(clientIdStr)
			if err != nil {
				return http.Error(c, errors.ErrValidationFailed.Code, "Invalid clientId parameter")
			}
			clientId = clientIdVal
		}
	} else {
		clientIdPtr := utils.GetClientId(c.UserContext())
		if clientIdPtr == nil {
			return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
		}
		clientId = *clientIdPtr
	}

	items, err := h.farmGroupService.GetDropdown(clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, items)
}
