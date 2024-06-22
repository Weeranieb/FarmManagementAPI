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

type farmControllerImp struct {
	FarmService services.IFarmService
}

func NewFarmController(farmService services.IFarmService) IFarmController {
	return &farmControllerImp{
		FarmService: farmService,
	}
}

func (c farmControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/farm")
		{
			eg.POST("", c.AddFarm)
			eg.GET("", c.GetFarmList)
			eg.GET(":id", c.GetFarm)
			eg.PUT("", c.UpdateFarm)
		}
	}
}

// POST /api/v1/farm
// Add a new farm entry.
// @Summary      Add a new farm entry
// @Description  Add a new farm entry with the provided details
// @Tags         farm
// @Accept       json
// @Produce      json
// @Param        body body models.AddFarm true "Farm data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farm [post]
func (c farmControllerImp) AddFarm(ctx *gin.Context) {
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

// GET /api/v1/farm
// Get list of farms associated with the current client.
// @Summary      Get list of farms
// @Description  Retrieve a list of farms associated with the current client
// @Tags         farm
// @Accept       json
// @Produce      json
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farm [get]

func (c farmControllerImp) GetFarm(ctx *gin.Context) {
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

// GET /api/v1/farm/{id}
// Get farm by ID.
// @Summary      Get farm by ID
// @Description  Retrieve details of a specific farm by its ID
// @Tags         farm
// @Accept       json
// @Produce      json
// @Param        id path int true "Farm ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farm/{id} [get]
func (c farmControllerImp) UpdateFarm(ctx *gin.Context) {
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

// PUT /api/v1/farm
// Update farm entry.
// @Summary      Update farm entry
// @Description  Update details of a farm entry
// @Tags         farm
// @Accept       json
// @Produce      json
// @Param        body body models.Farm true "Farm data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farm [put]
func (c farmControllerImp) GetFarmList(ctx *gin.Context) {
	var response httputil.ResponseModel
	var farmList []*models.Farm

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Farm_GetFarmList_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	// get clientId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_GetFarmList_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	farmList, err = c.FarmService.GetList(clientId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Farm_GetFarmList_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = farmList

	ctx.JSON(http.StatusOK, response)
}
