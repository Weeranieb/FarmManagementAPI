package router

import (
	"github.com/gofiber/fiber/v3"
)

func (r *Router) setupFishSizeGradeRoutes(group fiber.Router) {
	fishSizeGrade := group.Group("/fish-size-grade")

	fishSizeGrade.Get("/dropdown", r.handlers.FishSizeGradeHandler.GetDropdown)
}
