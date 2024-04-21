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

type IDailyFeedController interface {
	ApplyRoute(router *gin.Engine)
}

type dailyFeedControllerImp struct {
	DailyFeedService services.IDailyFeedService
}

func NewDailyFeedController(dailyFeedService services.IDailyFeedService) IDailyFeedController {
	return &dailyFeedControllerImp{
		DailyFeedService: dailyFeedService,
	}
}

func (c dailyFeedControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/dailyfeed")
		{
			eg.POST("", c.AddDailyFeed)
			eg.GET(":id", c.GetDailyFeed)
			eg.PUT("", c.UpdateDailyFeed)
		}
	}
}

func (c dailyFeedControllerImp) AddDailyFeed(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addDailyFeed models.AddDailyFeed

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_DailyFeed_AddDailyFeed_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addDailyFeed); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_AddDailyFeed_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_AddDailyFeed_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newPond, err := c.DailyFeedService.Create(addDailyFeed, username)
	if err != nil {
		httputil.NewError(ctx, "Err_DailyFeed_AddDailyFeed_04", err)
		return
	}

	response.Result = true
	response.Data = newPond

	ctx.JSON(http.StatusOK, response)
}

func (c dailyFeedControllerImp) GetDailyFeed(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_GetDailyFeed_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_DailyFeed_GetDailyFeed_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	pond, err := c.DailyFeedService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_GetDailyFeed_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = pond

	ctx.JSON(http.StatusOK, response)
}

func (c dailyFeedControllerImp) UpdateDailyFeed(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateDailyFeed *models.DailyFeed

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_DailyFeed_UpdateDailyFeed_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateDailyFeed); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_UpdateDailyFeed_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_UpdateDailyFeed_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.DailyFeedService.Update(updateDailyFeed, username)
	if err != nil {
		httputil.NewError(ctx, "Err_DailyFeed_UpdateDailyFeed_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
