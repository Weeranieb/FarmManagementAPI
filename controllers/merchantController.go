package controllers

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/httputil"
	"boonmafarm/api/utils/jwtutil"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IMerchantController interface {
	ApplyRoute(router *gin.Engine)
}

type merchantControllerImp struct {
	MerchantService services.IMerchantService
}

func NewMerchantController(merchantService services.IMerchantService) IMerchantController {
	return &merchantControllerImp{
		MerchantService: merchantService,
	}
}

func (c merchantControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/merchant")
		{
			eg.POST("", c.AddMerchant)
			eg.GET(":id", c.GetMerchant)
			eg.PUT("", c.UpdateMerchant)
		}
	}
}

// POST /api/v1/merchant
// Add a new merchant.
// @Summary      Add a new merchant
// @Description  Create a new merchant with the provided details
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        body body models.AddMerchant true "Merchant data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/merchant [post]
func (c merchantControllerImp) AddMerchant(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addMerchant models.AddMerchant

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Merchant_AddMerchant_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addMerchant); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Merchant_AddMerchant_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Merchant_AddMerchant_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newMerchant, err := c.MerchantService.Create(addMerchant, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Merchant_AddMerchant_04", err)
		return
	}

	response.Result = true
	response.Data = newMerchant

	ctx.JSON(http.StatusOK, response)
}

// GET /api/v1/merchant/{id}
// Get a merchant by ID.
// @Summary      Get a merchant by ID
// @Description  Retrieve a merchant by its ID
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        id path int true "Merchant ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      404  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/merchant/{id} [get]
func (c merchantControllerImp) GetMerchant(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Merchant_GetMerchant_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Merchant_GetMerchant_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	merchant, err := c.MerchantService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Merchant_GetMerchant_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = merchant

	ctx.JSON(http.StatusOK, response)
}

// PUT /api/v1/merchant
// Update a merchant.
// @Summary      Update a merchant
// @Description  Update an existing merchant with new details
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        body body models.Merchant true "Updated merchant data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/merchant [put]
func (c merchantControllerImp) UpdateMerchant(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateMerchant *models.Merchant

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Merchant_UpdateMerchant_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateMerchant); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Merchant_UpdateMerchant_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Merchant_UpdateMerchant_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.MerchantService.Update(updateMerchant, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Merchant_UpdateMerchant_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
