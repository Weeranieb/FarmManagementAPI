package router

import (
	"github.com/gofiber/fiber/v3"
)

func (r *Router) setupActivityRoutes(group fiber.Router) {
	// Farm-wide activity feed (fill/move/sell across all of the client's ponds)
	activity := group.Group("/activity")

	activity.Get("", r.handlers.ActivityHandler.GetActivityFeed)
}
