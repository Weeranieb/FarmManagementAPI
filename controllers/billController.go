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

type IBillController interface {
	ApplyRoute(router *gin.Engine)
}

type billControllerImp struct {
	BillService services.IBillService
}

func NewBillController(billService services.IBillService) IBillController {
	return &billControllerImp{
		BillService: billService,
	}
}

func (c billControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/bill")
		{
			eg.POST("", c.AddBill)
			eg.GET(":id", c.GetBill)
			eg.PUT("", c.UpdateBill)
		}
	}
}

func (c billControllerImp) AddBill(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addBill models.AddBill

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Bill_AddBill_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addBill); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_AddBill_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_AddBill_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newBill, err := c.BillService.Create(addBill, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Bill_AddBill_04", err)
		return
	}

	response.Result = true
	response.Data = newBill

	ctx.JSON(http.StatusOK, response)
}

func (c billControllerImp) GetBill(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_GetBill_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Bill_GetBill_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	bill, err := c.BillService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_GetBill_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = bill

	ctx.JSON(http.StatusOK, response)
}

func (c billControllerImp) UpdateBill(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateBill *models.Bill

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Bill_UpdateBill_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateBill); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_UpdateBill_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_UpdateBill_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.BillService.Update(updateBill, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Bill_UpdateBill_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
