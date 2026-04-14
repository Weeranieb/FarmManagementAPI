package router

import "github.com/gofiber/fiber/v2"

func (r *Router) setupScanRoutes(group fiber.Router) {
	pond := group.Group("/pond")
	pond.Post("/:pondId/daily-logs/scan", r.handlers.ScanHandler.ScanDailyLog)
}
