package router

import "github.com/gofiber/fiber/v2"

func (r *Router) setupDailyLogRoutes(group fiber.Router) {
	pond := group.Group("/pond")

	pond.Post("/:pondId/daily-logs/upload", r.handlers.DailyLogHandler.UploadExcel)
	pond.Get("/:pondId/daily-logs", r.handlers.DailyLogHandler.GetMonth)
	pond.Put("/:pondId/daily-logs", r.handlers.DailyLogHandler.BulkUpsert)
}
