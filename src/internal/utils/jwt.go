package utils

import (
	"context"
	"errors"
)

type contextKey string

const (
	userIdKey    contextKey = "userId"
	clientIdKey  contextKey = "clientId"
	userLevelKey contextKey = "userLevel"
	usernameKey  contextKey = "username"
)

// Context key getters for use in middleware
func UserIdKey() contextKey    { return userIdKey }
func ClientIdKey() contextKey  { return clientIdKey }
func UserLevelKey() contextKey { return userLevelKey }
func UsernameKey() contextKey  { return usernameKey }

func GetUserId(ctx context.Context) (int, error) {
	userId := ctx.Value(userIdKey)
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

func GetClientId(ctx context.Context) (int, error) {
	clientId := ctx.Value(clientIdKey)
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

func GetUserLevel(ctx context.Context) (int, error) {
	userLevel := ctx.Value(userLevelKey)
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

func GetUsername(ctx context.Context) (string, error) {
	username := ctx.Value(usernameKey)
	if username == nil {
		return "", errors.New("username not found")
	}

	if v, ok := username.(string); ok {
		return v, nil
	}
	return "", errors.New("invalid username type")
}

// IsSuperAdmin checks if the current user is a super admin
func IsSuperAdmin(ctx context.Context) (bool, error) {
	userLevel, err := GetUserLevel(ctx)
	if err != nil {
		return false, err
	}
	return userLevel >= 3, nil
}

// IsClientAdminOrAbove checks if user is client admin or super admin
func IsClientAdminOrAbove(ctx context.Context) (bool, error) {
	userLevel, err := GetUserLevel(ctx)
	if err != nil {
		return false, err
	}
	return userLevel >= 2, nil
}

// CanAccessClient checks if user can access a specific client
func CanAccessClient(ctx context.Context, targetClientId int) (bool, error) {
	userLevel, err := GetUserLevel(ctx)
	if err != nil {
		return false, err
	}

	// Super admin can access all clients
	if userLevel >= 3 {
		return true, nil
	}

	// Others can only access their own client
	userClientId, err := GetClientId(ctx)
	if err != nil {
		return false, err
	}

	return userClientId == targetClientId, nil
}
