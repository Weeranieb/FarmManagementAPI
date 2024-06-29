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
	AuthService services.IAuthService
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
			eg.POST("login", c.Login)
		}
	}
}

// POST /api/v1/auth/register
// Register a new user.
// @Summary      Register a new user
// @Description  Register a new user with the provided details
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.AddUser true "User data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/auth/register [post]
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

	newUser, err := c.AuthService.Create(addUser)
	if err != nil {
		httputil.NewError(ctx, "Err_Auth_Register_03", err)
		return
	}

	response.Result = true
	response.Data = newUser

	ctx.JSON(http.StatusOK, response)
}

// POST /api/v1/auth/login
// Login user and return JWT token.
// @Summary      Login user
// @Description  Login user with provided credentials and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.Login true "Login data"
// @Success      200  {object}  httputil.ResponseModel
// @Failure      400  {object}  httputil.ErrorResponseModel
// @Failure      401  {object}  httputil.ErrorResponseModel
// @Failure      500  {object}  httputil.ErrorResponseModel
// @Router       /api/v1/auth/login [post]
func (c authControllerImp) Login(ctx *gin.Context) {
	var response httputil.ResponseModel
	var login models.Login

	defer func() {
		if r := recover(); r != nil {
			errRes := httputil.ErrorResponseModel{}
			errRes.Error(ctx, "Err_Auth_Login_01", fmt.Sprint(r))
			response.Error = errRes
			ctx.JSON(http.StatusOK, response)
			return
		}
	}()

	if err := ctx.ShouldBindJSON(&login); err != nil {
		errRes := httputil.ErrorResponseModel{}
		errRes.Error(ctx, "Err_Auth_Login_02", err.Error())
		response.Error = errRes
		ctx.JSON(http.StatusOK, response)
		return
	}

	token, user, err := c.AuthService.Login(login)
	if err != nil {
		httputil.NewError(ctx, "Err_Auth_Login_03", err)
		return
	}

	type dataPayload struct {
		AccessToken string       `json:"accessToken"`
		User        *models.User `json:"user"`
	}

	response.Result = true
	response.Data = dataPayload{
		AccessToken: token,
		User:        user,
	}

	ctx.JSON(http.StatusOK, response)
}
