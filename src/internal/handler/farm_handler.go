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

//go:generate go run github.com/vektra/mockery/v2@latest --name=FarmHandler --output=./mocks --outpkg=handler --filename=farm_handler.go --structname=MockFarmHandler --with-expecter=false
type FarmHandler interface {
	AddFarm(c *fiber.Ctx) error
	GetFarm(c *fiber.Ctx) error
	GetFarmList(c *fiber.Ctx) error
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

// POST /farm
// Add a new farm entry.
// @Summary      Add a new farm entry
// @Description  Add a new farm entry with the provided details
// @Tags         farm
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body dto.CreateFarmRequest true "Farm data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm [post]
func (h *farmHandlerImpl) AddFarm(c *fiber.Ctx) error {
	var createFarmRequest dto.CreateFarmRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createFarmRequest); err != nil {
		return err
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	// Get clientId
	clientIdPtr, err := utils.GetClientId(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}
	if clientIdPtr == nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	newFarm, err := h.farmService.Create(createFarmRequest, username, *clientIdPtr)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newFarm)
}

// GET /farm/:id
// Get farm by ID.
// @Summary      Get farm by ID
// @Description  Retrieve details of a specific farm by its ID
// @Tags         farm
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path int true "Farm ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm/{id} [get]
func (h *farmHandlerImpl) GetFarm(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get id from param
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid farm ID")
	}

	// Get clientId
	clientIdPtr, err := utils.GetClientId(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}
	if clientIdPtr == nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	farm, err := h.farmService.Get(id, *clientIdPtr)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, farm)
}

// GET /farm
// Get list of farms associated with the current client.
// @Summary      Get list of farms
// @Description  Retrieve a list of farms associated with the current client
// @Tags         farm
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm [get]
func (h *farmHandlerImpl) GetFarmList(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get clientId
	clientIdPtr, err := utils.GetClientId(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}
	if clientIdPtr == nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	farmList, err := h.farmService.GetList(*clientIdPtr)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, farmList)
}

// PUT /farm
// Update farm entry.
// @Summary      Update farm entry
// @Description  Update details of a farm entry
// @Tags         farm
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body model.Farm true "Farm data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /farm [put]
func (h *farmHandlerImpl) UpdateFarm(c *fiber.Ctx) error {
	var updateFarm *model.Farm

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := c.BodyParser(&updateFarm); err != nil {
		return http.Error(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Message)
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.farmService.Update(updateFarm, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
