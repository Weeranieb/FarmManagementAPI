package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupFarmGroupRoutes(group fiber.Router) {
	farmGroup := group.Group("/farm-group")

	farmGroup.Post("", r.handlers.FarmGroupHandler.AddFarmGroup)
	farmGroup.Get("/dropdown", r.handlers.FarmGroupHandler.GetFarmGroupDropdown)
	farmGroup.Get("/:id", r.handlers.FarmGroupHandler.GetFarmGroup)
	farmGroup.Put("", r.handlers.FarmGroupHandler.UpdateFarmGroup)
	farmGroup.Get("", r.handlers.FarmGroupHandler.ListFarmGroup)
}
