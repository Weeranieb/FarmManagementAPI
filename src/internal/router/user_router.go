package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupUserRoutes(group fiber.Router) {
	// User routes
	user := group.Group("/user")

	// User CRUD operations
	user.Post("", r.handlers.UserHandler.AddUser)                // super-admin: create user
	user.Get("", r.handlers.UserHandler.GetUser)                 // self
	user.Put("", r.handlers.UserHandler.UpdateUser)              // self-update
	user.Get("/list", r.handlers.UserHandler.GetUserList)        // list (filtered)
	user.Put("/:id", r.handlers.UserHandler.AdminUpdateUser)     // super-admin: admin-update
	user.Delete("/:id", r.handlers.UserHandler.DeleteUser)       // super-admin: soft delete
}
