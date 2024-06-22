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
			eg.GET("", c.GetActivePondList)
		}
	}
}

// AddActivePond godoc
// @Summary      Add a new active pond
// @Description  Add a new active pond
// @Tags         activepond
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body models.AddActivePond true "New Active Pond data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/activepond [post]
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

// GetActivePond godoc
// @Summary      Get an active pond by ID
// @Description  Get an active pond by ID
// @Tags         activepond
// @Accept       json
// @Produce      json
// @Param        id path int true "Active Pond ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      404  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/activepond/{id} [get]
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

	activePond, err := c.ActivePondService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_GetActivePond_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = activePond

	ctx.JSON(http.StatusOK, response)
}

// UpdateActivePond godoc
// @Summary      Update an active pond
// @Description  Update an active pond
// @Tags         activepond
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body models.ActivePond true "Updated Active Pond data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/activepond [put]
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

func (c activePondControllerImp) GetActivePondList(ctx *gin.Context) {
	var response httputil.ResponseModel
	var farmId int
	// get param from query
	farmIds := ctx.Query("farmId")
	// convert string to int
	farmId, err := strconv.Atoi(farmIds)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_GetActivePondList_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_ActivePond_GetActivePondList_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	ponds, err := c.ActivePondService.GetList(farmId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_ActivePond_GetActivePondList_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = ponds

	ctx.JSON(http.StatusOK, response)
}
