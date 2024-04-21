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

type IFeedPriceHistoryController interface {
	ApplyRoute(router *gin.Engine)
}

type feedPriceHistoryControllerImp struct {
	FeedPriceHistoryService services.IFeedPriceHistoryService
}

func NewFeedPriceHistoryController(feedPriceHistoryService services.IFeedPriceHistoryService) IFeedPriceHistoryController {
	return &feedPriceHistoryControllerImp{
		FeedPriceHistoryService: feedPriceHistoryService,
	}
}

func (c feedPriceHistoryControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/feedpricehistory")
		{
			eg.POST("", c.AddFeedPriceHistory)
			eg.GET(":id", c.GetFeedPriceHistory)
			eg.PUT("", c.UpdateFeedPriceHistory)
		}
	}
}

func (c feedPriceHistoryControllerImp) AddFeedPriceHistory(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addPond models.AddFeedPriceHistory

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FeedPriceHistory_AddFeedPriceHistory_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addPond); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedPriceHistory_AddFeedPriceHistory_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedPriceHistory_AddFeedPriceHistory_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newPond, err := c.FeedPriceHistoryService.Create(addPond, username)
	if err != nil {
		httputil.NewError(ctx, "Err_FeedPriceHistory_AddFeedPriceHistory_04", err)
		return
	}

	response.Result = true
	response.Data = newPond

	ctx.JSON(http.StatusOK, response)
}

func (c feedPriceHistoryControllerImp) GetFeedPriceHistory(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedPriceHistory_GetFeedPriceHistory_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FeedPriceHistory_GetFeedPriceHistory_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	pond, err := c.FeedPriceHistoryService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedPriceHistory_GetFeedPriceHistory_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = pond

	ctx.JSON(http.StatusOK, response)
}

func (c feedPriceHistoryControllerImp) UpdateFeedPriceHistory(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateFeedPriceHistory *models.FeedPriceHistory

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FeedPriceHistory_UpdateFeedPriceHistory_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateFeedPriceHistory); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedPriceHistory_UpdateFeedPriceHistory_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedPriceHistory_UpdateFeedPriceHistory_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.FeedPriceHistoryService.Update(updateFeedPriceHistory, username)
	if err != nil {
		httputil.NewError(ctx, "Err_FeedPriceHistory_UpdateFeedPriceHistory_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
