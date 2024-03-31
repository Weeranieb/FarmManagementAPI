package controllers

import (
	"boonmafarm/api/pkg/models"
	"boonmafarm/api/pkg/services"
	"boonmafarm/api/utils/httputil"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	ApplyRoute(router *gin.Engine)
}

type authControllerImp struct {
	AuthService services.IUserService
}

func NewAuthController(authService services.IAuthService) IAuthController {
	return &authControllerImp{
		AuthService: authService,
	}
}

func (c authControllerImp) ApplyRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("/auth")
		{
			eg.POST("register", c.Register)
		}
	}
}

func (c authControllerImp) Register(ctx *gin.Context) {
	var response httputil.ResponseModel
	var addUser models.AddUser

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Auth_Register_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&addUser); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Auth_Register_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	newUser, err := c.AuthService.Create(addUser, "", 1)
	if err != nil {
		httputil.NewError(ctx, "Err_Auth_Register_03", err)
		return
	}

	response.Result = true
	response.Data = newUser

	ctx.JSON(http.StatusOK, response)
}
