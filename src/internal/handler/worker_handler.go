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

//go:generate go run github.com/vektra/mockery/v2@latest --name=WorkerHandler --output=./mocks --outpkg=handler --filename=worker_handler.go --structname=MockWorkerHandler --with-expecter=false
type WorkerHandler interface {
	AddWorker(c *fiber.Ctx) error
	GetWorker(c *fiber.Ctx) error
	UpdateWorker(c *fiber.Ctx) error
	ListWorker(c *fiber.Ctx) error
}

type workerHandlerImpl struct {
	workerService    service.WorkerService
	farmGroupService service.FarmGroupService
}

func NewWorkerHandler(workerService service.WorkerService, farmGroupService service.FarmGroupService) WorkerHandler {
	return &workerHandlerImpl{
		workerService:    workerService,
		farmGroupService: farmGroupService,
	}
}

func (h *workerHandlerImpl) AddWorker(c *fiber.Ctx) error {
	var createWorkerRequest dto.CreateWorkerRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createWorkerRequest); err != nil {
		return err
	}

	isAdmin, err := utils.IsClientAdminOrAbove(c.UserContext())
	if err != nil || !isAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	clientId, err := h.resolveClientId(c, createWorkerRequest.FarmGroupId)
	if err != nil {
		return err
	}

	newWorker, err := h.workerService.Create(c.UserContext(), createWorkerRequest, username, clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newWorker)
}

func (h *workerHandlerImpl) GetWorker(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid worker ID")
	}

	worker, err := h.workerService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	canAccess, accessErr := utils.CanAccessClient(c.UserContext(), worker.ClientId)
	if accessErr != nil || !canAccess {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	return http.Success(c, worker)
}

func (h *workerHandlerImpl) UpdateWorker(c *fiber.Ctx) error {
	var updateWorker dto.UpdateWorkerRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &updateWorker); err != nil {
		return err
	}

	isAdmin, err := utils.IsClientAdminOrAbove(c.UserContext())
	if err != nil || !isAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	existing, getErr := h.workerService.Get(updateWorker.Id)
	if getErr != nil {
		return http.NewError(c, errors.ErrGeneric.Code, getErr)
	}

	canAccess, accessErr := utils.CanAccessClient(c.UserContext(), existing.ClientId)
	if accessErr != nil || !canAccess {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.workerService.Update(c.UserContext(), updateWorker, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func (h *workerHandlerImpl) ListWorker(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	sPage := c.Query("page")
	sPageSize := c.Query("pageSize")
	orderBy := c.Query("orderBy")
	keyword := c.Query("keyword")

	page, err := strconv.Atoi(sPage)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid page number")
	}

	pageSize, err := strconv.Atoi(sPageSize)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid page size")
	}

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

	workerList, err := h.workerService.GetPage(clientId, page, pageSize, orderBy, keyword)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, workerList)
}

// resolveClientId derives the clientId: from the JWT token for regular users,
// or from the farmGroupId for super admins who may not have clientId in their token.
func (h *workerHandlerImpl) resolveClientId(c *fiber.Ctx, farmGroupId int) (int, error) {
	clientIdPtr := utils.GetClientId(c.UserContext())
	if clientIdPtr != nil {
		return *clientIdPtr, nil
	}

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil || !isSuperAdmin {
		return 0, http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	clientId, fgErr := h.farmGroupService.GetClientIdByFarmGroupId(farmGroupId)
	if fgErr != nil {
		return 0, http.NewError(c, errors.ErrGeneric.Code, fgErr)
	}
	return clientId, nil
}
