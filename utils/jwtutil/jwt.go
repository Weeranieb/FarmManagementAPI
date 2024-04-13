package jwtutil

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserId(ctx *gin.Context) (int, error) {
	userId, ok := ctx.Get("userId")
	if !ok {
		return 0, errors.New("user id not found")
	}
	return int(userId.(float64)), nil
}

func GetClientId(ctx *gin.Context) (int, error) {
	clientId, ok := ctx.Get("clientId")
	if !ok {
		return 0, errors.New("client id not found")
	}
	return int(clientId.(float64)), nil
}

func GetUserLevel(ctx *gin.Context) (int, error) {
	userLevel, ok := ctx.Get("userLevel")
	if !ok {
		return 0, errors.New("user level not found")
	}
	return int(userLevel.(float64)), nil
}

func GetUsername(ctx *gin.Context) (string, error) {
	username, ok := ctx.Get("username")
	if !ok {
		return "", errors.New("username not found")
	}
	return username.(string), nil
}
