package controllers

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/processors"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/httputil"
	"boonmafarm/api/utils/jwtutil"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IActivityController interface {
	ApplyRoute(router *gin.Engine)
}

type activityControllerImp struct {
	ActivityService   services.IActivityService
	ActivityProcessor processors.IActivityProcessor
}

func NewActivityController(activityService services.IActivityService, activityProcessor processors.IActivityProcessor) IActivityController {
	return &activityControllerImp{
		ActivityService:   activityService,
		ActivityProcessor: activityProcessor,
	}
}

func (c activityControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/activity")
		{
			// eg.POST("", c.AddActivity)
			eg.POST("fill", c.AddFillActivity)
			eg.POST("move", c.AddMoveActivity)
			eg.GET(":id", c.GetActivity)
			eg.PUT("", c.UpdateActivity)
			eg.GET("", c.ListActivity)
		}
	}
}

func (c activityControllerImp) AddFillActivity(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addFillActivity models.CreateFillActivityRequest

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Activity_AddFillActivity_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addFillActivity); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_AddAFillctivity_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_AddFillActivity_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newActivity, activePond, err := c.ActivityProcessor.CreateFill(addFillActivity, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Activity_AddFillActivity_04", err)
		return
	}

	type ret struct {
		Activity   *models.Activity   `json:"activity"`
		ActivePond *models.ActivePond `json:"activePond"`
	}

	response.Result = true
	response.Data = ret{
		Activity:   newActivity,
		ActivePond: activePond,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c activityControllerImp) AddMoveActivity(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addMoveActivity models.CreateMoveActivityRequest

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Activity_AddFillActivity_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addMoveActivity); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_AddFillctivity_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_AddFillActivity_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newActivity, fromActivePond, toActivePond, err := c.ActivityProcessor.CreateMove(addMoveActivity, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Activity_AddFillActivity_04", err)
		return
	}

	type ret struct {
		Activity       *models.Activity   `json:"activity"`
		FromActivePond *models.ActivePond `json:"fromActivePond"`
		ToActivePond   *models.ActivePond `json:"toActivePond"`
	}

	response.Result = true
	response.Data = ret{
		Activity:       newActivity,
		FromActivePond: fromActivePond,
		ToActivePond:   toActivePond,
	}

	ctx.JSON(http.StatusOK, response)
}

// func (c activityControllerImp) AddSellActivity(ctx *gin.Context) {
// 	var response httputil.ResponseModel
// 	var addFillActivity models.CreateFillActivityRequest

// 	defer func() {
// 		if r := recover(); r != nil {
// 			errRes := httputil.ErrorResponseModel{}
// 			errRes.Error(ctx, "Err_Activity_AddSellActivity_01", fmt.Sprint(r))
// 			response.Error = errRes
// 			ctx.JSON(http.StatusOK, response)
// 			return
// 		}
// 	}()

// 	if err := ctx.ShouldBindJSON(&addFillActivity); err != nil {
// 		errRes := httputil.ErrorResponseModel{}
// 		errRes.Error(ctx, "Err_Activity_AddSellctivity_02", err.Error())
// 		response.Error = errRes
// 		ctx.JSON(http.StatusOK, response)
// 		return
// 	}

// 	// get username
// 	username, err := jwtutil.GetUsername(ctx)
// 	if err != nil {
// 		errRes := httputil.ErrorResponseModel{}
// 		errRes.Error(ctx, "Err_Activity_AddSellActivity_03", err.Error())
// 		response.Error = errRes
// 		ctx.JSON(http.StatusOK, response)
// 		return
// 	}

// 	newActivity, err := c.ActivityService.CreateSell(addFillActivity, username)
// 	if err != nil {
// 		httputil.NewError(ctx, "Err_Activity_AddSellActivity_04", err)
// 		return
// 	}

// 	response.Result = true
// 	response.Data = newActivity

// 	ctx.JSON(http.StatusOK, response)
// }

// GetActivity retrieves an activity based on the provided ID.
// It handles the HTTP GET request and returns the activity as a JSON response.
// @Summary      Get an activity by ID
// @Description  Get an activity by ID
// @Tags         activity
// @Accept       json
// @Produce      json
// @Param        id path int true "Activity ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      404  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/activity/{id} [get]
func (c activityControllerImp) GetActivity(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_GetActivity_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Activity_GetActivity_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	activity, err := c.ActivityService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_GetActivity_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = activity

	ctx.JSON(http.StatusOK, response)
}

// UpdateActivity updates an activity.
// It handles the HTTP PUT request and expects the request body to contain JSON data representing the updated activity.
// It returns a JSON response indicating success or failure of the update operation, along with any relevant data.
// @Summary      Update an activity
// @Description  Update an activity
// @Tags         activity
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body models.ActivityWithSellDetail true "Updated Activity data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/activity [put]
func (c activityControllerImp) UpdateActivity(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateActivity *models.ActivityWithSellDetail

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Activity_UpdateActivity_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateActivity); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_UpdateActivity_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_UpdateActivity_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	var result []*models.SellDetail
	result, err = c.ActivityService.Update(updateActivity, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Activity_UpdateActivity_04", err)
		return
	}

	var ret = struct {
		SellDetail []*models.SellDetail `json:"sellDetails"`
	}{
		SellDetail: result,
	}
	response.Result = true
	response.Data = ret

	ctx.JSON(http.StatusOK, response)
}

func (c activityControllerImp) ListActivity(ctx *gin.Context) {
	var response httputil.ResponseModel

	sPage := ctx.Query("page")
	sPageSize := ctx.Query("pageSize")
	orderBy := ctx.Query("orderBy")
	keyword := ctx.Query("keyword")
	mode := ctx.Query("mode")
	sfarmId := ctx.Query("farmId")

	page, err := strconv.Atoi(sPage)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_ListActivity_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	pageSize, err := strconv.Atoi(sPageSize)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_ListActivity_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	farmId, _ := strconv.Atoi(sfarmId)

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Activity_ListActivity_04", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_ListActivity_05", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Data, err = c.ActivityService.TakePage(clientId, page, pageSize, orderBy, keyword, mode, farmId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Activity_ListActivity_06", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
