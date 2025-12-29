package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupUserRoutes(group fiber.Router) {
	// User routes
	user := group.Group("/user")

	// User CRUD operations matching old controller
	user.Post("", r.handlers.UserHandler.AddUser)
	user.Get("", r.handlers.UserHandler.GetUser)
	user.Put("", r.handlers.UserHandler.UpdateUser)
	user.Get("/list", r.handlers.UserHandler.GetUserList)
}
