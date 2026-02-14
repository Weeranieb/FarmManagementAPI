package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupFarmRoutes(group fiber.Router) {
	// Farm routes
	farm := group.Group("/farm")

	// Farm CRUD operations
	farm.Post("", r.handlers.FarmHandler.AddFarm)
	farm.Get("/hierarchy", r.handlers.FarmHandler.GetFarmHierarchy)
	farm.Get("/:id", r.handlers.FarmHandler.GetFarm)
	farm.Get("", r.handlers.FarmHandler.GetFarmList)
	farm.Put("", r.handlers.FarmHandler.UpdateFarm)
}
