package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupAuthRoutes(group fiber.Router) {
	// Auth routes
	auth := group.Group("/auth")

	// Auth operations
	auth.Post("/register", r.handlers.AuthHandler.Register)
	auth.Post("/login", r.handlers.AuthHandler.Login)
}
