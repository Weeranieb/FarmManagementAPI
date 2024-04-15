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

type IFarmOnFarmGroupController interface {
	ApplyRoute(router *gin.Engine)
}

type FarmOnFarmGroupControllerImp struct {
	FarmOnFarmGroupService services.IFarmOnFarmGroupService
}

func NewFarmOnFarmGroupController(farmOnFarmGroupService services.IFarmOnFarmGroupService) IFarmOnFarmGroupController {
	return &FarmOnFarmGroupControllerImp{
		FarmOnFarmGroupService: farmOnFarmGroupService,
	}
}

func (c FarmOnFarmGroupControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/farmGroup")
		{
			eg.POST("", c.AddFarmOnFarmGroup)
			eg.DELETE("/:id", c.DeleteFarmOnFarmGroup)
		}
	}
}

func (c FarmOnFarmGroupControllerImp) AddFarmOnFarmGroup(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addFarmOnFarmGroup models.AddFarmOnFarmGroup

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FarmGroup_AddFarmGroup_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addFarmOnFarmGroup); err != nil {
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
		errRes.Error(ctx, "Err_FarmGroup_AddFarmGroup_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newFarmGroup, err := c.FarmOnFarmGroupService.Create(addFarmOnFarmGroup, username)
	if err != nil {
		httputil.NewError(ctx, "Err_FarmGroup_AddFarmGroup_04", err)
		return
	}

	response.Result = true
	response.Data = newFarmGroup

	ctx.JSON(http.StatusOK, response)
}

func (c FarmOnFarmGroupControllerImp) DeleteFarmOnFarmGroup(ctx *gin.Context) {
	var response httputil.ResponseModel
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError(ctx, "Err_FarmGroup_DeleteFarmOnFarmGroup_01", err)
		return
	}

	err = c.FarmOnFarmGroupService.Delete(id)
	if err != nil {
		httputil.NewError(ctx, "Err_FarmGroup_DeleteFarmOnFarmGroup_02", err)
		return
	}

	response.Result = true
	ctx.JSON(http.StatusOK, response)
}
