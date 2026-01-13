package router

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Router struct {
	handlers *handler.Handler
}

func NewRouter() *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*", // In production, specify exact origins
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	}))

	return app
}

func SetupRoutes(app *fiber.App, conf *config.Config, handlers *handler.Handler) {
	r := &Router{handlers: handlers}

	r.setupPublicRoutes(app)

	api := app.Group("/api/v1")
	r.setupPublicAPIRoutes(api)
	r.setupProtectedRoutes(api)
}

func (r *Router) setupPublicRoutes(app *fiber.App) {
	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}

func (r *Router) setupPublicAPIRoutes(api fiber.Router) {
	// Setup auth routes (public, no JWT required)
	r.setupAuthRoutes(api)

}

func (r *Router) setupProtectedRoutes(api fiber.Router) {
	// Protected routes (require JWT authentication)
	protected := api.Group("", middleware.JWTAuthMiddleware())

	r.setupUserRoutes(protected)
	r.setupClientRoutes(protected)
	r.setupFarmRoutes(protected)
	r.setupMerchantRoutes(protected)
	r.setupPondRoutes(protected)
	r.setupWorkerRoutes(protected)
	r.setupFeedCollectionRoutes(protected)
	r.setupFeedPriceHistoryRoutes(protected)
}
