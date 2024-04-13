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
			eg.GET("", c.GetUsers)
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

	newUser, err := c.UserService.Create(addUser, "", 1)
	if err != nil {
		httputil.NewError(ctx, "Err_User_AddUser_03", err)
		return
	}

	response.Result = true
	response.Data = newUser

	ctx.JSON(http.StatusOK, response)
}

func (c userControllerImp) GetUsers(ctx *gin.Context) {
	var response httputil.ResponseModel

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_User_GetUsers_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	// get userId
	id, err := jwtutil.GetUserId(ctx)
	if err != nil {
		httputil.NewError(ctx, "Err_User_GetUsers_02", err)
		return
	}

	users, err := c.UserService.GetUser(id)
	if err != nil {
		httputil.NewError(ctx, "Err_User_GetUsers_03", err)
		return
	}

	response.Result = true
	response.Data = users

	ctx.JSON(http.StatusOK, response)
}
