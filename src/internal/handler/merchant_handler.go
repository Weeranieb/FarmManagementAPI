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

//go:generate go run github.com/vektra/mockery/v2@latest --name=MerchantHandler --output=./mocks --outpkg=handler --filename=merchant_handler.go --structname=MockMerchantHandler --with-expecter=false
type MerchantHandler interface {
	AddMerchant(c *fiber.Ctx) error
	GetMerchant(c *fiber.Ctx) error
	GetMerchantList(c *fiber.Ctx) error
	UpdateMerchant(c *fiber.Ctx) error
	DeleteMerchant(c *fiber.Ctx) error
}

type merchantHandlerImpl struct {
	merchantService service.MerchantService
}

func NewMerchantHandler(merchantService service.MerchantService) MerchantHandler {
	return &merchantHandlerImpl{
		merchantService: merchantService,
	}
}

// POST /merchant
// Add a new merchant.
// @Summary      Add a new merchant
// @Description  Create a new merchant with the provided details
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.CreateMerchantRequest true "Merchant data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /merchant [post]
func (h *merchantHandlerImpl) AddMerchant(c *fiber.Ctx) error {
	var createMerchantRequest dto.CreateMerchantRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createMerchantRequest); err != nil {
		return err
	}

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil || !isSuperAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	newMerchant, err := h.merchantService.Create(c.UserContext(), createMerchantRequest, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, newMerchant)
}

// GET /merchant/:id
// Get a merchant by ID.
// @Summary      Get a merchant by ID
// @Description  Retrieve a merchant by its ID
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        id path int true "Merchant ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /merchant/{id} [get]
func (h *merchantHandlerImpl) GetMerchant(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get id from param
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid merchant ID")
	}

	merchant, err := h.merchantService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, merchant)
}

// GET /merchant
// Get list of merchants.
// @Summary      Get list of merchants
// @Description  Retrieve a list of all merchants
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /merchant [get]
func (h *merchantHandlerImpl) GetMerchantList(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	merchantList, err := h.merchantService.GetList()
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, merchantList)
}

// PUT /merchant
// Update a merchant.
// @Summary      Update a merchant
// @Description  Update an existing merchant with new details
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.UpdateMerchantRequest true "Updated merchant data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /merchant [put]
func (h *merchantHandlerImpl) UpdateMerchant(c *fiber.Ctx) error {
	var updateMerchantRequest dto.UpdateMerchantRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &updateMerchantRequest); err != nil {
		return err
	}

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil || !isSuperAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.merchantService.Update(c.UserContext(), updateMerchantRequest, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// DELETE /merchant/:id
// Soft-delete a merchant.
// @Summary      Soft-delete a merchant
// @Description  Soft-delete a merchant by ID
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "Merchant ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /merchant/{id} [delete]
func (h *merchantHandlerImpl) DeleteMerchant(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil || !isSuperAdmin {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid merchant ID")
	}

	if err := h.merchantService.Delete(c.UserContext(), id); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
