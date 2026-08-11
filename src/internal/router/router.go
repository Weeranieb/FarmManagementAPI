package router

import (
	"strings"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/middleware"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/metrics"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"
)

type Router struct {
	handlers *handler.Handler
	conf     *config.Config
	db       *gorm.DB
}

func SetupRoutes(app *fiber.App, conf *config.Config, handlers *handler.Handler, db *gorm.DB) {
	r := &Router{
		handlers: handlers,
		conf:     conf,
		db:       db,
	}

	// ── Global middleware ──────────────────────────────────────────────
	// Recover + AccessLog + metrics are registered via logging.UseHTTP in cmd/api and app.

	// 1. Helmet — security headers (X-Frame-Options, X-Content-Type-Options, etc.)
	app.Use(helmet.New())

	// 2. CORS — restrict allowed origins in production.
	// v3 takes []string configs and panics on "*" + credentials, so split the
	// configured origins and drop credentials whenever a wildcard is present.
	corsOrigins := conf.Cors.AllowedOrigins
	if corsOrigins == "" {
		corsOrigins = "*"
	}
	var allowOrigins []string
	allowCredentials := true
	for _, origin := range strings.Split(corsOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		if origin == "*" {
			allowCredentials = false
		}
		allowOrigins = append(allowOrigins, origin)
	}
	app.Use(cors.New(cors.Config{
		AllowCredentials: allowCredentials,
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "traceparent", "tracestate"},
		ExposeHeaders:    []string{"X-Request-ID"},
	}))

	// ── Public routes (no rate limit) ─────────────────────────────────
	r.setupPublicRoutes(app)

	// ── API routes (with rate limiter) ────────────────────────────────
	api := app.Group("/api/v1")

	// 3. Rate limiter — applied only to /api/v1 routes
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
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
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
	// Swagger documentation — served through the net/http Swagger UI wrapped for
	// Fiber v3, since the gofiber swagger packages still target v2.
	app.Get("/swagger/*", adaptor.HTTPHandlerFunc(httpSwagger.WrapHandler))

	// Liveness (cheap) / readiness (DB) / Prometheus text
	app.Get("/health", handler.Live)
	app.Get("/ready", handler.Ready(r.db))
	app.Get("/metrics", metrics.Handler())
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
	r.setupActivityRoutes(protected)
}
