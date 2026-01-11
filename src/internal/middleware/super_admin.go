package middleware

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

// SuperAdminMiddleware checks if the user is a super admin
// This middleware should be used after JWTAuthMiddleware
func SuperAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
		if err != nil || !isSuperAdmin {
			return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
		}

		return c.Next()
	}
}
