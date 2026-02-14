package utils

import (
	"context"
	"errors"

	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
)

func GetUserId(ctx context.Context) (int, error) {
	userId := ctx.Value(constants.UserIDKey)
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

func GetClientId(ctx context.Context) *int {
	clientId := ctx.Value(constants.ClientIDKey)
	if clientId == nil {
		return nil
	}

	var id int
	switch v := clientId.(type) {
	case int:
		id = v
	case float64:
		id = int(v)
	case int64:
		id = int(v)
	default:
		return nil
	}
	return &id
}

func GetClientIdForAccess(ctx context.Context) (*int, bool) {
	isSuperAdmin, err := IsSuperAdmin(ctx)
	if err != nil {
		return nil, false
	}

	// Super admin doesn't need clientId restriction
	if isSuperAdmin {
		return nil, true
	}

	// Regular users must have clientId
	clientId := GetClientId(ctx)
	return clientId, clientId != nil && *clientId != 0
}

func GetUserLevel(ctx context.Context) (int, error) {
	userLevel := ctx.Value(constants.UserLevelKey)
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
	username := ctx.Value(constants.UsernameKey)
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
	userClientId := GetClientId(ctx)
	if userClientId == nil {
		return false, nil
	}

	return *userClientId == targetClientId, nil
}
