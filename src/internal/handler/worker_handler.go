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

type WorkerHandler interface {
	AddWorker(c fiber.Ctx) error
	GetWorker(c fiber.Ctx) error
	UpdateWorker(c fiber.Ctx) error
	ListWorker(c fiber.Ctx) error
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

func (h *workerHandlerImpl) AddWorker(c fiber.Ctx) error {
	var createWorkerRequest dto.CreateWorkerRequest

	if err := validateAndParse(c, &createWorkerRequest); err != nil {
		return err
	}

	if err := requireClientAdmin(c); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	clientId, err := h.resolveClientId(c, createWorkerRequest.FarmGroupId)
	if err != nil {
		return err
	}

	newWorker, err := h.workerService.Create(c.Context(), createWorkerRequest, username, clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newWorker)
}

func (h *workerHandlerImpl) GetWorker(c fiber.Ctx) error {
	id, err := parseParamInt(c, "id", "Invalid worker ID")
	if err != nil {
		return err
	}

	worker, err := h.workerService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	if err := requireClientAccess(c, worker.ClientId); err != nil {
		return err
	}

	return http.Success(c, worker)
}

func (h *workerHandlerImpl) UpdateWorker(c fiber.Ctx) error {
	var updateWorker dto.UpdateWorkerRequest

	if err := validateAndParse(c, &updateWorker); err != nil {
		return err
	}

	if err := requireClientAdmin(c); err != nil {
		return err
	}

	existing, getErr := h.workerService.Get(updateWorker.Id)
	if getErr != nil {
		return http.NewError(c, errors.ErrGeneric.Code, getErr)
	}

	if err := requireClientAccess(c, existing.ClientId); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.workerService.Update(c.Context(), updateWorker, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func (h *workerHandlerImpl) ListWorker(c fiber.Ctx) error {
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

	clientId, err := resolveListClientId(c)
	if err != nil {
		return err
	}

	workerList, err := h.workerService.GetPage(clientId, page, pageSize, orderBy, keyword)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, workerList)
}

// resolveClientId derives the clientId: from the JWT token for regular users,
// or from the farmGroupId for super admins who may not have clientId in their token.
func (h *workerHandlerImpl) resolveClientId(c fiber.Ctx, farmGroupId int) (int, error) {
	clientIdPtr := utils.GetClientId(c.Context())
	if clientIdPtr != nil {
		return *clientIdPtr, nil
	}

	isSuperAdmin, err := utils.IsSuperAdmin(c.Context())
	if err != nil || !isSuperAdmin {
		return 0, http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	clientId, fgErr := h.farmGroupService.GetClientIdByFarmGroupId(farmGroupId)
	if fgErr != nil {
		return 0, http.NewError(c, errors.ErrGeneric.Code, fgErr)
	}
	return clientId, nil
}
