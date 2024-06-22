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

type farmOnFarmGroupControllerImp struct {
	FarmOnFarmGroupService services.IFarmOnFarmGroupService
}

func NewFarmOnFarmGroupController(farmOnFarmGroupService services.IFarmOnFarmGroupService) IFarmOnFarmGroupController {
	return &farmOnFarmGroupControllerImp{
		FarmOnFarmGroupService: farmOnFarmGroupService,
	}
}

func (c farmOnFarmGroupControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/farmOnFarmGroup")
		{
			eg.POST("", c.AddFarmOnFarmGroup)
			eg.DELETE(":id", c.DeleteFarmOnFarmGroup)
		}
	}
}

// POST /api/v1/farmOnFarmGroup
// Add a farm to a farm group.
// @Summary      Add a farm to a farm group
// @Description  Add a farm to a farm group with the provided details
// @Tags         farmOnFarmGroup
// @Accept       json
// @Produce      json
// @Param        body body models.AddFarmOnFarmGroup true "Farm on farm group data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farmOnFarmGroup [post]
func (c farmOnFarmGroupControllerImp) AddFarmOnFarmGroup(ctx *gin.Context) {
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

// DELETE /api/v1/farmOnFarmGroup/{id}
// Remove a farm from a farm group.
// @Summary      Remove a farm from a farm group
// @Description  Remove a farm from a farm group by its ID
// @Tags         farmOnFarmGroup
// @Accept       json
// @Produce      json
// @Param        id path int true "Farm on farm group ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/farmOnFarmGroup/{id} [delete]
func (c farmOnFarmGroupControllerImp) DeleteFarmOnFarmGroup(ctx *gin.Context) {
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
