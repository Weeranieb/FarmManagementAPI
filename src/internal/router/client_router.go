package router

import "github.com/gofiber/fiber/v2"

func (r *Router) setupClientRoutes(group fiber.Router) {
	// Client routes
	client := group.Group("/client")

	// Client CRUD operations
	client.Post("", r.handlers.ClientHandler.AddClient)
	client.Get("/:id", r.handlers.ClientHandler.GetClient)
	client.Put("", r.handlers.ClientHandler.UpdateClient)
}
