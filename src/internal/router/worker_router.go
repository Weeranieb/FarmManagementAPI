package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupWorkerRoutes(group fiber.Router) {
	// Worker routes
	worker := group.Group("/worker")

	// Worker CRUD operations
	// Note: More specific routes (/:id) must come before less specific routes ("")
	worker.Post("", r.handlers.WorkerHandler.AddWorker)
	worker.Get("/:id", r.handlers.WorkerHandler.GetWorker)
	worker.Put("", r.handlers.WorkerHandler.UpdateWorker)
	worker.Get("", r.handlers.WorkerHandler.ListWorker)
}

