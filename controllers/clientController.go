package controllers

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/httputil"
	"boonmafarm/api/utils/jwtutil"
	"fmt"
	"net/http"

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
			eg.GET("", c.GetClient)
			eg.PUT("", c.UpdateClient)
		}
	}
}

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
