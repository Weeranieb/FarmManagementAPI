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

//go:generate go run github.com/vektra/mockery/v2@latest --name=WorkerHandler --output=./mocks --outpkg=handler --filename=worker_handler.go --structname=MockWorkerHandler --with-expecter=false
type WorkerHandler interface {
	AddWorker(c *fiber.Ctx) error
	GetWorker(c *fiber.Ctx) error
	UpdateWorker(c *fiber.Ctx) error
	ListWorker(c *fiber.Ctx) error
}

type workerHandlerImpl struct {
	workerService service.WorkerService
}

func NewWorkerHandler(workerService service.WorkerService) WorkerHandler {
	return &workerHandlerImpl{
		workerService: workerService,
	}
}

// POST /worker
// Add a new worker.
// @Summary      Add a new worker
// @Description  Add a new worker with the provided details
// @Tags         worker
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.CreateWorkerRequest true "Worker data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /worker [post]
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

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	// Get client id
	clientIdPtr := utils.GetClientId(c.UserContext())
	if clientIdPtr == nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	newWorker, err := h.workerService.Create(createWorkerRequest, username, *clientIdPtr)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newWorker)
}

// GET /worker/:id
// Get a worker by ID.
// @Summary      Get a worker by ID
// @Description  Retrieve a worker by its ID
// @Tags         worker
// @Accept       json
// @Produce      json
// @Param        id path int true "Worker ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /worker/{id} [get]
func (h *workerHandlerImpl) GetWorker(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get id from param
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid worker ID")
	}

	worker, err := h.workerService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, worker)
}

// PUT /worker
// Update a worker.
// @Summary      Update a worker
// @Description  Update an existing worker with new details
// @Tags         worker
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body model.Worker true "Updated worker data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /worker [put]
func (h *workerHandlerImpl) UpdateWorker(c *fiber.Ctx) error {
	var updateWorker *model.Worker

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &updateWorker); err != nil {
		return err
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.workerService.Update(updateWorker, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// GET /worker
// Get a list of workers with pagination.
// @Summary      Get a list of workers with pagination
// @Description  Retrieve a paginated list of workers for the current client
// @Tags         worker
// @Accept       json
// @Produce      json
// @Param        page query int true "Page number"
// @Param        pageSize query int true "Page size"
// @Param        orderBy query string false "Order by field"
// @Param        keyword query string false "Search keyword"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /worker [get]
func (h *workerHandlerImpl) ListWorker(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get query parameters
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

	// Get client id
	clientIdPtr := utils.GetClientId(c.UserContext())
	if clientIdPtr == nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	workerList, err := h.workerService.GetPage(*clientIdPtr, page, pageSize, orderBy, keyword)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, workerList)
}
