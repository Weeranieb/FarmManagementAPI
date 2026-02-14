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

// POST /pond
// Create multiple ponds for a farm (farmId, names array). New ponds are created with status maintenance.
// @Summary      Create multiple ponds
// @Description  Create multiple ponds for a farm. Request: farmId, array of names. New ponds have status maintenance.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.CreatePondsRequest true "farmId, names[]"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond [post]
func (h *pondHandlerImpl) AddPonds(c *fiber.Ctx) error {
	var createPondsRequest dto.CreatePondsRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createPondsRequest); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil || !isSuperAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	if err := h.pondService.CreatePonds(c.UserContext(), createPondsRequest, username); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, nil)
}

// GET /pond/:id
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
// @Router       /pond/{id} [get]
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

// GET /pond
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
// @Router       /pond [get]
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

// PUT /pond
// Update a pond.
// @Summary      Update a pond
// @Description  Update an existing pond with new details
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body model.Pond true "Updated pond data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond [put]
func (h *pondHandlerImpl) UpdatePond(c *fiber.Ctx) error {
	var updatePond *model.Pond

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &updatePond); err != nil {
		return err
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.pondService.Update(c.UserContext(), updatePond, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// DELETE /pond/:id
// Delete a pond.
// @Summary      Delete a pond
// @Description  Delete a pond by its ID
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "Pond ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{id} [delete]
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
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.pondService.Delete(id, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
