package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupUserRoutes(group fiber.Router) {
	// User routes
	user := group.Group("/user")

	// User CRUD operations
	user.Post("", r.handlers.UserHandler.AddUser)                        // super-admin: create user
	user.Get("", r.handlers.UserHandler.GetUser)                         // self
	user.Put("", r.handlers.UserHandler.UpdateUser)                      // self-update
	user.Put("/password", r.handlers.UserHandler.ChangePassword)         // self change password
	user.Get("/list", r.handlers.UserHandler.GetUserList)                // list (filtered)
	user.Put("/:id", r.handlers.UserHandler.AdminUpdateUser)             // super-admin: admin-update
	user.Put("/:id/password", r.handlers.UserHandler.AdminResetPassword) // super-admin: reset password
	user.Delete("/:id", r.handlers.UserHandler.DeleteUser)               // super-admin: soft delete
}
