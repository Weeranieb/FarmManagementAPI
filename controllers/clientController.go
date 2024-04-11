package controllers

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/httputil"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IClientController interface {
	ApplyRoute(router *gin.Engine)
}

type ClientControllerImp struct {
	ClientService services.IClientService
}

func NewClientController(clientService services.IClientService) IClientController {
	return &ClientControllerImp{
		ClientService: clientService,
	}
}

func (c ClientControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/client")
		{
			eg.POST("", c.AddClient)
			eg.GET("/:id", c.GetClient)
		}
	}
}

func (c ClientControllerImp) AddClient(ctx *gin.Context) {
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

func (c ClientControllerImp) GetClient(ctx *gin.Context) {
	var response httputil.ResponseModel
	id := ctx.Param("id")

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Client_GetClient_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	clientId, err := strconv.Atoi(id)
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
