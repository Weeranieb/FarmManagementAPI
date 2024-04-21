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

type IFeedCollectionController interface {
	ApplyRoute(router *gin.Engine)
}

type feedCollectionControllerImp struct {
	FeedCollectionService services.IFeedCollectionService
}

func NewFeedCollectionController(feedCollectionService services.IFeedCollectionService) IFeedCollectionController {
	return &feedCollectionControllerImp{
		FeedCollectionService: feedCollectionService,
	}
}

func (c feedCollectionControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/feedcollection")
		{
			eg.POST("", c.AddFeedCollection)
			eg.GET(":id", c.GetFeedCollection)
			eg.PUT("", c.UpdateFeedCollection)
		}
	}
}

func (c feedCollectionControllerImp) AddFeedCollection(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addFeedCollection models.AddFeedCollection

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FeedCollection_AddFeedCollection_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addFeedCollection); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_AddFeedCollection_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_AddFeedCollection_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newPond, err := c.FeedCollectionService.Create(addFeedCollection, username)
	if err != nil {
		httputil.NewError(ctx, "Err_FeedCollection_AddFeedCollection_04", err)
		return
	}

	response.Result = true
	response.Data = newPond

	ctx.JSON(http.StatusOK, response)
}

func (c feedCollectionControllerImp) GetFeedCollection(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_GetFeedCollection_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FeedCollection_GetFeedCollection_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	pond, err := c.FeedCollectionService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_GetFeedCollection_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = pond

	ctx.JSON(http.StatusOK, response)
}

func (c feedCollectionControllerImp) UpdateFeedCollection(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateFeedCollection *models.FeedCollection

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FeedCollection_UpdateFeedCollection_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateFeedCollection); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_UpdateFeedCollection_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_UpdateFeedCollection_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.FeedCollectionService.Update(updateFeedCollection, username)
	if err != nil {
		httputil.NewError(ctx, "Err_FeedCollection_UpdateFeedCollection_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
