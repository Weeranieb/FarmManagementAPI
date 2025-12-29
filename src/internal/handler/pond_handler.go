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

//go:generate go run github.com/vektra/mockery/v2@latest --name=PondHandler --output=./mocks --outpkg=handler --filename=pond_handler.go --structname=MockPondHandler --with-expecter=false
type PondHandler interface {
	AddPond(c *fiber.Ctx) error
	AddPonds(c *fiber.Ctx) error
	GetPond(c *fiber.Ctx) error
	GetPondList(c *fiber.Ctx) error
	UpdatePond(c *fiber.Ctx) error
	DeletePond(c *fiber.Ctx) error
}

type pondHandlerImpl struct {
	pondService service.PondService
}

func NewPondHandler(pondService service.PondService) PondHandler {
	return &pondHandlerImpl{
		pondService: pondService,
	}
}

// POST /api/v1/pond
// Add a new pond.
// @Summary      Add a new pond
// @Description  Create a new pond with the provided details
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body dto.CreatePondRequest true "Pond data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/pond [post]
func (h *pondHandlerImpl) AddPond(c *fiber.Ctx) error {
	var createPondRequest dto.CreatePondRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createPondRequest, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	// Get username
	username, err := utils.GetUsername(c)
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	newPond, err := h.pondService.Create(createPondRequest, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newPond)
}

// POST /api/v1/pond/batch
// Add multiple ponds.
// @Summary      Add multiple ponds
// @Description  Create multiple ponds with the provided details
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body dto.CreatePondsRequest true "Ponds data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/pond/batch [post]
func (h *pondHandlerImpl) AddPonds(c *fiber.Ctx) error {
	var createPondsRequest dto.CreatePondsRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createPondsRequest, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	// Get username
	username, err := utils.GetUsername(c)
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	newPonds, err := h.pondService.CreateBatch(createPondsRequest, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newPonds)
}

// GET /api/v1/pond/:id
// Get a pond by ID.
// @Summary      Get a pond by ID
// @Description  Retrieve a pond by its ID
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        id path int true "Pond ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/pond/{id} [get]
func (h *pondHandlerImpl) GetPond(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get id from param
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	pond, err := h.pondService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, pond)
}

// GET /api/v1/pond
// Get a list of ponds by farm ID.
// @Summary      Get a list of ponds by farm ID
// @Description  Retrieve a list of ponds belonging to a specific farm
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        farmId query int true "Farm ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/pond [get]
func (h *pondHandlerImpl) GetPondList(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get farmId from query
	farmIdStr := c.Query("farmId")
	farmId, err := strconv.Atoi(farmIdStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid farm ID")
	}

	pondList, err := h.pondService.GetList(farmId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, pondList)
}

// PUT /api/v1/pond
// Update a pond.
// @Summary      Update a pond
// @Description  Update an existing pond with new details
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body model.Pond true "Updated pond data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/pond [put]
func (h *pondHandlerImpl) UpdatePond(c *fiber.Ctx) error {
	var updatePond *model.Pond

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := c.BodyParser(&updatePond); err != nil {
		return http.Error(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Message)
	}

	// Get username
	username, err := utils.GetUsername(c)
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.pondService.Update(updatePond, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// DELETE /api/v1/pond/:id
// Delete a pond.
// @Summary      Delete a pond
// @Description  Delete a pond by its ID
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path int true "Pond ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/pond/{id} [delete]
func (h *pondHandlerImpl) DeletePond(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get id from param
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	// Get username
	username, err := utils.GetUsername(c)
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.pondService.Delete(id, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
