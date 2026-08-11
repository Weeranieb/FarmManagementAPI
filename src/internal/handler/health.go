package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/metrics"
	"github.com/weeranieb/boonmafarm-backend/src/internal/version"
	"gorm.io/gorm"
)

// Live is process liveness — cheap, no dependency checks.
func Live(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"version": version.Version,
	})
}

// Ready checks DB reachability. Returns 503 when the dependency is down.
func Ready(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		if db == nil {
			metrics.IncDBPingFailure()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "unavailable",
				"version": version.Version,
				"error":   "database",
			})
		}
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			metrics.IncDBPingFailure()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "unavailable",
				"version": version.Version,
				"error":   "database",
			})
		}
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": version.Version,
		})
	}
}
