package controllers

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/httputil"
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
		eg := v1.Group("/users")
		{
			eg.POST("add", c.AddUser)
		}
	}
}

func (c userControllerImp) AddUser(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addUser models.AddUsers

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

	newUser, err := c.UserService.Create(addUser, "")
	if err != nil {
		httputil.NewError(ctx, "Err_User_AddUser_03", err)
		return
	}

	response.Result = true
	response.Data = newUser

	ctx.JSON(http.StatusOK, response)
}
