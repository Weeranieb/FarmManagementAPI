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

type IActivePondController interface {
	ApplyRoute(router *gin.Engine)
}

type activePondControllerImp struct {
	ActivePondService services.IActivePondService
}

func NewActivePondController(activePondService services.IActivePondService) IActivePondController {
	return &activePondControllerImp{
		ActivePondService: activePondService,
	}
}

func (c activePondControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/activepond")
		{
			eg.POST("", c.AddActivePond)
			eg.GET(":id", c.GetActivePond)
			eg.PUT("", c.UpdateActivePond)
		}
	}
}

func (c activePondControllerImp) AddActivePond(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addActivePond models.AddActivePond

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_ActivePond_AddActivePond_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addActivePond); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_AddActivePond_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_AddActivePond_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newPond, err := c.ActivePondService.Create(addActivePond, username)
	if err != nil {
		httputil.NewError(ctx, "Err_ActivePond_AddActivePond_04", err)
		return
	}

	response.Result = true
	response.Data = newPond

	ctx.JSON(http.StatusOK, response)
}

func (c activePondControllerImp) GetActivePond(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_GetActivePond_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_ActivePond_GetActivePond_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	users, err := c.ActivePondService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_GetActivePond_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = users

	ctx.JSON(http.StatusOK, response)
}

func (c activePondControllerImp) UpdateActivePond(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateActivePond *models.ActivePond

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_ActivePond_UpdateActivePond_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateActivePond); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_UpdateActivePond_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_UpdateActivePond_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.ActivePondService.Update(updateActivePond, username)
	if err != nil {
		httputil.NewError(ctx, "Err_ActivePond_UpdateActivePond_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
