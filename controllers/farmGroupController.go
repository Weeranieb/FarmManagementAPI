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

type IFarmGroupController interface {
	ApplyRoute(router *gin.Engine)
}

type FarmGroupControllerImp struct {
	FarmGroupService services.IFarmGroupService
}

func NewFarmGroupController(farmGroupService services.IFarmGroupService) IFarmController {
	return &FarmGroupControllerImp{
		FarmGroupService: farmGroupService,
	}
}

func (c FarmGroupControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/farmGroup")
		{
			eg.POST("", c.AddFarmGroup)
			eg.GET(":id", c.GetFarmGroup)
			eg.PUT("", c.UpdateFarmGroup)
			eg.GET(":id/farmList", c.GetFarmList)
		}
	}
}

func (c FarmGroupControllerImp) AddFarmGroup(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addFarmGroupp models.AddFarmGroup

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FarmGroup_AddFarmGroup_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addFarmGroupp); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_AddFarmGroup_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_AddFarmGroup_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get clientId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_AddFarmGroup_05", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newFarmGroup, err := c.FarmGroupService.Create(addFarmGroupp, username, clientId)
	if err != nil {
		httputil.NewError(ctx, "Err_FarmGroup_AddFarmGroup_03", err)
		return
	}

	response.Result = true
	response.Data = newFarmGroup

	ctx.JSON(http.StatusOK, response)
}

func (c FarmGroupControllerImp) GetFarmGroup(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from param
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_GetFarmGroup_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FarmGroup_GetFarmGroup_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	// get clientId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_GetFarmGroup_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	farmGroup, err := c.FarmGroupService.Get(id, clientId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_GetFarmGroup_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = farmGroup

	ctx.JSON(http.StatusOK, response)
}

func (c FarmGroupControllerImp) UpdateFarmGroup(ctx *gin.Context) {
	var response httputil.ResponseModel
	var UpdateFarm *models.FarmGroup

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FarmGroup_UpdateFarmGroup_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&UpdateFarm); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_UpdateFarmGroup_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_UpdateFarmGroup_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.FarmGroupService.Update(UpdateFarm, username)
	if err != nil {
		httputil.NewError(ctx, "Err_FarmGroup_UpdateFarmGroup_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}

func (c FarmGroupControllerImp) GetFarmList(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from param
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_GetFarmList_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FarmGroup_GetFarmList_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	farmList, err := c.FarmGroupService.GetFarmList(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FarmGroup_GetFarmList_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = farmList

	ctx.JSON(http.StatusOK, response)
}
