package router

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupClientRoutes(group fiber.Router) {
	// Client routes
	client := group.Group("/client")

	// Client CRUD operations
	// AddClient requires super admin
	client.Post("", middleware.SuperAdminMiddleware(), r.handlers.ClientHandler.AddClient)
	// List for dropdown (super admin only) - must be before /:id
	client.Get("/list", r.handlers.ClientHandler.GetClientList)
	client.Get("/:id", r.handlers.ClientHandler.GetClient)
	client.Put("", r.handlers.ClientHandler.UpdateClient)
}
