package router

import "github.com/gofiber/fiber/v2"

func (r *Router) setupDailyFeedRoutes(group fiber.Router) {
	pond := group.Group("/pond")

	pond.Post("/:pondId/daily-feed/upload", r.handlers.DailyFeedHandler.UploadExcel)
	pond.Get("/:pondId/daily-feed", r.handlers.DailyFeedHandler.GetMonth)
	pond.Put("/:pondId/daily-feed", r.handlers.DailyFeedHandler.BulkUpsert)
	pond.Delete("/:pondId/daily-feed/:feedCollectionId", r.handlers.DailyFeedHandler.DeleteTable)
}
