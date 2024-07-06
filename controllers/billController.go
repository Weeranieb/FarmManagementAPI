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

type IBillController interface {
	ApplyRoute(router *gin.Engine)
}

type billControllerImp struct {
	BillService services.IBillService
}

func NewBillController(billService services.IBillService) IBillController {
	return &billControllerImp{
		BillService: billService,
	}
}

func (c billControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/bill")
		{
			eg.POST("", c.AddBill)
			eg.GET(":id", c.GetBill)
			eg.PUT("", c.UpdateBill)
			eg.GET("", c.ListBill)
		}
	}
}

// POST /api/v1/bill
// Add a new bill.
// @Summary      Add a new bill
// @Description  Add a new bill with the provided details
// @Tags         bill
// @Accept       json
// @Produce      json
// @Param        body body models.AddBill true "Bill data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/bill [post]
func (c billControllerImp) AddBill(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addBill models.AddBill

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Bill_AddBill_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addBill); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_AddBill_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_AddBill_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newBill, err := c.BillService.Create(addBill, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Bill_AddBill_04", err)
		return
	}

	response.Result = true
	response.Data = newBill

	ctx.JSON(http.StatusOK, response)
}

// GET /api/v1/bill/{id}
// Get bill by ID.
// @Summary      Get bill by ID
// @Description  Get a bill by its ID
// @Tags         bill
// @Accept       json
// @Produce      json
// @Param        id path int true "Bill ID"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      404  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/bill/{id} [get]
func (c billControllerImp) GetBill(ctx *gin.Context) {
	var response httputil.ResponseModel
	// get id from params
	ids := ctx.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_GetBill_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Bill_GetBill_02", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	bill, err := c.BillService.Get(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_GetBill_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = bill

	ctx.JSON(http.StatusOK, response)
}

// PUT /api/v1/bill
// Update an existing bill.
// @Summary      Update an existing bill
// @Description  Update an existing bill with the provided details
// @Tags         bill
// @Accept       json
// @Produce      json
// @Param        body body models.Bill true "Bill data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/bill [put]
func (c billControllerImp) UpdateBill(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateBill *models.Bill

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Bill_UpdateBill_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateBill); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_UpdateBill_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_UpdateBill_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.BillService.Update(updateBill, username)
	if err != nil {
		httputil.NewError(ctx, "Err_Bill_UpdateBill_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}

func (c billControllerImp) ListBill(ctx *gin.Context) {
	var response httputil.ResponseModel

	sPage := ctx.Query("page")
	sPageSize := ctx.Query("pageSize")
	orderBy := ctx.Query("orderBy")
	keyword := ctx.Query("keyword")
	billType := ctx.Query("type")
	sfarmGroupId := ctx.Query("farmGroupId")

	page, err := strconv.Atoi(sPage)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_ListBill_01", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	pageSize, err := strconv.Atoi(sPageSize)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_ListBill_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	farmGroupId, _ := strconv.Atoi(sfarmGroupId)

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Bill_ListBill_03", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_ListBill_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Data, err = c.BillService.TakePage(clientId, page, pageSize, orderBy, keyword, billType, farmGroupId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Bill_ListBill_05", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}
