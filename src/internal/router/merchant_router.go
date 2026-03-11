package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) setupMerchantRoutes(group fiber.Router) {
	// Merchant routes
	merchant := group.Group("/merchant")

	// Merchant CRUD operations
	// Note: More specific routes (/:id) must come before less specific routes ("")
	merchant.Post("", r.handlers.MerchantHandler.AddMerchant)
	merchant.Get("/:id", r.handlers.MerchantHandler.GetMerchant)
	merchant.Delete("/:id", r.handlers.MerchantHandler.DeleteMerchant)
	merchant.Get("", r.handlers.MerchantHandler.GetMerchantList)
	merchant.Put("", r.handlers.MerchantHandler.UpdateMerchant)
}
