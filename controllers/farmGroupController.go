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

type farmGroupControllerImp struct {
	FarmGroupService services.IFarmGroupService
}

func NewFarmGroupController(farmGroupService services.IFarmGroupService) IFarmController {
	return &farmGroupControllerImp{
		FarmGroupService: farmGroupService,
	}
}

func (c farmGroupControllerImp) ApplyRoute(router *gin.Engine) {
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

// POST /api/v1/farmGroup
// Add a new farm group entry.
// @Summary      Add a new farm group entry
// @Description  Add a new farm group entry with the provided details
// @Tags         farmGroup
// @Accept       json
// @Produce      json
// @Param        body body models.AddFarmGroup true "Farm group data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farmGroup [post]
func (c farmGroupControllerImp) AddFarmGroup(ctx *gin.Context) {
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

// GET /api/v1/farmGroup/{id}
// Get farm group by ID.
// @Summary      Get farm group by ID
// @Description  Retrieve details of a specific farm group by its ID
// @Tags         farmGroup
// @Accept       json
// @Produce      json
// @Param        id path int true "Farm group ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farmGroup/{id} [get]
func (c farmGroupControllerImp) GetFarmGroup(ctx *gin.Context) {
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

// PUT /api/v1/farmGroup
// Update farm group entry.
// @Summary      Update farm group entry
// @Description  Update details of a farm group entry
// @Tags         farmGroup
// @Accept       json
// @Produce      json
// @Param        body body models.FarmGroup true "Farm group data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farmGroup [put]
func (c farmGroupControllerImp) UpdateFarmGroup(ctx *gin.Context) {
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

// GET /api/v1/farmGroup/{id}/farmList
// Get list of farms associated with a specific farm group.
// @Summary      Get list of farms associated with a specific farm group
// @Description  Retrieve a list of farms associated with a specific farm group by its ID
// @Tags         farmGroup
// @Accept       json
// @Produce      json
// @Param        id path int true "Farm group ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farmGroup/{id}/farmList [get]
func (c farmGroupControllerImp) GetFarmList(ctx *gin.Context) {
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
