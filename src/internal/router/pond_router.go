package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupPondRoutes(group fiber.Router) {
	// Pond routes
	pond := group.Group("/pond")

	// Pond CRUD operations
	// Note: More specific routes (e.g. /:pondId/fill) before param routes (/:id) so /pond/1/fill is not matched as id
	pond.Post("", r.handlers.PondHandler.AddPonds)
	pond.Post("/:pondId/fill", r.handlers.PondHandler.FillPond)
	pond.Get("/:id", r.handlers.PondHandler.GetPond)
	pond.Put("/:id", r.handlers.PondHandler.UpdatePond)
	pond.Delete("/:id", r.handlers.PondHandler.DeletePond)
	pond.Get("", r.handlers.PondHandler.GetPondList)
}
