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

type IClientController interface {
	ApplyRoute(router *gin.Engine)
}

type clientControllerImp struct {
	ClientService services.IClientService
}

func NewClientController(clientService services.IClientService) IClientController {
	return &clientControllerImp{
		ClientService: clientService,
	}
}

func (c clientControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/client")
		{
			eg.POST("", c.AddClient)
			eg.GET("/:clientId/farms", c.GetClientWithFarms)
			eg.GET("/farms", c.GetAllClientWithFarms)
			eg.GET("", c.GetClient)
			eg.GET("list", c.GetAllClient)
			eg.PUT("", c.UpdateClient)
		}
	}
}

// POST /api/v1/client
// Add a new client.
// @Summary      Add a new client
// @Description  Add a new client with the provided details
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        body body models.AddClient true "Client data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/client [post]
func (c clientControllerImp) AddClient(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addClient models.AddClient

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Client_AddClient_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addClient); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_AddClient_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newClient, err := c.ClientService.Create(addClient, "")
	if err != nil {
		httputil.NewError(ctx, "Err_User_AddUser_03", err)
		return
	}

	response.Result = true
	response.Data = newClient

	ctx.JSON(http.StatusOK, response)
}

// GET /api/v1/client
// Get client details.
// @Summary      Get client details
// @Description  Retrieve details of the currently logged-in client
// @Tags         client
// @Accept       json
// @Produce      json
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/client [get]
func (c clientControllerImp) GetClient(ctx *gin.Context) {
	var response httputil.ResponseModel

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Client_GetClient_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	// get userId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_GetClient_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	client, err := c.ClientService.Get(clientId)
	if err != nil {
		httputil.NewError(ctx, "Err_Client_GetClient_03", err)
		return
	}

	response.Result = true
	response.Data = client

	ctx.JSON(http.StatusOK, response)
}

// PUT /api/v1/client
// Update client details.
// @Summary      Update client details
// @Description  Update details of the currently logged-in client
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        body body models.Client true "Client data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/client [put]
func (c clientControllerImp) UpdateClient(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateClient models.Client

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Client_UpdateClient_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateClient); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_UpdateClient_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	userName, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_UpdateClient_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.ClientService.Update(&updateClient, userName)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_UpdateClient_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}

func (c clientControllerImp) GetClientWithFarms(ctx *gin.Context) {
	var response httputil.ResponseModel
	sClientId := ctx.Param("clientId")

	clientId, err := strconv.Atoi(sClientId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_GetClientWithFarms_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Client_GetClientWithFarms_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	userLevel, err := jwtutil.GetUserLevel(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_GetClientWithFarms_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	clientWithFarm, err := c.ClientService.GetClientWithFarms(userLevel, &clientId)
	if err != nil {
		httputil.NewError(ctx, "Err_Client_GetClientWithFarms_04", err)
		return
	}

	response.Result = true
	response.Data = clientWithFarm

	ctx.JSON(http.StatusOK, response)
}

func (c clientControllerImp) GetAllClientWithFarms(ctx *gin.Context) {
	var response httputil.ResponseModel

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Client_GetAllClientWithFarms_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	userLevel, err := jwtutil.GetUserLevel(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_GetAllClientWithFarms_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	clientWithFarms, err := c.ClientService.GetClientWithFarms(userLevel, nil)

	if err != nil {
		httputil.NewError(ctx, "Err_Client_GetClientWithFarms_04", err)
		return
	}

	response.Result = true
	response.Data = clientWithFarms

	ctx.JSON(http.StatusOK, response)
}

func (c clientControllerImp) GetAllClient(ctx *gin.Context) {
	var response httputil.ResponseModel
	keyword := ctx.Query("keyword")

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Client_GetAllClient_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	userLevel, err := jwtutil.GetUserLevel(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Client_GetAllClient_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	clients, err := c.ClientService.GetAllClient(userLevel, keyword)
	if err != nil {
		httputil.NewError(ctx, "Err_Client_GetAllClient_03", err)
		return
	}

	response.Result = true
	response.Data = clients

	ctx.JSON(http.StatusOK, response)
}
