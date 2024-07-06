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
			eg.GET("", c.ListFeedCollection)
		}
	}
}

// POST /api/v1/feedcollection
// Add a feed collection.
// @Summary      Add a feed collection
// @Description  Add a new feed collection with the provided details
// @Tags         feedcollection
// @Accept       json
// @Produce      json
// @Param        body body models.AddFeedCollection true "Feed collection data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/feedcollection [post]
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

	// get client id
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_AddFeedCollection_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newPond, err := c.FeedCollectionService.Create(addFeedCollection, username, clientId)
	if err != nil {
		httputil.NewError(ctx, "Err_FeedCollection_AddFeedCollection_04", err)
		return
	}

	response.Result = true
	response.Data = newPond

	ctx.JSON(http.StatusOK, response)
}

// GET /api/v1/feedcollection/{id}
// Get a feed collection by ID.
// @Summary      Get a feed collection by ID
// @Description  Retrieve a feed collection by its ID
// @Tags         feedcollection
// @Accept       json
// @Produce      json
// @Param        id path int true "Feed collection ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      404  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/feedcollection/{id} [get]
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

// PUT /api/v1/feedcollection
// Update a feed collection.
// @Summary      Update a feed collection
// @Description  Update an existing feed collection with new details
// @Tags         feedcollection
// @Accept       json
// @Produce      json
// @Param        body body models.FeedCollection true "Updated feed collection data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/feedcollection [put]
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

func (c feedCollectionControllerImp) ListFeedCollection(ctx *gin.Context) {
	var response httputil.ResponseModel

	sPage := ctx.Query("page")
	sPageSize := ctx.Query("pageSize")
	orderBy := ctx.Query("orderBy")
	keyword := ctx.Query("keyword")

	page, err := strconv.Atoi(sPage)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_ListFeedCollection_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	pageSize, err := strconv.Atoi(sPageSize)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_ListFeedCollection_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_FeedCollection_ListFeedCollection_03", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_ListFeedCollection_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Data, err = c.FeedCollectionService.TakePage(clientId, page, pageSize, orderBy, keyword)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_FeedCollection_ListFeedCollection_05", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
