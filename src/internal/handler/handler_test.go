package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

// userContextFromRequest sets UserContext from the request context so req.WithContext() is honored
func userContextFromRequest(c *fiber.Ctx) error {
	c.SetUserContext(c.Context())
	return c.Next()
}

// withUserContext returns a context with username, clientId, userLevel set for testing
func withUserContext(username string, clientId, userLevel int) context.Context {
	ctx := context.Background()
	if username != "" {
		ctx = context.WithValue(ctx, utils.UsernameKey(), username)
	}
	if clientId != 0 {
		ctx = context.WithValue(ctx, utils.ClientIdKey(), clientId)
	}
	if userLevel != 0 {
		ctx = context.WithValue(ctx, utils.UserLevelKey(), userLevel)
	}
	return ctx
}
