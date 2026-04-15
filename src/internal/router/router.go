package router

import (
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Router struct {
	handlers *handler.Handler
	conf     *config.Config
}

func SetupRoutes(app *fiber.App, conf *config.Config, handlers *handler.Handler) {
	r := &Router{
		handlers: handlers,
		conf:     conf,
	}

	// ── Global middleware ──────────────────────────────────────────────
	// 1. Recover — catch panics and return 500 instead of crashing
	app.Use(recover.New())

	// 2. Logger — log every request
	app.Use(logger.New())

	// 3. Helmet — security headers (X-Frame-Options, X-Content-Type-Options, etc.)
	app.Use(helmet.New())

	// 4. CORS — restrict allowed origins in production
	corsOrigins := conf.Cors.AllowedOrigins
	if corsOrigins == "" {
		corsOrigins = "*"
	}
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     corsOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	}))

	// ── Public routes (no rate limit) ─────────────────────────────────
	r.setupPublicRoutes(app)

	// ── API routes (with rate limiter) ────────────────────────────────
	api := app.Group("/api/v1")

	// 5. Rate limiter — applied only to /api/v1 routes
	window := time.Duration(conf.Security.RateLimitWindow) * time.Second
	if window <= 0 {
		window = 60 * time.Second
	}
	maxReqs := conf.Security.RateLimitMax
	if maxReqs <= 0 {
		maxReqs = 100
	}
	api.Use(limiter.New(limiter.Config{
		Max:        maxReqs,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    "429",
				"message": "Too many requests. Please try again later.",
			})
		},
	}))

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
	protected := api.Group("", middleware.JWTAuthMiddleware(r.conf.Authentication.JWTSecret))

	r.setupUserRoutes(protected)
	r.setupClientRoutes(protected)
	r.setupFarmRoutes(protected)
	r.setupMerchantRoutes(protected)
	r.setupPondRoutes(protected)
	r.setupFishSizeGradeRoutes(protected)
	r.setupFarmGroupRoutes(protected)
	r.setupWorkerRoutes(protected)
	r.setupFeedCollectionRoutes(protected)
	r.setupFeedPriceHistoryRoutes(protected)
	r.setupDailyLogRoutes(protected)
}
