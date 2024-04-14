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

type IFarmController interface {
	ApplyRoute(router *gin.Engine)
}

type FarmControllerImp struct {
	FarmService services.IFarmService
}

func NewFarmController(farmService services.IFarmService) IFarmController {
	return &FarmControllerImp{
		FarmService: farmService,
	}
}

func (c FarmControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/farm")
		{
			eg.POST("", c.AddFarm)
			eg.GET("/:id", c.GetFarm)
			eg.PUT("", c.UpdateFarm)
		}
	}
}

func (c FarmControllerImp) AddFarm(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addFarm models.AddFarm

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Farm_AddFarm_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addFarm); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_AddFarm_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_AddFarm_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get clientId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_AddFarm_05", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newFarm, err := c.FarmService.Create(addFarm, username, clientId)
	if err != nil {
		httputil.NewError(ctx, "Err_Farm_AddFarm_03", err)
		return
	}

	response.Result = true
	response.Data = newFarm

	ctx.JSON(http.StatusOK, response)
}

func (c FarmControllerImp) GetFarm(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from param
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_GetFarm_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Farm_GetFarm_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	// get clientId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_GetFarm_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	farm, err := c.FarmService.Get(id, clientId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_GetFarm_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = farm

	ctx.JSON(http.StatusOK, response)
}

func (c FarmControllerImp) UpdateFarm(ctx *gin.Context) {
	var response httputil.ResponseModel
	var UpdateFarm *models.Farm

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Farm_UpdateFarm_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&UpdateFarm); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_UpdateFarm_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_UpdateFarm_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.FarmService.Update(UpdateFarm, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Farm_UpdateFarm_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
