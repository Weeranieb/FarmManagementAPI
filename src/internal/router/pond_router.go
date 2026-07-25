package router

import (
	"github.com/gofiber/fiber/v3"
)

func (r *Router) setupPondRoutes(group fiber.Router) {
	// Pond routes
	pond := group.Group("/pond")

	// Pond CRUD operations
	pond.Get("/template", r.handlers.PondHandler.DownloadTemplate)
	pond.Post("/bulk-import/:clientId", r.handlers.PondHandler.BulkImportFarmPond)
	pond.Post("", r.handlers.PondHandler.AddPonds)
	pond.Post("/fill/calc", r.handlers.PondHandler.FillPondCalc)
	pond.Post("/move/calc", r.handlers.PondHandler.MovePondCalc)
	pond.Post("/sell/calc", r.handlers.PondHandler.SellPondCalc)
	pond.Post("/:pondId/fill/preview", r.handlers.PondHandler.FillPondPreview)
	pond.Post("/:pondId/move/preview", r.handlers.PondHandler.MovePondPreview)
	pond.Post("/:pondId/sell/preview", r.handlers.PondHandler.SellPondPreview)
	pond.Post("/:pondId/fill", r.handlers.PondHandler.FillPond)
	pond.Post("/:pondId/move", r.handlers.PondHandler.MovePond)
	pond.Post("/:pondId/sell", r.handlers.PondHandler.SellPond)
	pond.Get("/:pondId/activities", r.handlers.PondHandler.GetPondActivities)
	pond.Get("/:pondId/cycles", r.handlers.PondHandler.GetPondCycles)
	pond.Get("/:id", r.handlers.PondHandler.GetPond)
	pond.Put("/:id", r.handlers.PondHandler.UpdatePond)
	pond.Delete("/:id", r.handlers.PondHandler.DeletePond)
	pond.Get("", r.handlers.PondHandler.GetPondList)
}
