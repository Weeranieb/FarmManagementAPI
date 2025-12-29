package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupFeedPriceHistoryRoutes(group fiber.Router) {
	feedPriceHistory := group.Group("/feed-price-history")

	feedPriceHistory.Post("", r.handlers.FeedPriceHistoryHandler.AddFeedPriceHistory)
	feedPriceHistory.Get("/:id", r.handlers.FeedPriceHistoryHandler.GetFeedPriceHistory)
	feedPriceHistory.Put("", r.handlers.FeedPriceHistoryHandler.UpdateFeedPriceHistory)
	feedPriceHistory.Get("", r.handlers.FeedPriceHistoryHandler.GetAllFeedPriceHistory)
}
