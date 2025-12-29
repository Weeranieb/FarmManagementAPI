package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupPondRoutes(group fiber.Router) {
	// Pond routes
	pond := group.Group("/pond")

	// Pond CRUD operations
	// Note: More specific routes (/:id) must come before less specific routes ("")
	pond.Post("", r.handlers.PondHandler.AddPond)
	pond.Post("/batch", r.handlers.PondHandler.AddPonds)
	pond.Get("/:id", r.handlers.PondHandler.GetPond)
	pond.Put("", r.handlers.PondHandler.UpdatePond)
	pond.Delete("/:id", r.handlers.PondHandler.DeletePond)
	pond.Get("", r.handlers.PondHandler.GetPondList)
}
