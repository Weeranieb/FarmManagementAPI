package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupFeedCollectionRoutes(group fiber.Router) {
	// Feed collection routes
	feedCollection := group.Group("/feed-collection")

	// Feed collection CRUD operations
	// Note: More specific routes (/:id) must come before less specific routes ("")
	feedCollection.Post("", r.handlers.FeedCollectionHandler.AddFeedCollection)
	feedCollection.Get("/:id", r.handlers.FeedCollectionHandler.GetFeedCollection)
	feedCollection.Put("", r.handlers.FeedCollectionHandler.UpdateFeedCollection)
	feedCollection.Get("", r.handlers.FeedCollectionHandler.ListFeedCollection)
}
