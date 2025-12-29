package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func GetUserId(c *fiber.Ctx) (int, error) {
	userId := c.Locals("userId")
	if userId == nil {
		return 0, errors.New("user id not found")
	}
	
	switch v := userId.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case int64:
		return int(v), nil
	default:
		return 0, errors.New("invalid user id type")
	}
}

func GetClientId(c *fiber.Ctx) (int, error) {
	clientId := c.Locals("clientId")
	if clientId == nil {
		return 0, errors.New("client id not found")
	}
	
	switch v := clientId.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case int64:
		return int(v), nil
	default:
		return 0, errors.New("invalid client id type")
	}
}

func GetUserLevel(c *fiber.Ctx) (int, error) {
	userLevel := c.Locals("userLevel")
	if userLevel == nil {
		return 0, errors.New("user level not found")
	}
	
	switch v := userLevel.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case int64:
		return int(v), nil
	default:
		return 0, errors.New("invalid user level type")
	}
}

func GetUsername(c *fiber.Ctx) (string, error) {
	username := c.Locals("username")
	if username == nil {
		return "", errors.New("username not found")
	}
	
	if v, ok := username.(string); ok {
		return v, nil
	}
	return "", errors.New("invalid username type")
}

