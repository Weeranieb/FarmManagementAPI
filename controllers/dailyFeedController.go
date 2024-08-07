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

type IDailyFeedController interface {
	ApplyRoute(router *gin.Engine)
}

type dailyFeedControllerImp struct {
	DailyFeedProcessor processors.IDailyFeedProcessor
	DailyFeedService   services.IDailyFeedService
}

func NewDailyFeedController(dailyFeedProcessor processors.IDailyFeedProcessor, dailyFeedService services.IDailyFeedService) IDailyFeedController {
	return &dailyFeedControllerImp{
		DailyFeedProcessor: dailyFeedProcessor,
		DailyFeedService:   dailyFeedService,
	}
}

func (c dailyFeedControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/dailyfeed")
		{
			eg.POST("", c.AddDailyFeed)
			eg.GET(":id", c.GetDailyFeed)
			eg.GET("", c.GetDailyFeedList)
			eg.PUT("", c.UpdateDailyFeed)
			eg.PUT("bulk", c.BulkUpdateDailyFeed)
			eg.GET("/download", c.DownloadExcelForm)
			eg.POST("/upload", c.Upload)
		}
	}
}

// POST /api/v1/dailyfeed
// Add a new daily feed entry.
// @Summary      Add a new daily feed entry
// @Description  Add a new daily feed entry with the provided details
// @Tags         dailyfeed
// @Accept       json
// @Produce      json
// @Param        body body models.AddDailyFeed true "Daily Feed data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/dailyfeed [post]
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

// GET /api/v1/dailyfeed/{id}
// Get daily feed entry by ID.
// @Summary      Get daily feed entry by ID
// @Description  Retrieve details of a specific daily feed entry by its ID
// @Tags         dailyfeed
// @Accept       json
// @Produce      json
// @Param        id path int true "Daily Feed ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/dailyfeed/{id} [get]
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

// GET /api/v1/dailyfeed
// Get daily feed list.
// @Summary      Get daily feed list
// @Description  Retrieve a list of daily feed entries
// @Tags         dailyfeed
// @Accept       json
// @Produce      json
// @Param        feedId query int true "Feed ID"
// @Param        farmId query int true "Farm ID"
// @Param        date query string true "Date (YYYY-MM-DD)"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/dailyfeed [get]
func (c dailyFeedControllerImp) GetDailyFeedList(ctx *gin.Context) {
	var response httputil.ResponseModel
	var dailyFeedList []*models.DailyFeed

	sFeedId := ctx.Query("feedId")
	feedId, err := strconv.Atoi(sFeedId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_GetDailyList_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	sFarmId := ctx.Query("farmId")
	farmId, err := strconv.Atoi(sFarmId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_GetDailyList_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	date := ctx.Query("date") // 2024-01-01

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_DailyFeed_GetDailyList_03", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	dailyFeedList, err = c.DailyFeedService.GetDailyFeedList(feedId, farmId, date)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_GetDailyList_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = dailyFeedList

	ctx.JSON(http.StatusOK, response)
}

// PUT /api/v1/dailyfeed
// Update daily feed entry.
// @Summary      Update daily feed entry
// @Description  Update details of a daily feed entry
// @Tags         dailyfeed
// @Accept       json
// @Produce      json
// @Param        body body models.DailyFeed true "Daily Feed data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/dailyfeed [put]
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

func (c dailyFeedControllerImp) DownloadExcelForm(ctx *gin.Context) {
	var response httputil.ResponseModel

	formType := ctx.Query("type")
	if formType == "" {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_DownloadExcelForm_02", "Missing form type")
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}
	sFeedId := ctx.Query("feedId")
	feedId, err := strconv.Atoi(sFeedId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_DownloadExcelForm_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	sFarmId := ctx.Query("farmId")
	farmId, err := strconv.Atoi(sFarmId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_DownloadExcelForm_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	date := ctx.Query("date") // 2024-01-01

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_DailyFeed_DownloadExcelForm_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_DownloadExcelForm_05", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	excelBytes, err := c.DailyFeedProcessor.DownloadExcelForm(clientId, formType, feedId, farmId, date)
	if err != nil {
		httputil.NewError(ctx, "Err_DailyFeed_DownloadExcelForm_02", err)
		return
	}

	// Set response headers to trigger download
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=dailyfeed.xlsx")
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Length", fmt.Sprint(len(excelBytes)))

	// Write the byte data to the response
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}

func (c dailyFeedControllerImp) Upload(ctx *gin.Context) {
	var response httputil.ResponseModel

	excelFile, err := ctx.FormFile("file")
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_Upload_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_Upload_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get client id
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_Upload_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.DailyFeedProcessor.UploadExcelForm(excelFile, username, clientId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_Upload_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	ctx.JSON(http.StatusOK, response)
}

func (c dailyFeedControllerImp) BulkUpdateDailyFeed(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateDailyFeedList []*models.DailyFeed

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_DailyFeed_BulkUpdateDailyFeed_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateDailyFeedList); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_BulkUpdateDailyFeed_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_DailyFeed_BulkUpdateDailyFeed_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.DailyFeedProcessor.BulkCreateAndUpdate(updateDailyFeedList, username)
	if err != nil {
		httputil.NewError(ctx, "Err_DailyFeed_BulkUpdateDailyFeed_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
