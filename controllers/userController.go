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

type IUserController interface {
	ApplyRoute(router *gin.Engine)
}

type userControllerImp struct {
	UserService services.IUserService
}

func NewUserController(userService services.IUserService) IUserController {
	return &userControllerImp{
		UserService: userService,
	}
}

func (c userControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/user")
		{
			eg.POST("", c.AddUser)
			eg.GET("", c.GetUser)
			eg.PUT("", c.UpdateUser)
			eg.GET("list", c.GetUserList)
		}
	}
}

func (c userControllerImp) AddUser(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addUser models.AddUser

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_User_AddUser_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addUser); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_AddUser_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_AddUser_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get clientId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_AddUser_04", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newUser, err := c.UserService.Create(addUser, username, clientId)
	if err != nil {
		httputil.NewError(ctx, "Err_User_AddUser_03", err)
		return
	}

	response.Result = true
	response.Data = newUser

	ctx.JSON(http.StatusOK, response)
}

func (c userControllerImp) GetUser(ctx *gin.Context) {
	var response httputil.ResponseModel

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_User_GetUser_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	// get userId
	id, err := jwtutil.GetUserId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_GetUser_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	user, err := c.UserService.GetUser(id)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_GetUser_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = user

	ctx.JSON(http.StatusOK, response)
}

func (c userControllerImp) UpdateUser(ctx *gin.Context) {
	var response httputil.ResponseModel
	var updateUser *models.User

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_User_UpdateUser_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_UpdateUser_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	// get username
	username, err := jwtutil.GetUsername(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_UpdateUser_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	err = c.UserService.Update(updateUser, username)
	if err != nil {
		httputil.NewError(ctx, "Err_User_UpdateUser_04", err)
		return
	}

	response.Result = true

	ctx.JSON(http.StatusOK, response)
}

func (c userControllerImp) GetUserList(ctx *gin.Context) {
	var response httputil.ResponseModel

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_User_GetUserList_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	// get clientId
	clientId, err := jwtutil.GetClientId(ctx)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_GetUserList_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	users, err := c.UserService.GetUserList(clientId)
	if err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_User_GetUserList_03", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	response.Result = true
	response.Data = users

	ctx.JSON(http.StatusOK, response)
}
