package router

import (
	"github.com/gofiber/fiber/v3"
)

func (r *Router) setupActivityRoutes(group fiber.Router) {
	// Farm-wide activity feed (fill/move/sell across all of the client's ponds)
	activity := group.Group("/activity")

	activity.Get("", r.handlers.ActivityHandler.GetActivityFeed)

	// Per-size-grade breakdown of one sale — the feed only carries the summed
	// total, and a sale's money is priced per grade.
	activity.Get("/:activityId/sell-details", r.handlers.ActivityHandler.GetActivitySellDetails)
}
