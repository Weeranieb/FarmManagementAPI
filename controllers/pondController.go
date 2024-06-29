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

type IPondController interface {
	ApplyRoute(router *gin.Engine)
}

type pondControllerImp struct {
	PondService services.IPondService
}

func NewPondController(pondService services.IPondService) IPondController {
	return &pondControllerImp{
		PondService: pondService,
	}
}

func (c pondControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/pond")
		{
			eg.POST("", c.AddPond)
			eg.POST("/batch", c.AddPonds)
			eg.GET(":id", c.GetPond)
			eg.PUT("", c.UpdatePond)
			eg.GET("", c.GetPondList)
		}
	}
}

// POST /api/v1/pond
// Add a new pond.
// @Summary Add a new pond
// @Description Create a new pond with the provided details
// @Tags pond
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param body body models.AddPond true "Pond data"
// @Success 200 {object} httputil.ResponseModel
// @Failure 400 {object} httputil.ErrorResponseModel
// @Failure 500 {object} httputil.ErrorResponseModel
// @Router /api/v1/pond [post]
func (c pondControllerImp) AddPond(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addPond models.AddPond

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Pond_AddPond_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addPond); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_AddPond_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_AddPond_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newPond, err := c.PondService.Create(addPond, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Pond_AddPond_04", err)
		return
	}

	response.Result = true
	response.Data = newPond

	ctx.JSON(http.StatusOK, response)
}

func (c pondControllerImp) AddPonds(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addPonds []models.AddPond

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Pond_AddPonds_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addPonds); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_AddPonds_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_AddPonds_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newPonds, err := c.PondService.CreateBatch(addPonds, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Pond_AddPonds_04", err)
		return
	}

	response.Result = true
	response.Data = newPonds

	ctx.JSON(http.StatusOK, response)
}

// GET /api/v1/pond/{id}
// Get a pond by ID.
// @Summary      Get a pond by ID
// @Description  Retrieve a pond by its ID
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        id path int true "Pond ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      404  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/pond/{id} [get]
func (c pondControllerImp) GetPond(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_GetPond_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Pond_GetPond_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	pond, err := c.PondService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_GetPond_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = pond

	ctx.JSON(http.StatusOK, response)
}

// PUT /api/v1/pond
// Update a pond.
// @Summary      Update a pond
// @Description  Update an existing pond with new details
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        body body models.Pond true "Updated pond data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/pond [put]
func (c pondControllerImp) UpdatePond(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updatePond *models.Pond

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Pond_UpdatePond_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updatePond); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_UpdatePond_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_UpdatePond_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.PondService.Update(updatePond, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Pond_UpdatePond_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}

// GET /api/v1/pond
// Get a list of ponds by farm ID.
// @Summary      Get a list of ponds by farm ID
// @Description  Retrieve a list of ponds belonging to a specific farm
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        farmId query int true "Farm ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/pond [get]
func (c pondControllerImp) GetPondList(ctx *gin.Context) {
	var response httputil.ResponseModel
	var farmId int
	// get param from query
	farmIds := ctx.Query("farmId")
	// convert string to int
	farmId, err := strconv.Atoi(farmIds)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_GetPondList_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Pond_GetPondList_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	ponds, err := c.PondService.GetList(farmId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Pond_GetPondList_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = ponds

	ctx.JSON(http.StatusOK, response)
}
